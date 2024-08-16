#!/bin/bash

export PERSISTENT=$(blkid -L COS_PERSISTENT)
sudo mount $PERSISTENT /tmp/persistent
sudo mkdir /tmp/persistent/cloud-config
sudo cp kairos-config.yaml /tmp/persistent/cloud-config/
sudo umount /tmp/persistent
