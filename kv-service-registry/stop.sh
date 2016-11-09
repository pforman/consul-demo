#!/bin/bash

for i in 0 1 2 3 4 5 ; do
  port="888$i"
  id=$(docker inspect $port | jq -r .[].Id)
  echo -n "removing port $port - container "
  docker stop $id
  docker rm   $id
done

