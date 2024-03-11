#!/bin/bash
lsof -t -i :50051 | xargs kill -9

cd exe
rm ./zigurat
rm -rf ./tabuas
./sacerdote &
./pilar 127.0.0.1:50051 8081 &
./pilar 127.0.0.1:50051 8082 &
./pilar 127.0.0.1:50051 8083 &
./pilar 127.0.0.1:50051 8084 &
./viga 127.0.0.1:8081 127.0.0.1:8082 127.0.0.1:8083 127.0.0.1:8084 &