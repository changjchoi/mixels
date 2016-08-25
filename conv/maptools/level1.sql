.echo on
attach database "build.db" as db1;
.output level1.csv
.echo off
select distinct substr(one, 1, 2), two from db1.build;
