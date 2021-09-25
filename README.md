# cncamp_homework

## httpserver

A simple HTTP server that you may:

- Access `/header` to find your Request Headers in the Response Headers
- Access `/version` to get the VERSION environment variable
- Access `/log` to write logs in the server
- Access `/healthz` for a healch check

### Note for Dockerfile

- When using Apple M1 to play with docker, it pulls and builds images for linux/arm/v8 platform by default.
- In order to build images for other platform, you may find [buildx](https://docs.docker.com/buildx/working-with-buildx/) helpful.
- OR, you may also make use of GitHub Actions to avoid the issue.
- When using `alpine` as the base image to run a go binary, `CGO_ENABLED=0` must be set when building due to a different libc implementation on `alpine`. Replacing the dynamic link library also helps.

### Note for Google Cloud Platform

- Running `gcloud --quiet auth configure-docker` requires the service account to have the permission to create bucket. `Storage Admin` role works, but it's clearly not the least
  privilege you can grant.
- You'll need `Kubernetes Engine Developer` role for your service account.

### Things to modify for a different golang app

- Target binary name in `Dockerfile`
- Entry command in `Dockerfile`
- Everywhere `httpserver` appears in `deployment.yml`
- (Optional) A `service.yml` when things get complicated
- (Optional) `kustomization.yml` to include other `.yml` representing Kubernetes resources
- `env` in `.github/workflows/gke.yml`
- `secrets.GKE_PROJECT` and `secrets.GKE_SA_KEY` in `.github/workflows/gke.yml`
