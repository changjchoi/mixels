.echo on
BEGIN TRANSACTION;
CREATE TABLE build (
  one text, two text, three text, 
  four text, five text, six text, 
  seven text, eight text, nine text, 
  ten text, eleven text, twelve text, 
  thirteen text, fourteen text, fifteen text, 
  sixteen text, seventeen text, eighteen text, 
  nineteen text, twenty text, twentyone text, 
  twentytwo text, twentythree text, twentyfour text, 
  twentyfive text, twentysix text, twentyseven text, 
  twentyeight text, twentynine text, thirty text, 
  thirtyone text
);
.separator "|"
.import build_busan_utf8.txt build
.import build_chungbuk_utf8.txt build
.import build_chungnam_utf8.txt build
.import build_daegu_utf8.txt build
.import build_daejeon_utf8.txt build
.import build_gangwon_utf8.txt build
.import build_gwangju_utf8.txt build
.import build_gyeongbuk_utf8.txt build
.import build_gyeongnam_utf8.txt build
.import build_gyunggi_utf8.txt build
.import build_incheon_utf8.txt build
.import build_jeju_utf8.txt build
.import build_jeonbuk_utf8.txt build
.import build_jeonnam_utf8.txt build
.import build_sejong_utf8.txt build
.import build_seoul_utf8.txt build
.import build_ulsan_utf8.txt build
COMMIT;
