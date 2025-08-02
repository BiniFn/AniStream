FROM redis:8-alpine as redis

# Copy custom entrypoint
COPY docker/redis-entrypoint.sh /usr/local/bin/redis-entrypoint.sh
RUN chmod +x /usr/local/bin/redis-entrypoint.sh

ENTRYPOINT ["redis-entrypoint.sh"]
