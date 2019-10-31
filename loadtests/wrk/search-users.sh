#!/bin/bash
    
# `wrk -t $1 -c $2 -d $3s --latency -s ./search-users.lua http://localhost/search-user\?search\=mar\%20kar`

cat << EOF > ./accountings/c$2
`wrk2 -t $1 -c $2 -d $3s -R2000 --latency -s ./search-users.lua http://localhost/search-user\?search\?mar\%20kar`
EOF
