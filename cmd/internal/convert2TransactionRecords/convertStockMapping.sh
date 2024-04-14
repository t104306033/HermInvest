#!/bin/bash

cat - | awk -F',' '!seen[$3,$4]++ {
    stockNo = $3
    stockName = $4

    # print out result
    printf "%s,%s\n", stockNo, stockName
}' | sort -k1,1
