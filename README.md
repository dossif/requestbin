# RequestBin

A minimal HTTP debug/echo server written in Go. It reflects back everything about
the incoming request (method, headers, body, remote address, etc.) as JSON, and
optionally lets you force an arbitrary HTTP response status via the URL path —
useful for testing webhooks, HTTP clients, and API integrations.

Docker image: https://hub.docker.com/r/dossif/requestbin

## Requirements

- Go 1.27+ (only needed to build from source; uses `net/http.ServeMux` path
  parameters and `log/slog`, both requiring a recent Go toolchain)

No external dependencies — the entire server is built on the standard library.

## Compile

    go build -o requestbin ./cmd/requestbin

## Run

    env RB_LISTEN="localhost:8080" ./requestbin

### Environment variables

| Variable    | Default       | Description                          |
|-------------|---------------|---------------------------------------|
| `RB_LISTEN` | `0.0.0.0:8080`| Address and port the server listens on |

## Docker compose

    docker-compose -f docker-compose.yml up

## Usage

Any request returns a JSON document describing the request and response:

    $ curl http://127.0.0.1:8080/status/201?foo=bar
    {
      "Request": {
        "RemoteAddr": "127.0.0.1:56035",
        "Method": "GET",
        "Host": "127.0.0.1:8080",
        "Proto": "HTTP/1.1",
        "ProtoMajor": 1,
        "ProtoMinor": 1,
        "Pattern": "/status/{status}",
        "Url": "/status/201?foo=bar",
        "Path": "/status/201",
        "RawQuery": "foo=bar",
        "Query": {
          "foo": [
            "bar"
          ]
        },
        "RequestURI": "/status/201?foo=bar",
        "ContentLength": 0,
        "Trailer": null,
        "Body": "",
        "Headers": {
          "Accept": [
            "*/*"
          ],
          "User-Agent": [
            "curl/7.79.1"
          ]
        },
        "Cookies": {}
      },
      "Response": {
        "Status": 201,
        "Headers": {
          "Content-Type": [
            "application/json"
          ]
        }
      }
    }

Use `/status/{code}` to force a specific HTTP response status (any integer in
the range `200`-`599`):

    $ curl http://127.0.0.1:8080/status/512 -I
    HTTP/1.1 512 status code 512
    Content-Type: application/json
    Date: Tue, 02 Aug 2022 09:48:02 GMT
    Content-Length: 490

An invalid or out-of-range status returns a `599` JSON error response.

## Release process

Pushing a `v*.*.*` git tag triggers `.github/workflows/build.yml`, which
compiles the binary, publishes a Docker image to `dossif/requestbin`, and
creates a GitHub release with the binary attached.
