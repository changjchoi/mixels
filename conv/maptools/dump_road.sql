.echo on
BEGIN TRANSACTION;
CREATE TABLE road (
  one text, two text, three text, 
  four text, five text, six text, 
  seven text, eight text, nine text, 
  ten text, eleven text, twelve text, 
  thirteen text, fourteen text, fifteen text, 
  sixteen text, seventeen text, eighteen text, 
  nineteen text, twenty text
);
.separator "|"
.import road_code_total_utf8.txt road
COMMIT;
