# Plan: erofs-backed serving + content-addressed preview site

> **Status:** implemented on branch `feat/erofs-serving-preview` (uncommitted).
> This doc has been reconciled with the final implementation; notable deviations
> from the original plan (non-reproducible volumes, sub-block pre-size fix,
> server-authoritative hashing, dev-mode branch discovery) are called out inline.

## Context

Today `internal/lume` builds the site with Deno/Lume into `lume/_site`, zips it
to `var/site.zip`, and serves the live site by swapping `f.fs = os.DirFS(destDir)`
(`internal/lume/lume.go:374`). A separate external "FutureSight" service receives
`site.zip` via a multipart POST to `{FutureSightURL}/upload` (`lume.go:651`) for
previews. There is no content addressing — only the git SHA is tracked.

We want two things:

1. **Build and serve xesite assets through `github.com/Xe/erofs`** (already an
   indirect dep at v0.4.0; pure Go, no external tools). The built site becomes a
   single immutable EROFS volume that `internal/lume.FS` serves via `fs.FS`.
2. **A content-addressed preview site**: a new standalone binary stores each
   build's erofs volume in Tigris (`github.com/tigrisdata/storage-go`) and serves
   each version on its own subdomain — a **permanent** host keyed by the truncated
   sha256 of the volume, plus a **moving** host keyed by a DNS-safe slug of the
   git branch.

`erofs.FS` implements `fs.FS`, so it drops into the existing `f.fs` slot with no
change to the HTTP serving path (`http.FileServerFS(fs)` in `cmd/xesite/main.go:86`).

## Relevant APIs (verified via `go doc`)

**erofs** (`github.com/Xe/erofs`):

- `NewBuilder(w io.WriterAt, opts ...BuildOption) *Builder`
- `(*Builder).AddFromFS(fsys fs.FS) error` — walks an `fs.FS`, adds every entry
- `(*Builder).Build() error` — finalizes
- `WithEpoch(t time.Time)`, `WithCompression(alg)` (`CompressionAutoLZ4` / `CompressionNone`)
- `Open(r io.ReaderAt) (*FS, error)` — returns an `fs.FS` (also `StatFS`)

**storage-go** (`github.com/tigrisdata/storage-go`): `type Client struct{ *s3.Client }`

- `New(ctx, ...Option) (*Client, error)`; `WithAccessKeypair(id, secret)`,
  `WithGlobalEndpoint()` (default endpoint `https://t3.storage.dev`). `New` calls
  `awsConfig.LoadDefaultConfig`, so it also picks up `AWS_ACCESS_KEY_ID` /
  `AWS_SECRET_ACCESS_KEY` from the environment.
- Because it embeds `*s3.Client`, use standard `PutObject`, `GetObject`,
  `HeadObject`, `ListObjectsV2`. `GetObject().Body` is an `io.ReadCloser` (NOT an
  `io.ReaderAt`), so volumes are read fully into a `bytes.Reader` before
  `erofs.Open`.

---

## Part A — Build & serve via erofs (`internal/lume`)

### A1. New file `internal/lume/erofs.go`

- `func buildEROFS(srcDir, destPath string, epoch time.Time) (hash string, err error)`:
  1. `fout, err := os.Create(destPath)` (an `*os.File` satisfies `io.WriterAt`),
     then `fout.Truncate(4096)` to pre-size to one block. **Why:** `Build()` reads
     the whole first block back to checksum the superblock, and `*os.File.ReadAt`
     returns `io.EOF` on a sub-block-sized image. The real site is megabytes so this
     only matters for tiny inputs (e.g. tests), and EROFS images are block-aligned
     anyway, so the padding is part of a valid image. (The library's own tests use a
     zero-padding in-memory buffer and never hit this.)
  2. `b := erofs.NewBuilder(fout, erofs.WithEpoch(epoch), erofs.WithCompression(erofs.CompressionAutoLZ4))`.
  3. `b.AddFromFS(os.DirFS(srcDir))`, then `b.Build()`, then close `fout`.
  4. Reopen the finished file and stream it through `sha256.New()`, returning the
     first `hashLen` hex chars (`const hashLen = 16`, 64 bits — short enough for a
     DNS label, wide enough to avoid collisions across builds). Hash the **volume
     bytes**, per the requirement.
