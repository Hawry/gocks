# GoMOCKA

[![Go Report Card](https://goreportcard.com/badge/github.com/hawry/gomocka)](https://goreportcard.com/report/github.com/hawry/gomocka) [![GoDoc](https://godoc.org/github.com/Hawry/gomocka?status.svg)](https://godoc.org/github.com/Hawry/gomocka)

A very lightweight and simple mocking service. The idea is that a JSON file configures possible endpoints, and the service is either run locally or in a docker container.

<!-- TOC depthFrom:2 -->

- [Getting started](#getting-started)
  - [gomocka --help](#gomocka---help)
- [Configuration](#configuration)
  - [Port](#port)
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
      --help                    Show context-sensitive help (also try --help-long and --help-man).
      --version                 print version of gock
  -c, --config="settings.json"  configuration file to create endpoints from
  -v, --verbose                 enabled verbose logging. if --silent is used, --verbose will be ignore
  -s, --silent                  disabled all output except for errors. overrides --verbose if set
  -g, --generate                generate a sample configuration
```

## Configuration

```json
{
  "port": 8080,
  "mocks": [
    {
      "path":"/",
      "method":"GET",
      "response_code": 200,
      "response_body": "{\"hello\":\"world\"}",
      "headers": {
        "Content-Type":"application/json",
        "x-trace-id": "thisisatrace"
      }
    },
    {
      "path":"/hello",
      "method":"POST",
      "response_code": 400,
      "response_body":"bad something"
    }
  ]
}
```

### Port
Port to listen to. The service automatically binds to `0.0.0.0:<port>`.

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

## Build from source

Prerequisite: [Go][2], recommended minimum version is 1.10.

Clone this repository and run `make static` to create statically linked binary without CGO. The created binary can then be used as part of the docker build process mentioned above.

## Run in Docker
Before you try to build the docker image, make sure you have compiled a statically linked binary as explained above.

The project contains a lightweight Dockerfile, which only should be about 12MB in size. To create the docker image; just run `docker build -t gomocka .` in the project directory. Modify `settings.json` before building to make sure the latest configuration is part of the docker image.

To run the docker image: `docker run -d -p 8080:8080 gomocka`. Replace `8080` with the port number that you've specified in your `settings.json` file.

## Roadmap
- [ ] TLS support
- [ ] Load configuration from dynamic location without a build step
- [ ] Handle authorization in mocked endpoints
- [ ] Bind to specific network interface/address
- [ ] Copy request data to response data (such as headers, etc.)


[1]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods
[2]: https://www.golang.org
