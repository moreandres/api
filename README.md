# API

[![Build Status](https://travis-ci.org/moreandres/api.svg)](https://travis-ci.org/moreandres/api)
[![codecov](https://codecov.io/gh/moreandres/api/branch/master/graph/badge.svg)](https://codecov.io/gh/moreandres/api)
[![Go Report Card](https://goreportcard.com/badge/github.com/moreandres/api)](https://goreportcard.com/report/github.com/moreandres/api)
[![GoDoc](https://pkg.go.dev/badge/github.com/moreandres/api?status.svg)](https://pkg.go.dev/github.com/moreandres/api?tab=doc)
[![Sourcegraph](https://sourcegraph.com/github.com/moreandres/api/-/badge.svg)](https://sourcegraph.com/github.com/moreandres/api?badge)
[![Open Source Helpers](https://www.codetriage.com/moreandres/api/badges/users.svg)](https://www.codetriage.com/moreandres/api)
[![Release](https://img.shields.io/github/release/moreandres/api.svg?style=flat-square)](https://github.com/moreandres/api/releases)
[![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/moreandres/api)](https://www.tickgit.com/browse?repo=github.com/moreandres/api)

## What?

Run [JSON-API](https://jsonapi.org/) compliant APIs defined using attributes [JSON-Schema](https://json-schema.org/) definitions.

## Why?

Quickly run compliant APIs without coding.

## How?

1. [sirupsen/logrus](https://github.com/sirupsen/logrus) for logging.

1. [qri-io/jsonschema](https://github.com/qri-io/jsonschema) for json schema.

1. [viper](https://github.com/spf13/viper) for command line handling.

1. [cobra](https://github.com/spf13/cobra) for reading configuration.

1. [gin](https://github.com/gin-gonic/gin) as web framework.

1. [gorm](https://github.com/go-gorm/gorm) for ORM.

## Example

Given `persons.json` defined as
 ```
{
    "type": "object",
    "properties": {
        "first_name": { "type": "string" },
        "last_name": { "type": "string" },
        "birthday": { "type": "string", "format": "date" },
        "address": {
            "type": "object",
            "properties": {
                "street_address": { "type": "string" },
                "city": { "type": "string" },
                "state": { "type": "string" },
                "country": { "type" : "string" }
            }
        }
    }
}
```

The running `api run --schema persons.json` will expose a compliant CRUD which can be exercised as below.

1. GET /persons -> 2001, []

1. GET /persons/{uuid} -> 404

1. POST /persons { a } -> 201, [ { a } ] 

1. POST /persons { b } -> 201, [ { b } ]

1. GET /persons?filter=a&sort=a&page=0&limit=100 -> 200, [ a, b ] 

1. GET /persons -> 200, [ a, b ]

1. GET /persons/{uuid} -> 200

1. DELETE /persons/{uuid} -> 200

Additional middlewares are enabled as well.

1. Basic Auth using JWT tokens

1. OPTIONS

1. Health probes /liveness /readiness /startup

1. Logging

1. Recovery 

1. CORS and CSRF middlewares are enabled

1. Transaction/Correlation/Revision ID

1. Rate-limits per Client

1. Stats middleware

*Note*: Client headers should contain `Content-Type: application/vnd.api+json`

## Design

Config (viper) -> Routes (gin) -> JSON Schema (jsonschema) -> Datbabasabase (Ggrgorm)

## Collaborate

sudo apt update
sudo apt upgrade
sudo apt install golang golang-golang-x-tools