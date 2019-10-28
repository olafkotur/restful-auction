#!/bin/bash

echo Compiling server...
cd ./server
go build -o ../out/
cd ../

echo Compiling client...
cd ./client
go build -o ../out/
cd ../

echo Compiling load-balancer...
cd ./load-balancer
go build -o ../out/
cd ../