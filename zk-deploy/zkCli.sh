#!/bin/bash

docker run -it --rm --net=zookeeper_default --link zookeeper_zoo2_1:zookeeper zookeeper zkCli.sh -server zookeeper


