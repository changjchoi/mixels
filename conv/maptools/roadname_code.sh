#!/bin/bash
sqlite3 < roadname_code.sql 
sort -u roadname_code.csv > tmp
mv tmp roadname_code.csv
