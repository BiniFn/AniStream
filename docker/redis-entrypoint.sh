#!/bin/sh
set -e

if [ -f /run/secrets/redis_password ]; then
  PASSWORD=$(cat /run/secrets/redis_password)
  echo "requirepass $PASSWORD" > /tmp/redis.conf
  exec redis-server /tmp/redis.conf
else
  echo "No redis_password secret found! Starting without password."
  exec redis-server
fi
