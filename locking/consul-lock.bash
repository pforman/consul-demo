#!/bin/bash

echo "here we go!"
while true; do
  sid=$(consul-cli kv lock service/demo/worklock --lock-delay 1s)
  echo =====================================
  echo =====================================
  echo "look at that, I got the lock!"
  echo =====================================
  echo =====================================
  for i in 1 2 3 4 5 6 ; do
    echo sleeping....
    sleep 1
  done
  
  consul-cli kv unlock service/demo/worklock --session $sid
  s=$RANDOM
  let "s %= 5"
  let "s += 5"
  echo "released the lock, snoozing $s"
  sleep $s
done
