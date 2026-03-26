#!/bin/bash
echo "Starting envd..."
/workspace/envd/envd > /proc/1/fd/1 2>&1 &

echo "Start Code Service at $(date '+%H:%M:%S')..."
cd /root/.server/
.venv/bin/uvicorn main:app --host 0.0.0.0 --port 49999 --workers 1 --no-access-log --no-use-colors --timeout-keep-alive 640 > /proc/1/fd/1 2>&1 &


PORT=${1:-49999}
TIMEOUT=${2:-10}

START_TIME=$(date +%s)
echo "Services started at $(date '+%H:%M:%S'), waiting for port $PORT..."

success=0
MAX_RETRIES=$((TIMEOUT * 5))
for i in $(seq 1 $MAX_RETRIES); do
    if (echo > /dev/tcp/localhost/$PORT) >/dev/null 2>&1; then
        END_TIME=$(date +%s)
        ELAPSED=$((END_TIME - START_TIME))
        echo "------------------------------------------"
        echo "Check Success: Port $PORT is reachable!"
        echo "Total Time Elapsed: ${ELAPSED}s"
        echo "------------------------------------------"
        success=1
        break
    fi
    echo "Waiting for port $PORT... (${i}/${TIMEOUT}s)"
    sleep 0.2
done

if [ $success -eq 0 ]; then
    END_TIME=$(date +%s)
    ELAPSED=$((END_TIME - START_TIME))
    echo "Error: Timeout waiting for port $PORT after ${ELAPSED}s."
    exit 1
fi

exit 0
