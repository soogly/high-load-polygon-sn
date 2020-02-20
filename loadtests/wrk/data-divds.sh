#!/bin/bash
    
export R=$4
cat << EOF > ./output/conomy/data/get-dividends-SBER/t$1-d$3-c$2-R$4
`wrk2 -t $1 -c $2 -d $3s -R$4 https://data.stage.conomy.ru/api/stock/dividends-by-ticker/SBER`
EOF
