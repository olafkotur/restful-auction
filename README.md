# SCC.311 Distributed Systems Coursework

## Vagrant: Setup
* `vagrant up`
* `vagrant ssh`

## Docker: Build & Run
* Running without optional arguments will build / launclsh all services
* `docker-compose build` optional `server` | `load-balancer`
* `docker-compose up` optional `server` | `load-balancer`

## Docker: Deployment
* `docker login -u "$USERNAME" -p "$PASSWORD" harbor.scc.lancs.ac.uk`
* `docker tag my-server harbor.scc.lancs.ac.uk/$USERNAME/server:latest`
* `docker push harbor.scc.lancs.ac.uk/$USERNAME/server:latest`

## Testing HTTP responses
* Simple request to uri
* `curl -Lvk localhost:8080/api/auctions`
* Passing additonal headers | specifying HTTP request methods
* `curl -Lvk -X POST -H "Header: Value" localhost:8080/api/auctions`
* Passing data via the -d flag
* `curl -Lvk -X POST -d "name=Motorbike&startBid=0.99" localhost:8080/api/auction`

## Running
* Compile: Make sure to be in the `/out/` dir and run: `make`
* Execute: `java -cp ".:../libs/*" Class` replacing Class

