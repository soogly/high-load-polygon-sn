#!/bin/bash
    
export R=$4
cat << EOF > ./output/users/search-users/repl_rslave/t$1-d$3-c$2-R$4
`wrk2 -t $1 -c $2 -d $3s -R$4 -s scripts/get-dividends.lua https://data.stage.conomy.ru/api/stock/dividend/?limit=300`


--wrk -t8 -c 1000 -d 30 -s scripts/get-dividends.lua https://10.129.0.29
