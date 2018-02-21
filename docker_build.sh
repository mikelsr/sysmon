#!/bin/bash

docker build . | tee /tmp/sysmon_build_out
did=$(cat /tmp/sysmon_build_out | grep "Successfully built " | cut -d" " -f3)

docker tag $did sysmon:latest

echo "Built with ID:${did} NAME:sysmon"
