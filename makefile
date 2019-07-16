BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
APPVersion?=0.1
ContainerTagName?=its/excel
APP?=app
PORT?=11004
ImageName?=devdockersrv.test5.sinica.edu.tw/its/excel
ContainerName?=doreexcel
MKFILE := $(abspath $(lastword $(MAKEFILE_LIST)))
CURDIR := $(dir $(MKFILE))

build:
	GOOS=linux GOARCH=amd64 go build -tags netgo \
	-ldflags "-s -w -X version.BuildTime=${BUILD_TIME}" \
	-o ${APP}

docker: build
	docker build -t ${ImageName}:${APPVersion} .
	rm -f app
	docker images

dockerpush:
	docker tag ${ImageName}:${APPVersion} ${ImageName}:${APPVersion}
	docker push ${ImageName}:${APPVersion}

run: docker
	docker run --restart=always -d --name ${ContainerName} -v /etc/localtime:/etc/localtime:ro \
        -v ${CURDIR}tmp:/tmp/excel \
	--env-file ./envfile \
	-p ${PORT}:80 ${ImageName}:${APPVersion}
	
stop:
	docker stop ${ContainerName}
	
test:
	curl -X POST -F "xlsx=@${CURDIR}2018.xlsx" https://devwteamapi.test5.sinica.edu.tw/excel/excel2json

log:
	docker logs -f -t --tail 20 ${ContainerName}

rm:
	docker rm ${ContainerName}
