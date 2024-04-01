#!/bin/bash

cat - | awk -F',' '{
    date = $5
    stockNo = $2; stockName = $1
    tranType = $3; unitPrice = $6; quantity = $4

    # time
    # The transaction time cannot overlap. If it is the same day,
    # the time will be increased by 10 seconds than the previous time.
    if (prev_date == date) {
        prev_time = time
        cmd = "date +\"%H:%M:%S\" -d\""prev_time " 10 seconds\""
        cmd | getline time
        close(cmd)
    } else {
        time = "09:00:00"
    }
    prev_date = date

    # Manual input status is always 1
    status = 1

	# print out result
	printf "%s,%s,%s,%s,%s,%s,%s,%s\n",
		date, time, stockNo, stockName, tranType, unitPrice, quantity, status
}
' | sort -k1,1 -k2,2
