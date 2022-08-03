FROM alpine:3.13
COPY requestbin /requestbin
EXPOSE 8080
ENTRYPOINT ["/requestbin"]
