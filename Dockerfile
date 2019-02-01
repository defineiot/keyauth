
FROM alpine:latest
ADD keyauthd /
ENTRYPOINT ./keyauthd service start -t env
EXPOSE 8080