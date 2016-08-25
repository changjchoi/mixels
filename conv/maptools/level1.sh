#!/bin/bash

sqlite3 < level1.sql
sort -u level1.csv > tmp
mv tmp level1.csv