- `type erofsFS struct { *erofs.FS; f *os.File }` with `func (e *erofsFS) Close() error { return e.f.Close() }`.
  The existing swap logic already closes the old `f.fs` if it is an `io.Closer`
  (`lume.go:94`, `lume.go:368`), so wrapping the file handle here makes the handle
  lifecycle correct without touching that logic.
- `func openEROFS(path string) (*erofsFS, error)`: `os.Open` → `erofs.Open(f)` →
  wrap. The `*os.File` is the `io.ReaderAt`; erofs reads lazily so the handle stays
  open for the FS lifetime.

### A2. Wire into `build()` (`internal/lume/lume.go:333`)

Replace the tail of `build()` (currently lines 362–374, the `ZipFolder` + `os.DirFS`
swap):

- Build the volume to `filepath.Join(f.opt.DataDir, "site.erofs")` using `buildEROFS`,
  passing `begin` (the build start time) as the epoch (inode mtimes). Note: the
  volume is **not** byte-reproducible — the EROFS superblock embeds `time.Now()`
  and the builder lays out directory entries in map order — so each build gets its
  own content address. Open the new volume *before* closing the old `f.fs` so a
  failed open keeps the previous build serving.
- Store the returned hash on the struct: add `volumeHash string` to `FS`
  (`lume.go:74`), guarded by the existing `f.lock`.
- Close the old `f.fs` (keep existing `io.Closer` block), then
  `f.fs = openEROFS(volumePath)`.
- **Zip pipeline:** the user chose "replace zip for serving." Keep `ZipFolder`/
  `site.zip` generation only if the internal `/zip` endpoint
  (`cmd/xesite/internalapi.go:34`) must keep working; otherwise drop it. Recommended:
  keep `ZipFolder` writing `site.zip` for that download endpoint (cheap, isolated),
  but remove zip from the serving path. `zip.go`'s `ZipServer` type is already
  unused and can be deleted.

### A3. Expose metadata for the preview upload

- `func (f *FS) VolumeHash() string` (lock-guarded getter) — mirrors `Commit()`.
- `func (f *FS) Branch() string` returns `f.opt.Branch`.

### A3a. Dev-mode branch discovery

In dev mode the build reads the working tree (`fs.repoDir = os.Getwd()`) but
`f.opt.Branch` would otherwise stay pinned to the `--git-branch` flag (default
`main`), so a developer on a feature branch would publish previews under the
`main` slug. In the `o.Development` block of `New()`, open the working tree with
go-git and adopt the checked-out branch:

```go
if wtRepo, err := git.PlainOpenWithOptions(fs.repoDir, &git.PlainOpenOptions{DetectDotGit: true}); err == nil {
    if head, err := wtRepo.Head(); err == nil && head.Name().IsBranch() {
        o.Branch = head.Name().Short()
    }
}
```

Production is unchanged: outside dev mode the branch still comes from the flag.

### A4. Tests — `internal/lume/erofs_test.go` (table-driven, per skill)

- `TestBuildEROFS`: write a small tree into `t.TempDir()` (text + binary files, a
  subdir), `buildEROFS`, reopen with `openEROFS`, assert each file round-trips via
  `fs.ReadFile`, and that the hash is `hashLen` long.
- `TestBuildEROFSDistinctContent`: different content yields different hashes.
  (Volumes are **not** byte-reproducible across builds — see A2 — so determinism
  cannot be asserted; each build legitimately gets its own permanent URL.)
- `TestServeEROFS`: build a volume, serve it through `httptest` +
  `http.FileServerFS`, and assert `GET /`, `/index.html`, and a nested path return
  the right bytes — this is the actual "serve erofs over HTTP" path. Inner subtests
  must **not** call `t.Parallel()` (the parent's `defer srv.Close()` would fire
  first).
- `TestOpenEROFSMissing`: opening a missing volume errors.

---

## Part B — New preview binary `cmd/futuresight`

Reuses the existing `--future-sight-url` flag name and the FutureSight upload
contract, but now in-repo, Tigris-backed, and serving content-addressed subdomains.
(Name kept for continuity; `cmd/xesite-preview` is an alternative.)

