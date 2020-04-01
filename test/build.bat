set prjPath=%cd%
echo %prjPath%
cd ../../../
set GOPATH=%cd%
set GOARCH=amd64
set GOOS=windows
REM cd src/tuyue/tuyue_query/cmd/test
cd %prjPath%
go build