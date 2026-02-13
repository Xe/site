#!/usr/bin/env node

import fs from 'fs';
import path from 'path';
import { execSync } from 'child_process';
import matter from 'gray-matter';
import { Octokit } from '@octokit/rest';

/**
 * Parses the PR body for a "DO NOT MERGE until YYYY-MM-DD UTC" instruction
 * @param {string} prBody - The PR body text
 * @returns {string|null} - The date string if found, null otherwise
 */
function parseMergeUntilDate(prBody) {
  if (!prBody) return null;

  // Match patterns like:
  // - "DO NOT MERGE until 2026-02-20 UTC"
  // - "DO NOT MERGE UNTIL 2026-02-20 UTC"
  // - "do not merge until 2026-02-20 UTC"
  const pattern = /do not merge until (\d{4}-\d{2}-\d{2}) UTC/i;
  const match = prBody.match(pattern);

  return match ? match[1] : null;
}

/**
 * Posts a comment to the PR warning about the merge date
 * @param {Octokit} octokit - Octokit instance
 * @param {string} owner - Repository owner
 * @param {string} repo - Repository name
 * @param {number} prNumber - PR number
 * @param {string} mergeUntilDate - The date until which the PR should not be merged
 * @param {string} filesInfo - Information about the future-dated files
 */
async function postMergeWarningComment(octokit, owner, repo, prNumber, mergeUntilDate, filesInfo) {
  const body = `## ⏳ Scheduled Publication Detected

This PR contains content with future dates that are scheduled for publication:

${filesInfo}

**Do not merge this PR until ${mergeUntilDate} UTC.**

The content dates have been validated against the "DO NOT MERGE until" date in the PR description.

> [!NOTE]
> This is an automated message from the validate-blog-dates workflow.`;

  try {
    await octokit.rest.issues.createComment({
      owner,
      repo,
      issue_number: prNumber,
      body
    });
    console.log('  ✅ Posted warning comment to PR');
  } catch (error) {
    console.log(`  ⚠️  Failed to post comment: ${error.message}`);
  }
}

/**
 * Fetches the PR body using Octokit
 * @param {Octokit} octokit - Octokit instance
 * @param {string} owner - Repository owner
 * @param {string} repo - Repository name
 * @param {number} prNumber - PR number
 * @returns {Promise<string|null>} - The PR body or null
 */
async function fetchPRBody(octokit, owner, repo, prNumber) {
  try {
    const { data } = await octokit.rest.pulls.get({
      owner,
      repo,
      pull_number: prNumber
    });
    return data.body;
  } catch (error) {
    console.log(`Could not fetch PR body: ${error.message}`);
    return null;
  }
}

/**
 * Validates that content files don't have future dates
 */
