APPVersion?=0.1
PORT?=11004
ImageName?=sinicaapp/excel
ContainerName?=doreexcel
MKFILE := $(abspath $(lastword $(MAKEFILE_LIST)))
CURDIR := $(dir $(MKFILE))

run: 
	mkdir -p ${CURDIR}tmp
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
