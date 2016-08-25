.echo on
attach database "build.db" as db1;
attach database "jibun.db" as db2;
.output level2.csv
.echo off
select distinct substr(one, 1, 5), three from db1.build;
select distinct substr(one, 1, 5), three from db2.jibun;
