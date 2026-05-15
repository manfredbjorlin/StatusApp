#!/bin/bash

cd /home/manfred/Development/StatusApp/cmd

GOOS=windows GOARCH=amd64 go build -o ../deployments/StatusApp.exe .
unset GOOS
unset GOARCH

mv ../deployments/StatusApp.exe /home/manfred/Sync/
