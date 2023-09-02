@echo off

cd ../

go build -ldflags "-s -X web-server-ui/commons.version=0.0.1" .