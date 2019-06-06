#!/bin/bash
docker build -f ci/Dockerfile.build -t flod-build:latest .
id=$(docker create flod-build:latest)
for BIN in addblock findcheckpoint floctl flod gencerts; do docker cp $id:/go/$BIN ./$BIN; done