### B1. CLI skeleton (`cmd/futuresight/main.go`) — follows xe-go style

`flagenv.Parse()` then `flag.Parse()`, `internal.Slog()`, kebab-case flags:

- `--bind` (default `:3000`), `--internal-api-bind` (`:3001`, for `/healthz`)
- `--tigris-bucket` (preview volume bucket)
- `--base-domain` (e.g. `preview.xeiaso.net`) — used to strip the host label
- `--cache-dir` (default `./var/volumes`) — local disk cache of fetched erofs
  volumes; backed by a Kubernetes `emptyDir` mount in prod
- `--upload-token` (bearer/HMAC shared secret guarding `/upload`)
- Tigris creds from `AWS_ACCESS_KEY_ID` / `AWS_SECRET_ACCESS_KEY` env (already the
  Tigris pattern in this repo) or explicit flags passed to `storage.WithAccessKeypair`.

Construct one `*storage.Client` via `storage.New(ctx, storage.WithGlobalEndpoint(), ...)`.

### B2. Object layout in Tigris

- `volumes/<hash>.erofs` — immutable, content-addressed.
- `branches/<branch-slug>` — tiny object whose body is the current `<hash>` for that
  branch (the moving pointer). Updated on each upload.

### B3. `POST /upload` (admin host / token-guarded)

Mirror today's multipart contract (`lume.go`) plus metadata. Accept multipart
form fields: `file` (the erofs volume), `branch`, `hash`.

- Verify the `Authorization: Bearer <token>` header against `--upload-token`
  (constant-time compare; empty token disables auth for local dev).
- Stream `file` to a temp file in the cache dir while computing sha256 — the
  **server is authoritative on the content hash**, so the client's `hash` field is
  only advisory. `PutObject(volumes/<hash>.erofs)`
  (`Content-Type: application/octet-stream`), then `os.Rename` the temp file into
  the cache to prime it (the first serve skips the download round-trip).
- Slugify `branch` → `PutObject(branches/<slug>)` with body `<hash>`.
- Respond `200` with the hash so the client can log the resulting URL.

**Upload host:** `api.preview.xeiaso.net`. `handleUpload` is host-agnostic
(routed by path), so the existing `*.preview.xeiaso.net` wildcard ingress + cert
already cover it — no extra ingress rule or cert. The endpoint is public and
token-guarded; it caps the body via `http.MaxBytesReader` (`--max-upload-bytes`,
default 512 MiB → `413` on overflow). Dev boxes set
`FUTURE_SIGHT_URL=https://api.preview.xeiaso.net`; in-cluster xesite keeps using
`http://futuresight.default.svc`.

### B4. Host-based serving (`r.Host`)

No host-routing exists today (only `internal/domain_redirect.go:31`), so add a
small handler that:

1. Strips `--base-domain` suffix from `r.Host` to get the leftmost label
   (`<label>.preview.xeiaso.net` → `<label>`).
2. If `label` is `hashLen` hex chars → treat as a content hash directly.
   Else treat as a branch slug → `GetObject(branches/<label>)` to resolve `<hash>`.
3. Serve the volume for `<hash>` via `http.FileServerFS(volumeFS)`.

`resolveVolume(ctx, hash) (fs.FS, error)` — **disk-backed cache** in `--cache-dir`:

- Cache path `filepath.Join(cacheDir, hash+".erofs")`. In-process
  `map[string]*erofs.FS` (guarding the open handles) under a `sync.RWMutex`;
  content-addressed volumes are immutable, so entries never expire.
