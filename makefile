BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
APPVersion?=0.0.1
ContainerTagName?=justgps/excel
APP?=app
PORT?=11004
ImageName?=justgps/excel
ContainerName?=doreexcel
MKFILE := $(abspath $(lastword $(MAKEFILE_LIST)))
CURDIR := $(dir $(MKFILE))

build:
	go env -w GOPRIVATE=github.com/asccclass
	GOOS=linux GOARCH=amd64 GO111MODULE=off go build -tags netgo \
	-ldflags "-s -w -X version.BuildTime=${BUILD_TIME}" \
	-o ${APP}

docker: build
	docker build -t ${ImageName}:${APPVersion} .
	rm -f app
	docker images

push: docker
	docker tag ${ImageName}:${APPVersion} ${ImageName}:${APPVersion}
	docker push ${ImageName}:${APPVersion}

run: docker
	docker run --restart=always -d --name ${ContainerName} \
	-v /etc/localtime:/etc/localtime:ro \
        -v ${CURDIR}tmp:/tmp \
	-v ${CURDIR}envfile:/app/envfile \
	--env-file ./envfile \
	-p ${PORT}:80 ${ImageName}:${APPVersion}
	sh clean.sh
	
stop:
	docker stop ${ContainerName}
	
test:
	curl -X POST -F "xlsx=@${CURDIR}2018.xlsx" https://devwteamapi.test5.sinica.edu.tw/excel/excel2json

log:
	docker logs -f -t --tail 20 ${ContainerName}

rm:
	docker rm ${ContainerName}

re:stop rm run

login:
	docker exec -it ${ContainerName} /bin/sh
