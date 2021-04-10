FROM alpine:3.13

ADD requestbin /
EXPOSE 8080
ENTRYPOINT ["./requestbin"]
