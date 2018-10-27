$env:GOOS = "linux"
$ErrorActionPreference = "Stop"

echo "Building App"
go build

echo "Building containers"
docker-compose build

echo "Removing existing stack"
docker-compose down
docker-compose rm -vf

echo "Running stack"
docker-compose up -d