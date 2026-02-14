package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/go-github/v82/github"
	"github.com/google/uuid"

	"xeiaso.net/v4/cmd/sponsor-panel/templates"
)

// inviteHandler handles POST /invite - invites a user to the GitHub team.
func (s *Server) inviteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	slog.Debug("inviteHandler: processing team invite request")

	// Get user from session
	user, err := s.getSessionUser(r)
	if err != nil {
		slog.Error("inviteHandler: failed to get session user", "err", err)
		renderError(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	slog.Debug("inviteHandler: authenticated user", "user_id", user.ID, "login", user.Login)

	// Check $50+ sponsorship tier (5000 cents)
	if !user.IsSponsorAtTier(5000) {
		slog.Error("inviteHandler: user not eligible for team invitation", "user", user.Login, "user_id", user.ID)
		renderError(w, "Requires $50+/month sponsorship", http.StatusForbidden)
		return
	}

	slog.Debug("inviteHandler: user is eligible for team invitation", "user", user.Login, "user_id", user.ID)

	// Parse form
	if err := r.ParseForm(); err != nil {
		slog.Error("inviteHandler: failed to parse form", "err", err)
		renderError(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	if username == "" {
		slog.Error("inviteHandler: empty username provided", "user_id", user.ID)
		renderError(w, "Username required", http.StatusBadRequest)
		return
	}

	// Strip @ if present
	username = strings.TrimPrefix(username, "@")

	slog.Debug("inviteHandler: inviting user to team", "invited_by", user.Login, "username_to_invite", username)

	// Invite to team using go-github
	teamSlug := "botstopper-customers"
	org := "TecharoHQ"

	membership, _, err := s.ghClient.Teams.AddTeamMembershipBySlug(
		r.Context(), org, teamSlug, username, nil,
	)

	if err != nil {
		slog.Error("inviteHandler: failed to invite to team", "user", username, "err", err, "invited_by", user.Login)
		// Check for common errors
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "422") {
			renderError(w, "User not found or already invited", http.StatusBadRequest)
			return
		}
		renderError(w, "Failed to invite: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Determine membership state
	state := "pending"
	if membership != nil && membership.State != nil && *membership.State == "active" {
		state = "active"
	}

	slog.Info("inviteHandler: team invitation successful",
		"invited_by", user.Login,
		"username_invited", username,
		"state", state,
		"team", teamSlug,
		"org", org)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	renderInviteSuccess(w, username, state)
}

// logoHandler handles POST /logo - submits a logo to a GitHub issue.
func (s *Server) logoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	slog.Debug("logoHandler: processing logo submission request")

	// Get user from session
	user, err := s.getSessionUser(r)
	if err != nil {
		slog.Error("logoHandler: failed to get session user", "err", err)
		renderError(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	slog.Debug("logoHandler: authenticated user", "user_id", user.ID, "login", user.Login)

	// Check user is a sponsor (any tier)
	if !user.IsSponsorAtTier(100) {
		slog.Error("logoHandler: user not a sponsor", "user", user.Login, "user_id", user.ID)
		renderError(w, "Requires active sponsorship", http.StatusForbidden)
		return
	}

	slog.Debug("logoHandler: user is eligible for logo submission", "user", user.Login, "user_id", user.ID)

	// Parse multipart form (5MB max)
	if err := r.ParseMultipartForm(5 * 1024 * 1024); err != nil {
		slog.Error("logoHandler: failed to parse multipart form", "err", err)
		renderError(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	companyName := r.FormValue("company")
	website := r.FormValue("website")

	if companyName == "" || website == "" {
		slog.Error("logoHandler: missing required fields", "user_id", user.ID, "company", companyName, "website", website)
		renderError(w, "Company name and website are required", http.StatusBadRequest)
		return
	}

	// Get uploaded file
	file, header, err := r.FormFile("logo")
	if err != nil {
		slog.Error("logoHandler: failed to get logo file", "err", err)
		renderError(w, "Logo file required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	slog.Debug("logoHandler: received logo file",
		"user_id", user.ID,
		"company", companyName,
		"filename", header.Filename,
		"size", header.Size)

	// Validate file size
	if header.Size > 5*1024*1024 {
		slog.Error("logoHandler: file too large", "user_id", user.ID, "size", header.Size)
		renderError(w, "File too large (max 5MB)", http.StatusBadRequest)
		return
	}

	// Read file into memory for S3 upload
	fileData, err := io.ReadAll(file)
	if err != nil {
		slog.Error("logoHandler: failed to read file", "err", err)
		renderError(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Upload to S3 with UUID v7 folder structure
	var logoURL string
	var s3Key string
	if s.bucketName != "" && s.s3Client != nil {
		// Generate UUID v7 for folder
		folderID := uuid.Must(uuid.NewV7())
		ext := path.Ext(header.Filename)
		s3Key = fmt.Sprintf("%s/%s%s", folderID.String(), "logo", ext)

		// Detect content type
		contentType := mime.TypeByExtension(ext)
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		slog.Debug("logoHandler: uploading to S3",
			"user_id", user.ID,
			"bucket", s.bucketName,
			"key", s3Key,
			"content_type", contentType)

		putInput := &s3.PutObjectInput{
			Bucket:      &s.bucketName,
			Key:         &s3Key,
			Body:        bytes.NewReader(fileData),
			ContentType: &contentType,
		}

		_, err := s.s3Client.PutObject(r.Context(), putInput)
		if err != nil {
			slog.Error("logoHandler: failed to upload to S3", "err", err, "user_id", user.ID)
			renderError(w, "Failed to upload logo: "+err.Error(), http.StatusInternalServerError)
			return
		}

		logoURL = fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucketName, s3Key)
		slog.Info("logoHandler: uploaded to S3",
			"user_id", user.ID,
			"url", logoURL,
			"key", s3Key)
	}

	// Create GitHub issue
	issueTitle := fmt.Sprintf("Logo Submission: %s", companyName)
	issueBody := fmt.Sprintf(`# Logo Submission: %s

**Submitted by:** @%s

## Details

- **Company:** %s
- **Website:** %s
- **File:** %s
- **File Size:** %d bytes
%s

## Next Steps

1. Review logo for quality and appropriateness
2. Add to Anubis README sponsors section
3. Close this issue

/label logo-submission
/label needs-review
`, companyName, user.Login, companyName, website, header.Filename, header.Size, func() string {
		if logoURL != "" {
			return fmt.Sprintf("- **S3 Bucket:** %s\n- **S3 Key:** `%s`\n- **Logo URL:** %s", s.bucketName, s3Key, logoURL)
		}
		return "- **Storage:** Not configured (bucket-name not set)"
	}())

	issue := &github.IssueRequest{
		Title:  github.Ptr(issueTitle),
		Body:   github.Ptr(issueBody),
		Labels: &[]string{"logo-submission", "needs-review"},
	}

	slog.Debug("logoHandler: creating GitHub issue", "user_id", user.ID, "company", companyName, "title", issueTitle)

	createdIssue, _, err := s.ghClient.Issues.Create(r.Context(), "TecharoHQ", *logoSubmissionRepo, issue)
	if err != nil {
		slog.Error("logoHandler: failed to create GitHub issue", "err", err, "user_id", user.ID, "company", companyName)
		renderError(w, "Failed to create issue: "+err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("logoHandler: GitHub issue created successfully",
		"user_id", user.ID,
		"login", user.Login,
		"company", companyName,
		"issue_number", createdIssue.GetNumber(),
		"issue_url", createdIssue.GetHTMLURL())

	// Store submission record
	submission := &LogoSubmission{
		UserID:            user.ID,
		CompanyName:       companyName,
		Website:           website,
		LogoURL:           logoURL,
		GitHubIssueURL:    createdIssue.GetHTMLURL(),
		GitHubIssueNumber: createdIssue.GetNumber(),
	}
	if err := createLogoSubmission(r.Context(), s.pool, submission); err != nil {
		slog.Error("logoHandler: failed to store submission", "err", err, "user_id", user.ID, "issue_number", createdIssue.GetNumber())
	} else {
		slog.Debug("logoHandler: submission stored in database", "user_id", user.ID, "issue_number", createdIssue.GetNumber())
	}

	slog.Info("logoHandler: logo submission completed",
		"company", companyName,
		"issue", createdIssue.GetNumber(),
		"user", user.Login,
		"user_id", user.ID)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	renderLogoSuccess(w, companyName, createdIssue.GetHTMLURL(), createdIssue.GetNumber())
}

// renderError renders an error message for HTMX.
func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(statusCode)
	templates.FormResult(message, false).Render(context.Background(), w)
}

// renderInviteSuccess renders success response for team invitation.
func renderInviteSuccess(w http.ResponseWriter, username, state string) {
	stateText := "Invited"
	if state == "active" {
		stateText = "Active Member"
	}
	templates.InviteSuccess(username, stateText).Render(context.Background(), w)
}

// renderLogoSuccess renders success response for logo submission.
func renderLogoSuccess(w http.ResponseWriter, company, issueURL string, issueNumber int) {
	templates.LogoSuccess(company, issueURL, issueNumber).Render(context.Background(), w)
}
