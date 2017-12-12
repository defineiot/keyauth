FROM alpine:latest

MAINTAINER yumaojun "18108053819@163.com"

WORKDIR /root/openauth/
ADD openauth /root/openauth/openauth
ADD conf /root/openauth/conf
VOLUME ["/root/openauth/log"]

ENV OA_MYSQL_HOST="127.0.0.1" \
    OA_MYSQL_PORT="3306" \
    OA_MYSQL_USER="openauth" \
    OA_MYSQL_PASS="openauth" \
    OA_MYSQL_DB="openauth" \
    OA_MYSQL_MAX_OPEN_CONN=1000 \
    OA_MYSQL_MAX_IDEL_CONN=200 \
    OA_MYSQL_MAX_LIFE_TIME=60 \

    OA_APP_HOST="0.0.0.0" \
    OA_APP_PORT="8080" \
    OA_APP_KEY="your app secret key" \
    OA_APP_NAME="openauth"

    OA_LOG_FILE_PATH="log/debug.log"
    OA_LOG_LEVEL="debug"


EXPOSE 8080

CMD ["./openauth", "service", "-t", "env", "start"]