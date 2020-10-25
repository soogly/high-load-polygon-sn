#!/bin/bash


echo "starting boil app"
./boil &
echo "ok"
echo 

echo "starting nginx"
nginx -g "daemon off;"
echo "ok"
echo