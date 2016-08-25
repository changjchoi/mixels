# Some os test 
# someday, it will be used 
ifeq ($(OS),Windows_NT)
    OS += WIN32
    ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
        OS += AMD64
    endif
    ifeq ($(PROCESSOR_ARCHITECTURE),x86)
        OS += IA32
    endif
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
	      RELEASE := $(shell lsb_release -r | cut -d: -f2)
	      DISTRI := $(shell lsb_release -i | cut -d: -f2)
        OS += LINUX
    endif
    ifeq ($(UNAME_S),Darwin)
        OS += OSX
    endif
    UNAME_P := $(shell uname -p)
    ifeq ($(UNAME_P),x86_64)
        MACHINE += AMD64
    endif
    ifneq ($(filter %86,$(UNAME_P)),)
        MACHINE += IA32
    endif
    ifneq ($(filter arm%,$(UNAME_P)),)
        MACHINE += ARM
    endif
endif

# get user name
USER = $(shell whoami)
# get user major group
GROUP = $(shell id -g -n $(USER) )
# check sphinx installed
SEARCHED = $(shell which searchd)
# ld returns waring, if the library exists
LIB_PROJ = $(shell ld -lproj 2>&1 | grep -c "warning" && rm -f a.out)
# delete a.out if the ld found some
DEL = $(shell rm -f a.out)

all: proj.4 go-packages sphinx mixels

# Build proj.4 source
proj.4:
ifeq ($(LIB_PROJ),0)
  @sudo git clone https://github.com/OSGeo/proj.4 tmp && \
		cd tmp && sudo ./autogen.sh && \
		sudo ./configure --prefix=/usr/local && \
		sudo make install && \
		cd .. && sudo rm -rf tmp 
  ifeq ($(OS),LINUX)
    @sudo ln -s /usr/local/lib/libproj.so.10 /usr/lib
  endif
else
	@echo "skip proj.4. already installed"
#	@rm -f a.out
endif

# Build go-proj-4. it tests a library and include path
go-packages:
	go get github.com/pebbe/go-proj-4/proj
	go get github.com/jonas-p/go-shp
	go get github.com/suapapa/go_hangul
	go get github.com/volkerp/goquadtree/quadtree
	go get github.com/go-sql-driver/mysql
#	@rm -f a.out

# Main mixels build
mixels: clean
	go build
#	@rm -f a.out

# First, the sphinx was installed. 
# Make a proper directory for mixels
sphinx:
ifneq ("$(wildcard /home1/dragon/sphinx/bin/searchd)","")	
	@echo "skip sphinx setting. already installed"
#	@rm -f a.out
else 
#	@rm -f a.out
	@echo "make directory /home1/dragon/sphinx ~"
	@sudo mkdir -p /home1/dragon/sphinx/bin
	@sudo mkdir -p /home1/dragon/sphinx/etc
	@sudo mkdir -p /home1/dragon/sphinx/var/data
	@sudo mkdir -p /home1/dragon/sphinx/var/log
	@sudo mkdir -p /home1/dragon/sphinx/source
	@sudo chgrp -R $(GROUP) /home1
	@sudo chown -R $(USER) /home1
#  ifeq ($(OS),LINUX)
#		@sudo add-apt-repository ppa:builds/sphinxsearch-daily
#		@sudo apt-get update
#		@sudo apt-get install sphinxsearch
#  else 
#    ifeq ($(OS),Darwin)
#      ifeq ($(shell which port),/opt/local/bin/port)
#	      @sudo port install sphinx
#      else ifeq ($(shell which brew),/usr/local/bin/brew)
#	      @sudo brew install sphinx
#      else 
#	      echo "do nothing"
#      endif
#    endif
#  endif
	@sudo cp `which searchd` /home1/dragon/sphinx/bin
	@sudo cp `which indexer` /home1/dragon/sphinx/bin
endif

map-download:
#	@rm -f a.out
	@mkdir -p /home1/dragon/Data
	@cd /home1/dragon/Data && git init && \
	  git remote add -f origin https://github.com/airof98/Data && \
		git config core.sparseCheckout true && \
		echo "/addr-map" >> .git/info/sparse-checkout && \
		git pull origin master && \
		git lfs pull

conv-map:
	@go test -v -run "ConvAll"

dump-address:
#	@rm -f a.out
	@go test -v -run "DumpAddress"

push-sphinx:
#	@rm -f a.out
	@go test -v -run "PushSphinx"

start-sphinx:
	@/home1/dragon/sphinx/bin/searchd -c /home1/dragon/sphinx/etc/sphinx.conf

stop-sphinx:
	@/home1/dragon/sphinx/bin/searchd -c /home1/dragon/sphinx/etc/sphinx.conf \
		--stop
run:
	@./mixels > /dev/null 2>&1 &

clean:
	@rm -f a.out
	@rm -f mixels

