set prjPath=%cd%
echo %prjPath%
cd ../../
set GOPATH=%cd%
set GOARCH=amd64
set GOOS=linux
cd %prjPath%
go build -a -v -ldflags="-s -w"