ROOT_PATH=$(shell pwd)
RELEASE=/bin/
APP_NAME=onbio_web

all : clean proc

proc :
	cd $(ROOT_PATH) && \
        go build -o $(ROOT_PATH)$(RELEASE)$(APP_NAME) -mod vendor

clean:
	@rm -rf $(ROOT_PATH)$(RELEASE)*
