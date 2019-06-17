BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
APP?=app
PORT?=11004
ImageName?=its/excel
ContainerName?=doreexcel
MKFILE := $(abspath $(lastword $(MAKEFILE_LIST)))
CURDIR := $(dir $(MKFILE))

build:
	GOOS=linux GOARCH=amd64 go build -tags netgo \
	-ldflags "-s -w -X version.BuildTime=${BUILD_TIME}" \
	-o ${APP}

docker: build
	docker build -t ${ImageName} .
	rm -f app
	docker images


run: docker
	docker run --rm -d --name ${ContainerName} -v /etc/localtime:/etc/localtime:ro \
	--env-file ./envfile -p ${PORT}:80 ${ImageName}
	
stop:
	docker stop ${ContainerName}
	
test:
	curl -X POST -F "xlsx=@${CURDIR}2018.xlsx" https://devwteamapi.test5.sinica.edu.tw/excel/excel2json

log:
	docker logs -f -t --tail 20 ${ContainerName}

rm:
	docker rm ${ContainerName}
