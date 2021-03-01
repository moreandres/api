# API

API quickly runs [JSON-API](https://jsonapi.org/) compliant resources defined by using [JSON-Schema](https://json-schema.org/) attributes definition.

## Badges

[![Build Status](https://travis-ci.org/moreandres/api.svg)](https://travis-ci.org/moreandres/api)
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
$ docker build -t my-api .
$ docker run -it --rm --name my-running-api my-api
```

## Details

API will expose a compliant resource data including:

1. Creating
```
POST /v1/resource
```

2. Fetching
```
GET /v1/resource
GET /v1/resource?sort=-id&page=0&limit=10
GET /v1/resource/{id}
```

3. Updating
```
PATCH /v1/resource/{id}
```

4. Deleting
```
DELETE /v1/resource/{id}
```

API will also expose supporting resources for health check and profiling

1. Health
```
GET /v1/health
```

2. Profiling
```
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

## Implementation

API uses:

1. [sirupsen/logrus](https://github.com/sirupsen/logrus) for logging.

2. [qri-io/jsonschema](https://github.com/qri-io/jsonschema) for JSON-Schema validation.

3. [cobra](https://github.com/spf13/cobra) for reading configuration from environment and configuration files.

4. [gin](https://github.com/gin-gonic/gin) as web framework to expose resources.

5. [gorm](https://github.com/go-gorm/gorm) for ORM to store and retrieve data.
