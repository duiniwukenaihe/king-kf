#!/bin/sh
[ "$MYSQL_USER" ] || MYSQL_USER='root'
[ "$MYSQL_PASSWORD" ] || MYSQL_PASSWORD='1232456'
[ "$MYSQL_HOST" ] || MYSQL_HOST='10.10.20.13'
[ "$MYSQL_PORT" ] || PORT="3306"
[ "$MYSQL_DB" ] || DB="kingfisher"
[ "$DB_URL" ] || DB_URL="${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/${MYSQL_DB}"

[ "$MQ_USER" ] || MQ_USER=admin
[ "$MQ_PASSWORD" ] || MQ_PASSWORD=admin	
[ "$MQ_HOST" ] || MQ_HOST=10.10.20.13
[ "$MQ_PORT" ] || MQ_PORT=5672
[ "$RABBITMQ_URL" ] || RABBITMQ_URL="amqp://${MQ_USER}:${MQ_PASSWORD}@${MQ_HOST}:${MQ_PORT}/"

[ "$LISTEN" ] || LISTEN=0.0.0.0
[ "$PORT" ] || PORT=8080
[ "$TIME_ZONE" ] || TIME_ZONE="Asia/Shanghai"
[ "$ALPINE_REPO" ] || ALPINE_REPO="mirrors.aliyun.com"

sed -i "s/dl-cdn.alpinelinux.org/${ALPINE_REPO}/g" /etc/apk/repositories     
apk --no-cache add tzdata 
echo "${TIME_ZONE}" > /etc/timezone 
ln -sf /usr/share/zoneinfo/${TIME_ZONE} /etc/localtime 
mkdir /lib64 
ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

/usr/local/bin/king-kf -dbURL=$DB_URL -listen=$LISTEN:$PORT -rabbitMQURL=$RABBITMQ_URL
