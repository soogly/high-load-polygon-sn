#!/bin/bash
    
export R=$4
cat << EOF > ./output/users/search-users/with_index_UNION/t$1-d$3-c$2-R$4
`wrk2 -t $1 -c $2 -d $3s -R$4 -s ./scripts/search-users-one-pref.lua http://95.213.237.186/search-user`
EOF
