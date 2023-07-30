#!/bin/sh
gcloud compute instances describe chatsvc \
  --format='get(networkInterfaces[0].accessConfigs[0].natIP)'
