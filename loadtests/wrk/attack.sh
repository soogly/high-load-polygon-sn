#!/bin/bash

for ((i=1; i<=1000; i*=10))
do 
    if [[ $i == 1 ]]
    then
        p=1
    else
        p=8
    fi
    source $1 $p $i 60 1
    sleep 10
done