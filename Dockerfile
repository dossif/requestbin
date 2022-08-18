FROM alpine:3.13
COPY requestbin /requestbin
RUN chmod +x /requestbin
EXPOSE 8080
ENTRYPOINT ["/requestbin"]
