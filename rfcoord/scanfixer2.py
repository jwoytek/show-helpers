#!/usr/bin/env python3

import os
import sys
import csv

with open(sys.argv[1], 'r') as in_f, open(sys.argv[2], 'w') as out_f:
    csv_read = csv.reader(in_f, delimiter=',')
    csv_write = csv.writer(out_f, delimiter=',')

    for row in csv_read:
        # start capturing after this value
        if row[0] == "Trace Data":
            break

    for row in csv_read:
        if len(row) != 2:
            # we're probably done
            break
        try:
            f = float(row[0])
            f = f / 1000000
            s = float(row[1])
            csv_write.writerow([f, s])
        except ValueError:
            break

