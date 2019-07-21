# maxx
A simple grpc client server using bi-directional stream 

- client server communicate using encrypted messages (SSL/TLS) over grp bidirectional stream.
- client sends number over stream.
- server returns the current maximum number received over the stream.

**Example:**
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
Output:
```sh
Starting maxx_server_1 ... done
Starting maxx_client_1               ... done
Starting maxx_client_2               ... done
Starting maxx_client_3               ... done
Starting maxx_client_4               ... done
Attaching to maxx_server_1, maxx_client_1, maxx_client_2, maxx_client_3, maxx_client_4
server_1  | 2019/07/21 12:52:53 listen on port 8443
server_1  | 2019/07/21 12:52:53 client disconnected.
server_1  | 2019/07/21 12:52:53 client disconnected.
server_1  | 2019/07/21 12:52:53 client disconnected.
server_1  | 2019/07/21 12:52:53 client disconnected.
client_1  | 2019/07/21 12:52:53 sent:[65 11 1 26 13] recv:[65]
client_3  | 2019/07/21 12:52:53 sent:[22 84 65 89 3] recv:[22 84 89]
client_4  | 2019/07/21 12:52:53 sent:[20 31 22 68 21] recv:[20 31 68]
client_2  | 2019/07/21 12:52:53 sent:[4 83 66 47 90] recv:[4 83 90]
maxx_client_2 exited with code 0
maxx_client_1 exited with code 0
maxx_client_4 exited with code 0
maxx_client_3 exited with code 0
```

##### Run test (without creating tcp connection)
```sh
go test ./...
```

##### Note: 
- Change environment variable COUNT in docker-compose.yml to increase the random numbers being sent from client (default:5)
