@echo off
cd /d %~dp0
set profle=%1

set GO111MODULE=on
set GOPROXY=https://mirrors.aliyun.com/goproxy

setlocal
set PROTOC_DIR=protoc
set PATH=%PATH%;%PROTOC_DIR%

set PROTO_DIR=..\\proto
set JAVASCRIPT_DIR=..\\proto\\js
set DBPROTO_DIR=..\\db

protoc -I=%PROTO_DIR% --go_out=%PROTO_DIR% --go-grpc_out=%PROTO_DIR% %PROTO_DIR%\*.proto
protoc -I=%PROTO_DIR% --js_out=import_style=commonjs,binary:%JAVASCRIPT_DIR% %PROTO_DIR%\*.proto

protoc -I=%DBPROTO_DIR% --go_out=%DBPROTO_DIR% %DBPROTO_DIR%\*.proto

@if %errorlevel%==0 (
    @echo protobuf generated!
) else (
    @echo protobuf generate fail!
    goto :fail
)

protoc.exe -I=%PROTO_DIR% --go_out=%PROTO_DIR% %PROTO_DIR%\*.proto
go generate
@if %errorlevel%==1 (
    @echo go generate fail!
    goto fail:
) else (
    @echo go code generated!
)

set GOOS=windows

set GOARCH=amd64
set CGO_ENABLED = 0
set GOOS
set GOARCH
set CGO_ENABLED
go build -o goserver.exe
@if %errorlevel%==0 (
    @echo go build success!
    goto :end
) else (
    echo go build failed!
)

:fail
rem echo failed.

:end
endlocal
rem set PATH
echo bye