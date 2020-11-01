TARGET := ./bin/kali

BASEDIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY: check
check:
	shellcheck ${TARGET}

.PHONY: install
install:
	ln -i -s -r ${BASEDIR}/${TARGET} ${HOME}/bin/

.PHONY: clean
clean:
	./bin/kali clean
