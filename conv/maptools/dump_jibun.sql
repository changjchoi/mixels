.echo on
BEGIN TRANSACTION;
CREATE TABLE jibun (
  one text, two text, three text, 
  four text, five text, six text, 
  seven text, eight text, nine text, 
  ten text, eleven text, twelve text, 
  thirteen text, fourteen text
);
.separator "|"
.import jibun_Jeonbuk_utf8.txt jibun
.import jibun_busan_utf8.txt jibun
.import jibun_chungbuk_utf8.txt jibun
.import jibun_chungnam_utf8.txt jibun
.import jibun_daegu_utf8.txt jibun
.import jibun_daejeon_utf8.txt jibun
.import jibun_gangwon_utf8.txt jibun
.import jibun_gwangju_utf8.txt jibun
.import jibun_gyeongbuk_utf8.txt jibun
.import jibun_gyeongnam_utf8.txt jibun
.import jibun_gyunggi_utf8.txt jibun
.import jibun_incheon_utf8.txt jibun
.import jibun_jeju_utf8.txt jibun
.import jibun_jeonnam_utf8.txt jibun
.import jibun_sejong_utf8.txt jibun
.import jibun_seoul_utf8.txt jibun
.import jibun_ulsan_utf8.txt jibun
COMMIT;
