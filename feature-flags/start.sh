#!/bin/bash

# If consul isn't up we're going to have a bad demo
consulcheck=$(docker inspect consul | jq .[].State.Running)
if [ $consulcheck != "true" ]; then
  echo "cannot find the consul container"
  echo "are you doing this right?"
  exit 1
fi

for i in 0 1 2 3 4 5 ; do
  port="888$i"
  echo -n "started port $port - container "
  docker run -e "CONSUL_ADDR=consul:8500" --link consul -d -p $port:80 --name $port feature-flags
done

