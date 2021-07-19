#!/bin/bash
attempt=0
sleep 1
while [ $attempt -le 10 ]; do
    echo "Waiting for PostgreSQL to be up (attempt: $attempt)..."
    docker exec -i test-server-db psql -U test_user -d test </dev/null && sleep 1 && exit 0
    sleep 1
    attempt=$(( $attempt + 1 ))
done
echo "PostgreSQL still not ready, giving up"
exit 1
