$env:GOOS = "linux"

echo "Building App"
go build

echo "Removing existing container"
docker rm -vf av

echo "Building container"
docker build -t patnaikshekhar/clamavrest .

echo "Running container"
docker run -p 8080:8000 -e AZ_ACC_NAME=$AZ_ACC_NAME -e AZ_ACC_KEY=$AZ_ACC_KEY --name av patnaikshekhar/clamavrest