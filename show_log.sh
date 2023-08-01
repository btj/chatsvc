#!/bin/sh
gcloud logging read 'jsonPayload."cos.googleapis.com/container_name":chatsvc' --format=json | jq -j '.[].jsonPayload.message'
