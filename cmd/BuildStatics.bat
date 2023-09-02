@echo off

cd ../

go test -timeout 3600s -v -run TestStatic ./static