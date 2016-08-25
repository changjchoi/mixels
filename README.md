Mixels - geocoding
==================

Google Like Geocoding Server and Converter

### Quick Start
* [지도](https://dl.dropboxusercontent.com/u/64402891/addr-map.tar.gz)를 다운 받습니다.
 * /home1/dragon/data 에 압축을 풉니다.
* config.json 파일을 수정합니다.
 * OutputPath를 /home1/dragon/data/addr-map으로 설정합니다.
* mixels 디렉토리에서 아래 명령을 수행합니다.
```bash
$ make 
$ make push-sphinx
$ make start-sphinx
$ make run
```

### Installation
#### goquadtree
```sh
$ go get github.com/volkerp/goquadtree/quadtree
```
* QuadTree Algorithm 
  * Shape의 Boundary를 이용하여 4분할 Tree를 생성합니다.

#### go_hangul
```sh
$ go get github.com/suapapa/go_hangul
```
* cp949를 UTF-8로 변경하기 위해 사용합니다.

#### go-proj-4
##### proj.4
* go-proj-4가 wrapper한 원본 c/c++ 입니다.
* [proj.4][PROJ] 소스를 다운 받아서 /usr/local 에 설치 하세요 
* 다른 곳에 설치 했다면 CGO_CFLAGS와 CGO_LDFLAGS를 적절히 변경하세요.
```sh
$ export CGO_LDFLAGS=-L/opt/local/lib
$ export CGO_CFLAGS=-I/opt/local/include
```
##### go-proj-4
```sh
$ go get github.com/pebbe/go-proj-4/proj
```
* shape 좌표를 경위도 좌표로 변경하기 위해 사용합니다.

#### go-shp
```sh
$ go get github.com/jonas-p/go-shp
```
* shape 파일을 읽고 처리하기 위해서 사용합니다.

#### go-sql-driver
```sh
$ go get github.com/go-sql-driver/mysql
```
* sphinx에서 mysql CLI를 사용하므로 go에서 사용가능한 mysql driver

#### sphinx install
* 메뉴얼 보시고 알아서 설치해 주세요. 

### How to Convert a Map
#### Tools - sqlite3 or mysql
도로명과 주소를 컨버전하기 위해 사용했습니다. 단순히 코드와 명칭을 뽑기 때문에 다른 DB를 사용해도 무방합니다.

#### 지도 데이터 디렉토리 구성
```sh
choeui-Mac-mini:map choechangjin$ ls
addr-map-raw/ addr-map/    
choeui-Mac-mini:addr-map-raw choechangjin$ ls
201510RDNMADR/     building1/         level2.csv         roadname_code.csv
201510RDNMCODE/    building2/         level3a.csv
boundary/          level1.csv         level3b.csv
```
* 컨버전된 지도 : ~/addr-map
* 법정동 및 도로명 지도 : 압축 풀리는 이름으로
* 행정경계 지도 : ~/boundary
* 건축물 용도별 전국지도 : ~/building1
  * 단독주택만 넣었습니다.
* 건축물 용도별 전국지도 : ~/building2
  * 단독주택을 제외한 나머지 shapefile

#### 법정동 코드 만들기
지도는 [여기서][NAD] 다운 받습니다.압축을 풀면 201510RDNMADR 디렉토리가 생성됩니다.

* 전체주소를 다운 받습니다. (2015년 10월분 - 개발시점)
* maptools의 script를 이용하여 cp949 한글을 코드를 utf-8 type으로 변경합니다.
* sqlite를 이용하여 build.db와 jibun.db를 생성합니다.
* 해당 Db를 이용하여 level1.csv, level2.csv, level3a.csv, level3b.csv를 생성합니다.
* 만들어진 csv 파일을 addr-map-raw 디렉토리로 이동합니다.

```sh
$ ./conv_hangul.sh
$ sqlite3 build.db < dump_build.sql
$ sqlite3 jibun.db < dump_jibun.sql
$ ./level1.sh
$ ./level2.sh
$ ./level3a.sh
$ ./level3b.sh
$ mv *.csv ../addr-map-raw
```

#### 도로명 코드 만들기
* [위와 같은 주소][NAD]에서 도로명코드를 다운 받습니다. 압축을 풀면 201510RDNMCODE 디렉토리가 생성됩니다.
* maptools의 script를 이용하여 cp949 한글 코드를 utf-8 코드로 변경합니다.
* road_code.db를 생성하고 road_code.csv를 생성합니다.
* 만들어진 csv를 addr-map-raw 디렉토리로 이동합니다.

```sh
$ ./conv_hangul.sh
$ sqlite3 road_code.db < dump_road.sql
$ ./roadname_code.sh
$ mv *.csv ../addr-map
```

#### 시도/시군구/읍면동/리 Shapefile
지도는 [여기][SIG]에 있습니다. 4개의 지도를 다운 받고 압축을 풀어서 적당한 폴더에 넣습니다. 다만, 2015년 6월 ‘시도’ 데이터가 깨져서 2014년 7월 데이터를 사용하였습니다.

```sh
choeui-Mac-mini:boundary choechangjin$ ls
emd.dbf         li.dbf          sig.dbf         tl_scco_ctprvn.dbf
emd.prj         li.prj          sig.prj         tl_scco_ctprvn.prj
emd.shp         li.shp          sig.shp         tl_scco_ctprvn.shp
emd.shx         li.shx          sig.shx         tl_scco_ctprvn.shx
```

#### 건축물 용도별 지도 Shapefile
지도는 [여기서][BDS] 다운 받습니다. 적당한 디렉토리에 넣습니다.약 831개 파일로
구성되어 있습니다. 

```sh
choeui-Mac-mini:building choechangjin$ ls 
이(미)용원.dbf*  이(미)용원.shp*  이(미)용원.shx*  제.dbf*  제.shp*  제.shx*
교회.dbf*  차고.dbf*  교회.shp*  차고.shp*  교회.shx*  차고.shx*
…
```

### Usages
* 기본 서비스 지도는 아래와 같은 명령으로 만들 수 있습니다.
* 10분 이내로 끝나지 않으면, timeout flag를 setting 해야 합니다.
* 10분 이내로 끝내기 위해서 건축물 용도별 지도를 2개로 분리했습니다.
* 컨버전 시간도 약간 줄었습니다. 약 8분 정도 소요됩니다.
```sh
$ go test -v -run "ConvAll"
```
* sphinx에 index를 걸어줄 주소 데이터를 dump 합니다. 
* format은 아래와 같습니다. 
```xml
id[TAB]주소[TAB]{id:level7-id, i:"level7-filename"}
```
* level7 파일을 2개로 갈라놓았기 때문에 검색에 level7 파일명을 넣어 두었습니다.
* 이런 format으로 만들어진 주소데이터는 약 819M입니다. 
* 원본 데이터는 중복을 포함하여 10587430개 입니다. 
  * 중복제거하고 collected 7299532 docs, 539.1 MB 로 출력이 되는 군요
```sh
$ go test -v -run "DumpAddress"
```
* 간략한 테스트로 index를 줄여서 만드려면 index parameter를 주면 됩니다.
```sh
$ go test -v -run "DumpAddress" -index=0,1
```
* 만들어진 address.tsv 파일을 sphinx 인덱를 걸어주려면 희수형이 만들어 놓은
sphinx.sh 를 직접 찾아서 실행하거나 아래와 같이 테스트를 걸어 줘도 됩니다.
```sh
$ go test -v -run "PushSphinx" 
```
* 다 만들어진 지도 위에서 서비스를 실행하려면 아래와 같이 합니다.
```sh
$ make run
or
$ go run main.go
or
$ go build
$ ./mixels
```
* 웹 브라우저 8000 port에 아래 주소를 치면 결과를 볼 수 있습니다.
* 주소 형식은 wgs84, latitude와 longitude를 넣는 순서로 합니다.
  * 외국 사이트도 이런 순서로 넣고
  * 비행기, 배에서도 이런 식으로 좌표를 불러 주는 것으로 보입니다.
```
http://localhost:8000/geocoding?latlon=37.49884,127.04750
```
* 결과는 아래와 같이 나옵니다.
```json
{
  "Results": [
    {
      "Components": [
        {
          "Name": "서울특별시",
          "Type": "Level1"
        },
        {
          "Name": "강남구",
          "Type": "Level2"
        },
        {
          "Name": "역삼동",
          "Type": "Level3A"
        },
        {
          "Name": "역삼로 306",
          "Type": "Level4"
        },
        {
          "Name": "개나리 래미안",
          "Type": "Level5"
        },
        {
          "Name": "754번지",
          "Type": "Level6"
        },
        {
          "Name": "106동",
          "Type": "Level7"
        }
      ],
      "DetailAddress": "서울특별시 강남구 역삼로 306, 106동 (역삼동, 개나리 래미안)",
      "Geometry": {
        "Location": {
          "X": 127.0475,
          "Y": 37.49884
        },
        "ViewPort": {
          "MinX": 127.04718103867589,
          "MinY": 37.49869554466216,
          "MaxX": 127.0477057891719,
          "MaxY": 37.49898877511604
        }
      }
    }
  ],
  "Status": "OK"
}
```
* 주소는 comma로 구분하여 넣습니다. 
* 공백글자는 +로 넣습니다.
* 맨 끝에 있는 토큰으로 검색을 합니다. 
* 앞에 있는 지역은 단지 검색 결과를 필터링하기 위해서 넣습니다.
  * 시도, 시군구, 읍면동, 리, 검색어
  * 시도, 읍면동, 검색어
  * 시도, 시군구, 검색어
  * 시도, 검색어
  * 읍면동, 검색어
  * 검색어
```
http://localhost:8000/geocoding?address=경기도,용인시+수지구,서홍마을+4단지+401
```
```json
{
  "Results": [
    {
      "Components": [
        {
          "Name": "경기도",
          "Type": "Level1"
        },
        {
          "Name": "용인시 수지구",
          "Type": "Level2"
        },
        {
          "Name": "신봉동",
          "Type": "Level3A"
        },
        {
          "Name": "신봉1로 28",
          "Type": "Level4"
        },
        {
          "Name": "서홍마을 4단지",
          "Type": "Level5"
        },
        {
          "Name": "893번지",
          "Type": "Level6"
        },
        {
          "Name": "401동",
          "Type": "Level7"
        }
      ],
      "DetailAddress": "경기도 용인시 수지구 신봉1로 28, 401동 (신봉동, 서홍마을 4단지)",
      "Geometry": {
        "Location": {
          "X": 127.08398172422055,
          "Y": 37.321261052765685
        },
        "ViewPort": {
          "MinX": 127.08212470518461,
          "MinY": 37.320699105639406,
          "MaxX": 127.08270312947857,
          "MaxY": 37.3209057864556
        }
      }
    }
  ],
  "Status": "OK"
}
```

### 환경설정 - mapconfig.json
```json
{                                                                                
  "ServicePort" : ":8000",                                                       
  "InputPath" : "/Users/choechangjin/happiness/map/addr-map-raw",                
  "OutputPath" : "/Users/choechangjin/happiness/map/addr-map",                   
  "Mapping" : {                                                                  
    "boundary/tl_scco_ctprvn.shp" : {                                            
      "File":"level1.bin", "FieldType": {                                        
        "CTPRVN_CD" : "level1_code"                                              
      }                                                                          
    },                                                                           
    "boundary/sig.shp" : {                                                       
      "File":"level2.bin", "FieldType": {                                        
        "SIG_CD" : "level2_code"                                                 
      }                                                                          
    },                                                                           
    "boundary/emd.shp" : {                                                       
      "File":"level3a.bin", "FieldType": {                                       
        "EMD_CD" : "level3a_code"                                                
      }                                                                          
    },                                                                           
    "boundary/li.shp" : {                                                        
      "File":"level3b.bin", "FieldType": {                                       
        "LI_CD" : "level3a_code"                                                 
      }                                                                          
    }, 
    "building1/*" : {                                                            
      "File":"level71.bin", "FieldType": {                                       
        "BULD_MNNM" : "build_main_num",                                          
        "BULD_NM" : "build_name",                                                
        "BULD_NM_DC" : "build_name_dtl",                                         
        "BULD_SLNO" : "build_sub_num",                                           
        "EMD_CD" : "level3a_code",                                               
        "GRO_FLO_CO" : "floor",                                                  
        "LI_CD" : "level3b_code",                                                
        "LNBR_MNNM" : "jibun_main_num",                                          
        "LNBR_SLNO" : "jibun_sub_num",                                           
        "MNTN_YN" : "is_mountain",                                               
        "RN_CD" : "roadname_code",                                               
        "SIG_CD" : "level2_code"                                                 
      }                                                                          
    },                                                                           
    "building2/*" : {                                                            
      "File":"level72.bin", "FieldType": {                                       
        "BULD_MNNM" : "build_main_num",                                          
        "BULD_NM" : "build_name",                                                
        "BULD_NM_DC" : "build_name_dtl",                                         
        "BULD_SLNO" : "build_sub_num",                                           
        "EMD_CD" : "level3a_code",                                               
        "GRO_FLO_CO" : "floor",                                                  
        "LI_CD" : "level3b_code",                                                
        "LNBR_MNNM" : "jibun_main_num",                                          
        "LNBR_SLNO" : "jibun_sub_num",                                           
        "MNTN_YN" : "is_mountain",                                               
        "RN_CD" : "roadname_code",                                               
        "SIG_CD" : "level2_code"                                                 
      }                                                                          
    }                                                                            
  },                                                                             
  "AreaLevel" : {                                                                
    "level1idx.bin" : "Level1",                                                  
    "level2idx.bin" : "Level2",                                                  
    "level3aidx.bin" : "Level3A",                                                
    "level3bidx.bin" : "Level3B",                                                
    "level71idx.bin" : "Level4",                                                 
    "level72idx.bin" : "Level4"                                                  
  }, 
  "RegionCodeFile" : {                                                           
    "RoadNameCode" : "roadname_code.csv",                                        
    "Level1" : "level1.csv",                                                     
    "Level2" : "level2.csv",                                                     
    "Level3A" : "level3a.csv",                                                   
    "Level3B" : "level3b.csv"                                                    
  },                                                                             
  "LevelList" : [                                                                
    {"Level":"Level1", "Priority":1},                                            
    {"Level":"Level2", "Priority":2},                                            
    {"Level":"Level3A", "Priority":3},                                           
    {"Level":"Level3B", "Priority":4},                                           
    {"Level":"Level4", "Priority":5},                                            
    {"Level":"Level5", "Priority":6},                                            
    {"Level":"Level6", "Priority":7},                                            
    {"Level":"Level7", "Priority":8}                                             
  ],
  "TSVFile" : ["address0.tsv", "address1.tsv", "address2.tsv", "address3.tsv"]
} 
```
* ServicePort : 웹 서비스 포트
* InputPath : 원본지도가 들어있는 디렉토리
* OutputPath : 컨버전된 파일이 들어갈 디렉토리
* Mapping : 시도/시군구/읍면동 Shapefile에서 읽어야 하는 속성값
  * boundary/sig.shp : 컨버전할 파일명 
    * File : 컨버전된 결과값을 저장할 파일명
    * FieldType : Shapefile에서 읽어야 하는 필드와 그 값을 저장할 변수 
* AreaLevel : 해당 파일이 어떤 Level인지 정의
* RegionCodeFile : 법정동 및 도로명 코드를 읽기 위한 파일
  * 건축물 용도별 속성에 들어있는 코드값을 명칭으로 변경하기 위한 파일명 

### 지도공유
* [지도](https://dl.dropboxusercontent.com/u/64402891/addr-map.tar.gz)

### Bugs
#### 지도 관련
* 대학교.shp 에 서울대학교 명칭이 들어가 있지 않음
* 희수형 아파트 이름이 없음
* 원본 데이터의 형상이 googlemap이나 daum map과 정확히 일치 하는 것이 거의 없음.
  * 약간 위 아래로 shift되어 있음. 

#### Server 관련


### References
* [건축물 용도별 전국 건물 ShapeFile Field 정의][BDS]
* [도로명주소 구축 가이드][RODG]

[NAD]: <http://www.juso.go.kr/support/AddressBuild.do>
[BDS]: <http://www.gisdeveloper.co.kr/entry/건축물-용도별-전국-건물-SHP-파일-다운로드?category=24>
[SIG]: <http://www.gisdeveloper.co.kr/979>
[PROJ]: <https://github.com/OSGeo/proj.4>
[BLEV]: <http://www.blevesearch.com>
[RODG]: <http://www.juso.go.kr/dn.do?fileName=%5B가이드%5D주소구축%28전체주소%29활용방법.pdf&realFileName=address_build_total_guide.pdf&regYmd=2014>
