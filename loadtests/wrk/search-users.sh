#!/bin/bash
    
# `wrk -t $1 -c $2 -d $3s --latency -s ./search-users.lua http://localhost/search-user\?search\=mar\%20kar`
export R=$4
cat << EOF > ./output/with_index/t$1-d$3-c$2-R$4
`wrk2 -t $1 -c $2 -d $3s -R$4 -s ./scripts/search-users.lua http://95.213.237.186/search-user\?search\?mar\%20kar`
EOF
