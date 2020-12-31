REM You Need MiniGW-64 to run these builds, this is usefull in case you need to build a package that requires CGO, like sqlite
@REM set CGO_ENABLED=0
@REM set CC=gcc
@REM set CXX=g++
REM set GOOS=windows
REM set GOARCH=amd64
REM go build -o output/windows/x64/OfflineSyncExporter_x64.exe main.go
set GOOS=linux
set GOARCH=amd64
go build -o output/linux/x64/WgManager .
@REM set CGO_ENABLED=0
@REM set GOOS=windows
@REM set GOARCH=amd64