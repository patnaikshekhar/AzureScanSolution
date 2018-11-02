$ErrorActionPreference = "Stop"

echo "Building App"
go build

echo "Building containers"
docker build -t patnaikshekhar/windows-defender .

echo "Removing existing container"
docker rm -vf scan_service

echo "Running container"
docker run -p 8000:80 --name scan_service -e AZ_ACC_NAME=$env:AZ_ACC_NAME -e AZ_ACC_KEY=$env:AZ_ACC_KEY patnaikshekhar/windows-defender