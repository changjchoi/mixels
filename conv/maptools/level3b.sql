.echo on
attach database "build.db" as db1;
attach database "jibun.db" as db2;
.output level3b.csv
.echo off
select distinct one, five from db1.build where five !="";
select distinct one, five from db2.jibun where five !="";
