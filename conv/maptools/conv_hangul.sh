#!/bin/bash
for DIR in `find . -name "*.txt"`
do
  echo "$DIR done......."
  iconv -f CP949 -t utf-8 "$DIR" > "tmp"
  rm -f "$DIR"; mv "tmp" "${DIR%.txt}_utf8.txt"
done
