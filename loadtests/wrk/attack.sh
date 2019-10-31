#!/bin/bash

for ((i=10; i<=1000; i*=10))
do 
    source $1 8 $i 30 
done