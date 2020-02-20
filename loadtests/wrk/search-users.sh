#!/bin/bash
    
export R=$4
cat << EOF > ./output/users/search-users/repl_rslave/t$1-d$3-c$2-R$4
`wrk2 -t $1 -c $2 -d $3s -R$4 http://127.0.0.1/search-user`
EOF
