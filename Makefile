BINARY_NAME=openauthd

ifndef	_info
define _info
	echo -e "\033[44;37m [INFO] $(shell date "+%Y-%m-%d %H:%M:%S"): $1 \033[0m"
endef
endif

default:run

run:build_in_docker
	$(call _info, "put your run command to there.")

build_in_docker:
	if [ -x ${BINARY_NAME} ];then rm -rf ${BINARY_NAME}; fi
	bash ./hack/build.sh
	$(call _info, "build openauth binary file success.")

build:build_in_docker
	$(call _info, "build binary file with docker")
