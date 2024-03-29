# API

API quickly runs [JSON-API](https://jsonapi.org/) compliant resources defined by using [JSON-Schema](https://json-schema.org/) attributes definition.

## Badges

[![CircleCI](https://circleci.com/gh/moreandres/api.svg?style=shield)](https://circleci.com/gh/moreandres/api)
[![codecov](https://codecov.io/gh/moreandres/api/branch/master/graph/badge.svg)](https://codecov.io/gh/moreandres/api)
[![Go Report Card](https://goreportcard.com/badge/github.com/moreandres/api)](https://goreportcard.com/report/github.com/moreandres/api)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fmoreandres%2Fapi.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fmoreandres%2Fapi?ref=badge_shield)
[![GoDoc](https://pkg.go.dev/badge/github.com/moreandres/api?status.svg)](https://pkg.go.dev/github.com/moreandres/api?tab=doc)
[![Sourcegraph](https://sourcegraph.com/github.com/moreandres/api/-/badge.svg)](https://sourcegraph.com/github.com/moreandres/api?badge)
[![Open Source Helpers](https://www.codetriage.com/moreandres/api/badges/users.svg)](https://www.codetriage.com/moreandres/api)
[![Release](https://img.shields.io/github/release/moreandres/api.svg?style=flat-square)](https://github.com/moreandres/api/releases)
[![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/moreandres/api)](https://www.tickgit.com/browse?repo=github.com/moreandres/api)

## Usage

```sh
docker build -t my-api .
docker run -p 8080:8080 -it --rm --name my-running-api my-api
```

## Details

API will expose a compliant resource data including:

1. Creating

```sh
POST /v1/resource
```

1. Fetching

```sh
GET /v1/resource
GET /v1/resource?sort=-id&page=0&limit=10
GET /v1/resource/{id}
```

1. Updating

```sh
PATCH /v1/resource/{id}
```

1. Deleting

```sh
DELETE /v1/resource/{id}
```

API will also expose supporting resources for health check and profiling

1. Health

```sh
GET /v1/health
```

1. Profiling

```sh
GET /debug/pprof
```

API adds middleware to:

1. Failover from failures
2. Include service revision
3. Limit concurrent connections
4. Include basic security controls
5. Restrict access
6. Keep usage statistics
7. Avoid CORS issues
8. Support compression
9. Include request ID

Check API tests for further usage details.

## Configuration

API uses API_XXX environment variables or configuration files

1. LogLevel (Info)
2. SchemaName (objects.json)
3. DbUri (file::memory:?cache=shared)
4. HttpPort (8080)
5. UseSSL (false)
6. HttpsPort (8443)
7. CertFile (api.cer)
8. KeyFile (api.key)
9. URL
10. QueryLimit (512)
11. JwtSecret (password)
12. MaxAllowed (20)
13. AccessCidr (0.0.0.0/0)

## Implementation

API uses:

1. [sirupsen/logrus](https://github.com/sirupsen/logrus) for logging.
2. [qri-io/jsonschema](https://github.com/qri-io/jsonschema) for JSON-Schema validation.
3. [cobra](https://github.com/spf13/cobra) for reading configuration from environment and configuration files.
4. [gin](https://github.com/gin-gonic/gin) as web framework to expose resources.
5. [gorm](https://github.com/go-gorm/gorm) for ORM to store and retrieve data.

## Deployment

API can be deployed using classic infrastructure. CloudWatch metrics and logs are enabled. Session Manager can be used to SSH into instances.

```sh
terraform init
terraform apply
terraform destroy
```

## Windows Development

1. Enable [WSL 2](https://docs.microsoft.com/en-us/windows/wsl/install-win10#manual-installation-steps)
2. Install [chocolatey](https://chocolatey.org/install)
3. Install tooling

```sh
choco install git git-lfs github-desktop docker-desktop golang mingw packer terraform vscode python minikube vagrant virtualbox
```
