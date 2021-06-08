set -x 
set -e
rm -rf dbserver
go build
./dbserver