- On miss: if the cache file is absent, `GetObject(volumes/<hash>.erofs)` and stream
  `.Body` to the cache file (download to a `<hash>.tmp` then atomic `os.Rename`, so a
  crash mid-download can't leave a truncated volume). Then `os.Open` the cache file
  (the `*os.File` is the `io.ReaderAt`) → `erofs.Open` → store in the map. Because
  erofs reads lazily from the file, volumes are never loaded fully into memory.
- Branch pointers are re-resolved per request (one small `GetObject(branches/<slug>)`),
  so a new build is visible on the branch host immediately; the resolved `<hash>`
  then hits the disk cache.

### B5. Slugify helper + tests (`cmd/futuresight/slug.go`, `slug_test.go`)

- `func slugifyBranch(s string) string`: lowercase, replace any run of non
  `[a-z0-9]` with a single `-`, trim leading/trailing `-`, truncate to 63 chars
  (DNS label limit). E.g. `feat/Foo_Bar` → `feat-foo-bar`.
- Table-driven test covering: simple branch, slashes, uppercase, leading/trailing
  junk, over-length truncation, unicode → ascii fallback. `t.Parallel()`, map or
  slice of cases with `name`, `want`.

### B6. Replace `futureSight()` in `internal/lume/lume.go`

Change the POST body from `site.zip` to `site.erofs`, and write the `branch` and
`hash` form fields (from `f.opt.Branch` — now working-tree-aware in dev mode, see
A3a — and `f.VolumeHash()`) **before** the `file` part so a streaming reader sees
them first. Send `Authorization: Bearer <token>` from a new
`FutureSightToken` option, wired from a new xesite `--future-sight-token` flag
(`cmd/xesite/main.go`). Keep the same `{FutureSightURL}/upload` endpoint and
multipart shape so the call site (`go f.FutureSight(...)`) is unchanged; existing
`futureSightPokes`/`futureSightErrors` metrics stay valid.

---

## Part C — Deployment

Model on the existing `manifest/xesite/*` and `manifest/sponsor-panel/*` (nginx
ingress + cert-manager + 1Password secrets).

- New `manifest/futuresight/`: `deployment.yaml`, `service.yaml`,
  `kustomization.yaml`, `1password.yaml` (Tigris keypair + upload token — reuse the
  `xesite-anubis-tigris` OnePasswordItem pattern at `manifest/xesite/1password.yaml:18`).
  Mount an `emptyDir` volume at the container's `--cache-dir` (mirror the existing
  xesite `emptyDir` data mount in `manifest/xesite/deployment.yaml`). The cache is
  pure scratch — losing it on pod restart just triggers a re-download from Tigris.
- `ingress.yaml` with a **wildcard host** `*.preview.xeiaso.net`. Note: nginx
  supports wildcard hosts, but cert-manager needs a **DNS-01** issuer for a
  wildcard TLS cert (the existing `letsencrypt-prod` issuer uses HTTP-01, which
  can't issue wildcards). This is the one real infra prerequisite — flag it for the
  user; a DNS-01 ClusterIssuer (or a pre-provisioned wildcard cert) must exist.
- Register the binary in the Docker build / Tekton pipeline alongside the other
  `cmd/*` binaries (`docker/`, `.tekton/xe-site-pipeline.yaml`).
- **Production xesite does not set `--future-sight-url`** — previews are a
  development feature, so the in-cluster deployment leaves it empty and the
  `if o.FutureSightURL != ""` guard skips the upload entirely (no token needed in
  prod). Only dev boxes set `FUTURE_SIGHT_URL` (→ `https://api.preview.xeiaso.net`)
  and `FUTURE_SIGHT_TOKEN` in their local `.env`.

---

## Verification (done)

- `go build ./...`, `go vet`, `gofmt` clean. `go mod tidy` promoted `Xe/erofs` and
  `tigrisdata/storage-go` to direct deps.
- New table-driven tests pass: erofs build/open/**HTTP-serve** round-trips
  (`internal/lume`), slugify, subdomain parsing, hash detection (`cmd/futuresight`).
- **End-to-end (lume):** `go test ./internal/lume -run TestCanBuildSite` built the
  real 782-file site into `site.erofs` (`hash=9fdda79a858ceadb`), opened it, and
  served it — confirming the full Deno → erofs → `fs.FS` path.
- **Preview binary (manual, needs a test Tigris bucket):** run `cmd/futuresight`,
  `POST /upload` a built `site.erofs` with `branch=main`, then
  `curl -H 'Host: <hash>.preview.localhost' …` and `-H 'Host: main.preview.localhost' …`
  and confirm both serve the same content. The Tigris-backed `Store` is not unit
  tested (no local creds); the serving/routing logic around it is.

## Open prerequisite for the user

Wildcard TLS for `*.preview.xeiaso.net` requires a DNS-01 cert-manager issuer (the
current `letsencrypt-prod` is HTTP-01). Confirm whether such an issuer exists or
should be added.
