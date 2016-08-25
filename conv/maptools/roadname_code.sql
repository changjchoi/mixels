.echo on
attach database "roadname_code.db" as db1;
.output roadname_code.csv
.echo off
select distinct two, three from db1.road;
