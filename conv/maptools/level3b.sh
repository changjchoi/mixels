#!/bin/bash

sqlite3 < level3b.sql
sort -u level3b.csv > tmp
mv tmp level3b.csv
