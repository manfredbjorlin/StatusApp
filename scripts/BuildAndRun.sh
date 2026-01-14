#!/bin/bash

cd /home/manfred/Development/StatusApp/cmd
go build -o ../deployments/StatusApp
cd ../deployments/
./StatusApp
