# maxx
A simple grpc client server using bi-directional stream 

- client server communicate using encrypted messages (SSL/TLS) over grp bidirectional stream.
- client sends number over stream.
- server returns the current maximum number received over the stream.

Example:
Client sends over stream [21, 81, 73, 94, 54]
Server sends back over stream [21, 81, 94] i.e. current maximum

## Prerequiste:
- Install Docker Compose [doc](https://docs.docker.com/compose/install/)
- Install grpc [doc](https://grpc.io/docs/quickstart/go/)


## Run

##### Build
```sh
# build docker images
docker-compose build
```

##### Run - 1 server and 4 client container
```sh
# run the docker containers
docker-compose up --scale client=4
```

##### Run test (without creating tcp connection)
```sh
go test ./...
```

##### Note: 
- Change environment variable COUNT in docker-compose.yml to increase the random numbers being sent from client (default:5)
