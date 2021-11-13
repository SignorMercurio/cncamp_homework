# cncamp_homework

<a href="https://github.com/SignorMercurio/cncamp_homework/actions"><img src="https://img.shields.io/github/workflow/status/SignorMercurio/cncamp_homework/Go?logo=GitHub" /></a>
<a href="https://codecov.io/gh/SignorMercurio/cncamp_homework"><img src="https://codecov.io/gh/SignorMercurio/cncamp_homework/branch/main/graph/badge.svg?token=PKWZK3BR9R"/></a>

<details>
<summary><img src="https://img.shields.io/badge/HW01-httpserver-4285f4?logo=google-chrome" /></summary>

## httpserver

A simple HTTP server that you may:

- Access `/header` to find your Request Headers in the Response Headers
- Access `/version` to get the VERSION environment variable
- Access `/log` to write logs in the server
- Access `/healthz` for a health check

### Sample usage

Start a server on `0.0.0.0:8000`:

```shell
$ ./httpserver :8000
```

### Note for Dockerfile

- When using Apple M1 to play with docker, it pulls and builds images for linux/arm/v8 platform by default.
- In order to build images for other platform, you may find [buildx](https://docs.docker.com/buildx/working-with-buildx/) helpful.
- OR, you may also make use of GitHub Actions to avoid the issue.
- When using `alpine` as the base image to run a go binary, `CGO_ENABLED=0` must be set when building due to a different libc implementation on `alpine`. Replacing the dynamic link library also helps.

### Note for Google Cloud Platform

- Running `gcloud --quiet auth configure-docker` requires the service account to have the permission to create bucket. `Storage Admin` role works, but it's clearly not the least
  privilege you can grant.
- You'll need `Kubernetes Engine Developer` role for your service account.
- `secrets.GKE_PROJECT`: GKE's Project ID
- `secrets.GKE_SA_KEY`: Base64 encoded JSON key of your service account

### Things to modify for a different golang app

- Target binary name in `Dockerfile`
- Entry command in `Dockerfile`
- Everywhere `httpserver` appears in `template.yml`
- (Optional) `kustomization.yml` to include other `.yml` representing Kubernetes resources
- `env` in `.github/workflows/gke.yml`
- `secrets.GKE_PROJECT` and `secrets.GKE_SA_KEY` in `.github/workflows/gke.yml`

</details>

<details>
<summary><img src="https://img.shields.io/badge/HW02-Docker-2496ed?logo=docker" /></summary>

## Docker

Build a multi-stage docker image for httpserver.

> See [Dockerfile](Dockerfile).

</details>

<details>
<summary><img src="https://img.shields.io/badge/HW03-Kubernetes-326ce5?logo=kubernetes" /></summary>

## Kubernetes

Deploy httpserver on Kubernetes. Based on the first homework, I choose to deploy it on Google Kubernetes Engine.

### Changes in httpserver

- Deprecate `valyala/fasthttp`, use `net/http` and `gorilla/mux`
- Add unit tests, coverage 100%
- CI / CD with GitHub Actions
  - CI: Codecov
  - CD: Deploy to GKE
- Add graceful termination when receiving SIGTERM

### Features \| 🚧 In construction

- [x] Add resource limit and request
- [x] Add readiness probe
- [x] <del>Add liveness probe</del> No need for liveness probe
- [x] Add graceful termination (see httpserver source code)
- [ ] Add graceful initialization with postStart
- [ ] Add configurations with ConfigMap
- [ ] Add support for leveled logging

</details>
