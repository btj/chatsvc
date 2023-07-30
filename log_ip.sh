#!/bin/sh
IP="$(./get_ip.sh)"
echo "$(date -Iseconds)\t$IP" >> ../ip.log
gcloud logging write --severity=INFO --payload-type=json ip-log '{ "external-ip": "'"$IP"'"}'
