#!/bin/sh
set -ex
LOG_DIR="/service/logs/${ENV}/nonick/nonick-notifier-service"
while [ ! -d $LOG_DIR ];
do
    echo "sleep wait ${LOG_DIR}"
    sleep 1
done

./nonick-notifier-service
