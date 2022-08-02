# RequestBin

Build:

    go build -o requestbin ./src/main

Run:

    env RB_LISTEN="localhost:8080" ./requestbin

Docker compose:

    docker-compose -f docker-compose.yml up --build

Usage:

    $ curl http://127.0.0.1:8080/
    {
      "Request": {
        "RemoteAddr": "127.0.0.1:56035",
        "Method": "GET",
        "Host": "127.0.0.1:8080",
        "Proto": "HTTP/1.1",
        "Url": "/",
        "ReqURI": "/",
        "Trailer": null,
        "Body": "",
        "Headers": {
          "Accept": [
            "*/*"
          ],
          "User-Agent": [
            "curl/7.79.1"
          ]
        }
      },
      "Response": {
        "Status": 200,
        "Headers": {
          "Content-Type": [
            "application/json"
          ]
        }
      }
    }
    $ curl http://127.0.0.1:8080/status/512 -I
    HTTP/1.1 512 status code 512
    Content-Type: application/json
    Date: Tue, 02 Aug 2022 09:48:02 GMT
    Content-Length: 490

