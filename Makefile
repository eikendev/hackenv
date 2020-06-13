TARGET := kali

.PHONY: check
check:
	shellcheck ${TARGET}

.PHONY: start
start:
	./kali $@

.PHONY: stop
stop:
	./kali $@

.PHONY: install
install:
	./kali $@

.PHONY: download
download:
	./kali $@

.PHONY: ssh
ssh:
	./kali $@

.PHONY: gui
gui:
	./kali $@

.PHONY: clean
clean:
	./kali $@

.PHONY: list
list:
	./kali $@

.PHONY: share
share:
	./kali $@

.PHONY: permissions
permissions:
	./kali $@
