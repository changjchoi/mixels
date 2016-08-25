#!/bin/bash

sqlite3 < level3a.sql
sort -u level3a.csv > tmp
mv tmp level3a.csv
