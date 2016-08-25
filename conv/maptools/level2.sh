#!/bin/bash

sqlite3 < level2.sql
sort -u level2.csv > tmp
mv tmp level2.csv
