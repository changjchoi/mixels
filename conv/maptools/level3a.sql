.echo on
attach database "build.db" as db1;
attach database "jibun.db" as db2;
.output level3a.csv
.echo off
select distinct substr(one, 1, 8), four from db1.build;
select distinct substr(one, 1, 8), four from db2.jibun;