async function validateContentDates() {
  const currentDate = new Date().toISOString().split('T')[0]; // YYYY-MM-DD format
  const isPullRequest = process.env.GITHUB_EVENT_NAME === 'pull_request';
  let changedFiles = [];
  let hasErrors = false;
  let mergeUntilDate = null;
  let futureDatedFiles = [];

  console.log(`Current date: ${currentDate}`);
  console.log('');

  // Initialize Octokit if we have a token
  const token = process.env.GITHUB_TOKEN;
  let octokit = null;
  let owner = null;
  let repo = null;

  if (token && process.env.GITHUB_REPOSITORY) {
    octokit = new Octokit({ auth: token });
    [owner, repo] = process.env.GITHUB_REPOSITORY.split('/');
  }

  // Check for merge-until instruction in PR body
  if (isPullRequest && octokit && owner && repo && process.env.GITHUB_PR_NUMBER) {
    const prBody = await fetchPRBody(octokit, owner, repo, parseInt(process.env.GITHUB_PR_NUMBER));
    mergeUntilDate = parseMergeUntilDate(prBody);
    if (mergeUntilDate) {
      console.log(`Found "DO NOT MERGE until" instruction: ${mergeUntilDate} UTC`);
      console.log('');
    }
  }

  // Get changed files
  if (isPullRequest) {
    try {
      const output = execSync('git diff --name-only origin/main...HEAD', { encoding: 'utf8' });
      changedFiles = output
        .trim()
        .split('\n')
        .filter(file => file.startsWith('lume/src/') && (file.endsWith('.md') || file.endsWith('.mdx')));
    } catch (error) {
      console.error('Error getting changed files:', error.message);
      process.exit(1);
    }
  } else {
    // For pushes to main, get all content files
    try {
      const output = execSync('find lume/src -name "*.md" -o -name "*.mdx"', { encoding: 'utf8' });
      changedFiles = output.trim().split('\n').filter(file => file);
    } catch (error) {
      console.error('Error finding content files:', error.message);
      process.exit(1);
    }
  }

  if (changedFiles.length === 0) {
    console.log('No content files changed');
    return;
  }

  console.log(`Checking ${changedFiles.length} content file(s)...`);
  console.log('');

  for (const file of changedFiles) {
    if (!fs.existsSync(file)) continue;

    try {
      const fileContent = fs.readFileSync(file, 'utf8');
      const { data: frontmatter } = matter(fileContent);

      console.log(`Checking: ${file}`);

      if (!frontmatter.date) {
        console.log('  ⚠️  WARNING: No date field found in frontmatter');
        console.log('');
        continue;
      }

      const postDate = frontmatter.date;

      // Handle different date formats
      let dateToCheck;
      if (typeof postDate === 'string') {
        dateToCheck = postDate.split('T')[0]; // Remove time component if present
      } else if (postDate instanceof Date) {
        dateToCheck = postDate.toISOString().split('T')[0];
      } else {
        console.log(`  ⚠️  WARNING: Invalid date format: ${postDate}`);
        console.log('');
        continue;
      }

      console.log(`  Date: ${dateToCheck}`);

      if (dateToCheck > currentDate) {
        // Check if this is allowed by a merge-until instruction
        if (mergeUntilDate && dateToCheck <= mergeUntilDate) {
          console.log(`  ⏳ Future date allowed by "DO NOT MERGE until ${mergeUntilDate} UTC" instruction`);
          futureDatedFiles.push({ file, date: dateToCheck });
        } else if (mergeUntilDate && dateToCheck > mergeUntilDate) {
          console.log(`  ❌ ERROR: Post date (${dateToCheck}) is after merge-until date (${mergeUntilDate})!`);
          hasErrors = true;
        } else {
          console.log('  ❌ ERROR: Future date detected!');
          hasErrors = true;
        }
      } else {
        console.log('  ✅ OK: Date is not in the future');
      }
      console.log('');

    } catch (error) {
      console.log(`  ❌ ERROR: Failed to parse ${file}: ${error.message}`);
      hasErrors = true;
    }
  }

  // Post a comment if there are future-dated files that are allowed
  if (futureDatedFiles.length > 0 && octokit && owner && repo && process.env.GITHUB_PR_NUMBER) {
    const filesInfo = futureDatedFiles.map(f => `- \`${f.file}\` (date: ${f.date})`).join('\n');
    await postMergeWarningComment(octokit, owner, repo, parseInt(process.env.GITHUB_PR_NUMBER), mergeUntilDate, filesInfo);
  }

  if (hasErrors) {
    console.log('');
    console.log('❌ Validation failed: Found content with future dates!');
    if (mergeUntilDate) {
      console.log(`Some future dates exceed the "DO NOT MERGE until" date (${mergeUntilDate} UTC)`);
    } else {
      console.log(`All content dates must be on or before the current date (${currentDate})`);
      console.log('');
      console.log('To fix this:');
      console.log('1. Update the "date" field in the frontmatter to today\'s date or earlier');
      console.log('2. Use format: date: YYYY-MM-DD');
      console.log('3. Or add "DO NOT MERGE until YYYY-MM-DD UTC" to your PR description');
    }
    process.exit(1);
  } else {
    console.log('✅ All content dates are valid!');
    if (futureDatedFiles.length > 0) {
      console.log(`⚠️  Remember: Do not merge until ${mergeUntilDate} UTC`);
    }
  }
}

// Run the validation
validateContentDates();
