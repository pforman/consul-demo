#!/bin/bash

set -e
host=$(hostname)

unlock ()
{
  consul-cli --consul consul:8500 kv unlock service/demo/worklock --session $1
  cp /tmp/unlocked /var/www/html/index.nginx-debian.html
}

cat >/tmp/locked<<EOF
<html>
  <body bgcolor="#aa0000">
    <h1>Working!</h1>
    <h3>(in bash...)</h3>
    <p>Hostname: $host</p>
  </body>
</html>
EOF

cat >/tmp/unlocked<<EOF
<html>
  <body bgcolor="#ffffff">
      <h1>Not working. Sad.</h1>
    <h3>(still in bash...)</h3>
    <p>Hostname: $host</p>
  </body>
</html>
EOF

echo "here we go!"
# Start with the unlock page instead of the Nginx default
cp /tmp/unlocked /var/www/html/index.nginx-debian.html
while true; do
  sid=$(consul-cli --consul consul:8500 kv lock service/demo/worklock --lock-delay 0)
  trap "unlock $sid; echo exiting on signal; exit" SIGINT SIGTERM
  echo =====================================
  echo =====================================
  echo "look at that, I got the lock!"
  echo =====================================
  echo =====================================
  cp /tmp/locked /var/www/html/index.nginx-debian.html
  for i in 1 2 3 4 5 6 ; do
    echo sleeping....
    sleep 1
  done
  
  #consul-cli --consul consul:8500 kv unlock service/demo/worklock --session $sid
  #cp /tmp/unlocked /var/www/html/index.nginx-debian.html
  unlock $sid
  s=$RANDOM
  let "s %= 5"
  let "s += 5"
  echo "released the lock, snoozing $s"
  sleep $s
done


  
