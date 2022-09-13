#!/bin/bash

check_tcp() {
  echo 2>/dev/null > "/dev/tcp/$1/$2" && echo "passing"
}

echo -n "[$(date +"%T")] Waiting MySQL ..."
until [[ $(check_tcp db 3306) == "passing" ]]; do
  echo -n "."
  sleep 1
done
echo -ne " done\n"

sleep 1