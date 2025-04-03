#!/usr/bin/env bash

# Wait for all dockers to be healthy for 1 minute if not error
TIMEOUT=60
RETRY_SLEEP=1
PENDING_DOCKERS=1
START_TIME=$(date +%s)
while [ $PENDING_DOCKERS -ne 0 ]; do
    PENDING_DOCKERS=$(docker compose ps -q \
    | xargs -I {} docker inspect \
        --format='{{if .Config.Healthcheck}}{{.Name}} => {{.State.Health.Status}}{{end}}' {} \
    | sed  '/^$/d' | grep -v healthy | wc -l)
    if [ $PENDING_DOCKERS -gt 0 ]; then
        CURRENT_TIME=$(date +%s)
        ELAPSED_TIME=$(($CURRENT_TIME - $START_TIME))
        if [ $ELAPSED_TIME -gt $TIMEOUT ]; then
            echo "❌ Timeout waiting for dockers to be healthy! error"
            exit 1
        fi

        echo "⏳ Waiting for $PENDING_DOCKERS dockers to be healthy..."
        sleep $RETRY_SLEEP
    fi
done
