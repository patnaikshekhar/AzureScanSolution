echo "Building App"
go build

echo "Removing existing container"
docker rm -vf av

echo "Building container"
docker build -t patnaikshekhar/clamavrest .

echo "Running container"
docker run -p 8080:8000 --name av patnaikshekhar/clamavrest