TARGET := kali

.PHONY: check
check:
	shellcheck ${TARGET}

.PHONY: start
start:
	./bin/kali $@

.PHONY: stop
stop:
	./bin/kali $@

.PHONY: install
install:
	./bin/kali $@

.PHONY: download
download:
	./bin/kali $@

.PHONY: ssh
ssh:
	./bin/kali $@

.PHONY: gui
gui:
	./bin/kali $@

.PHONY: clean
clean:
	./bin/kali $@

.PHONY: list
list:
	./bin/kali $@

.PHONY: share
share:
	./bin/kali $@

.PHONY: permissions
permissions:
	./bin/kali $@
