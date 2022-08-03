FROM alpine:3.13
COPY requestbin /requestbin123
EXPOSE 8080
ENTRYPOINT ["/requestbin123"]
