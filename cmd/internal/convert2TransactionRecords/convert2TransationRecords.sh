#!/bin/bash

cat - | iconv -f BIG5 | awk -F',' '{
	# date and time
	split($1, datetime, " ")
	split(datetime[1], dateArray, "/")
	date = dateArray[1] "-" dateArray[2] "-" dateArray[3]
	time = datetime[2]

	# stockNo and stockName
	split($2, stockDetails, "(")
	split(stockDetails[2], stockNumberArray, ")")
	stockNo = stockNumberArray[1]
	stockName = stockDetails[1]
	stockName = (stockName == "") ? "N/A" : stockName


	# tranType
	if ($3 ~ "現股買進") {
		tranType = 1
	} else if ($3 ~ "現股賣出") {
		tranType = -1
	} else {
		tranType = 0
	}

	# unitPrice, quantity and status
	unitPrice = $6
	quantity = $8
	if ($12 ~ "完全成交") {
		status = 1
	} else {
		status = $12
	}


	# print out result
	printf "%s,%s,%s,%s,%s,%s,%s,%s\n",
		date, time, stockNo, stockName, tranType, unitPrice, quantity, status
}
' | grep -v "委託成功" | sort -k1,1 -k2,2
