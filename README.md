# GoMOCKA

[![Go Report Card](https://goreportcard.com/badge/github.com/hawry/gomocka)](https://goreportcard.com/report/github.com/hawry/gomocka) [![GoDoc](https://godoc.org/github.com/Hawry/gomocka?status.svg)](https://godoc.org/github.com/Hawry/gomocka) [![Build Status](https://travis-ci.org/hawry/gomocka.svg?branch=master)](https://travis-ci.org/hawry/gomocka)

A very lightweight and simple mocking service. The idea is that a JSON file configures possible endpoints, and the service is either run locally or in a docker container.

<!-- TOC depthFrom:2 -->

- [Getting started](#getting-started)
  - [gomocka --help](#gomocka---help)
- [Configuration](#configuration)
  - [Port](#port)
  - [Authorization](#authorization)
    - [OpenID](#openid)
    - [Basic Auth](#basic-auth)
    - [Header](#header)
  - [Mocks](#mocks)
    - [Path](#path)
    - [Method](#method)
    - [Response Code](#response-code)
    - [Response Body](#response-body)
    - [Headers](#headers)
- [Examples](#examples)
  - [Hard coded response](#hard-coded-response)
  - [Dynamic response](#dynamic-response)
  - [With response headers](#with-response-headers)
  - [Authorization: bearer token (live endpoint)](#authorization-bearer-token-live-endpoint)
  - [Authorization: bearer token (static)](#authorization-bearer-token-static)
  - [Authorization: custom header](#authorization-custom-header)
  - [Authorization: basic auth](#authorization-basic-auth)
- [Build from source](#build-from-source)
- [Run in Docker](#run-in-docker)
- [Roadmap](#roadmap)

<!-- /TOC -->

## Getting started

Download or clone this repository and build the code (requires [Go][2]). Create a settings file according to the configuration part of this documentation, or create a sample configuration file by running `gomocka -g`.

Start the mock server: `gomocka`. By default, the application will look for a configuration file named `settings.json` in the current working dir. To change configuration file use the `--config` flag: `gomocka --config=./example.json`.

### gomocka --help

```
usage: gomocka [<flags>]

Flags:
      --help                    Show context-sensitive help (also try
                                --help-long and --help-man).
      --version                 print version of gock
  -c, --config="settings.json"  configuration file to create endpoints from
  -v, --verbose                 enabled verbose logging. if --silent is used,
                                --verbose will be ignore
  -s, --silent                  disabled all output except for errors. overrides
                                --verbose if set
  -g, --generate                generate a sample configuration
  -d, --gendocker               generates a docker file - specify the config
                                file to add with the flag --config
```

## Configuration

```json
{
  "port": 8080,
  "authorization": {
    "openid": {
      "jwks": "https://www.anopenidprovider.com/.well-known/keys.json"
    },
    "basic_auth": {
      "username": "ausername",
      "password": "apassword"
    },
    "header": {
      "Authorization": "Bearer thisisatoken"
    }
  },
  "mocks": [
    {
      "path": "/",
      "method": "GET",
      "response_code": 200,
      "response_body": "{\"hello\":\"world\"}",
      "headers": {
        "Content-Type": "application/json",
        "x-trace-id": "thisisatrace"
      }
    },
    {
      "path": "/hello",
      "method": "POST",
      "response_code": 400,
      "response_body": "bad something"
    },
    {
      "path": "/health",
      "method": "GET",
      "response_code": 200,
      "disable_auth": true
    }
  ]
}
```

### Port

Port to listen to. The service automatically binds to `0.0.0.0:<port>`.

### Authorization

Authorization configurations required for all paths. Use either only one of the possible configurations, or all of them.

#### OpenID

Specify an OpenID endpoint of which to validate incoming bearer tokens against. The tokens must be signed by one of the keys provided at the endpoint - but no further validations are made.

#### Basic Auth

Basic authentication data in the form of `username` and `password`.

#### Header

Header key and value, which both can be completely customized.

### Mocks

Array of what paths to respond to.

#### Path

Path to listen on. Relative to the root, to access `http://example.com/resource/action`, this attribute should contain `/resource/action`.

#### Method

One of the standard [HTTP Methods][1].

#### Response Code

Response code to return. The response code will be sent regardless of the payload body.

#### Response Body

String representation of the response body. The response body will be sent regardless of header value. To send an empty body, keep this attribute as a blank string: `"response_body":""`.

#### Disable Auth

Boolean which can disable the authorization settings in the Authorization part of the configuration. A path with the `disable_auth` set to `true` will ignore the authorization configuration and respond if it's matched to a request. Default value is `false`.

#### Headers

Headers to set in the response. You can use any string as a header key, and any string as a header value.

## Examples

### Hard coded response

```json
"mocks": [
  {
    "path":"/resource",
    "method": "GET",
    "response_code": 200,
    "response_body": "this is a hard coded response"
  }
]
```

```
$ curl localhost:8080/resource -i

HTTP/1.1 200 OK
Date: Sat, 02 Mar 2019 13:12:41 GMT
Content-Length: 29
Content-Type: text/plain; charset=utf-8

this is a hard coded response
```

### Dynamic response

```json
"mocks": [
  {
    "path":"/users/{userid}",
    "method":"GET",
    "response_code": 200,
    "response_body": "your id is {userid}"
  }
]
```

```
$ curl localhost:8080/users/1337 -i

HTTP/1.1 200 OK
Date: Sat, 02 Mar 2019 13:14:34 GMT
Content-Length: 15
Content-Type: text/plain; charset=utf-8

your id is 1337
```

### With response headers

```json
"mocks": [
  {
    "path":"/withheaders",
    "method":"GET",
    "response_code": 200,
    "response_body": "a hardcoded response",
    "headers": {
      "Content-Type": "application/json",
      "X-Custom-Header": "custom-header"
    }
  }
]
```

```
$ curl localhosi:8080/withheaders -i

HTTP/1.1 200 OK
Content-Type: application/json
X-Custom-Header: custom-header
Date: Sat, 02 Mar 2019 13:15:52 GMT
Content-Length: 20

a hardcoded response
```

### Authorization: bearer token (live endpoint)

```json
{
  "openid": {
    "jwks": "https://provider.com/.well-known/keys.json"
  }
}
```

```
$ curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U" localhost:8080/resource -i

HTTP/1.1 200 OK
Date: Tue, 12 Nov 2019 20:42:56 GMT
Content-Length: 29
Content-Type: text/plain; charset=utf-8

this is a hard coded response
```

### Authorization: bearer token (static)

```json
{
  "authorization": {
    "header": {
      "Authorization": "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImp0aSI6IjZhZjcwZTkxLWE2OGMtNDk3Ny1iZjRkLTYyYTgzNWJlZTRhMCIsImlhdCI6MTU1MTYzODMzMSwiZXhwIjoxNTUxNjQxOTMxfQ.6vo3Jgsrac7cn3V-RUNWWeTPPQFmpWJXhyNoRIp-FyE"
    }
  }
}
```

```
$ curl -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImp0aSI6IjZhZjcwZTkxLWE2OGMtNDk3Ny1iZjRkLTYyYTgzNWJlZTRhMCIsImlhdCI6MTU1MTYzODMzMSwiZXhwIjoxNTUxNjQxOTMxfQ.6vo3Jgsrac7cn3V-RUNWWeTPPQFmpWJXhyNoRIp-FyE" localhost:8080/resource -i
```

### Authorization: custom header

```json
{
  "authorization": {
    "header": {
      "x-api-key": "GD+6fJCdCYObdZt4oK+yvK/rsnY2LFUxNayBYxDUu34="
    }
  }
}
```

```
$ curl -H "x-api-key: GD+6fJCdCYObdZt4oK+yvK/rsnY2LFUxNayBYxDUu34=" localhost:8080/resource -i
```

### Authorization: basic auth

```json
{
  "authorization": {
    "basic_auth": {
      "username": "user",
      "password": "pass"
    }
  }
}
```

```
$ curl --user user:pass localhost:8080/resource -i
```

## Build from source

Prerequisite: [Go][2], recommended minimum version is 1.10.

Clone this repository and run `make static` to create statically linked binary without CGO. The created binary can then be used as part of the docker build process mentioned above.

## Run in Docker

Before you try to build the docker image, make sure you have compiled a statically linked binary as explained above.

To generate the dockerfile just run `gomocka --gendocker --config settings.json` which will create a default dockerfile in the project root. If you wish to use another settings file, just specify that with the `--config` flag. The dockerfile will bundle the settings-file in the docker image, and if you make any changes you'll need to rebuild the docker image.

To run the docker image: `docker run -d -p 8080:8080 gomocka`. Replace `8080` with the port number that you've specified in your `settings.json` file.

## Roadmap

- [ ] TLS support
- [ ] Load configuration from dynamic location without a build step
- [x] Handle authorization in mocked endpoints
- [ ] Bind to specific network interface/address
- [ ] Copy request data to response data (such as headers, etc.)

[1]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods
[2]: https://www.golang.org
