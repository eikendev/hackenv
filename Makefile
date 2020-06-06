SHARED_DIR := shared

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

.PHONY: permissions
permissions:
	chmod 777 "./${shared}"
	sudo semanage fcontext -a -t svirt_image_t "${PWD}/${shared}(/.*)?"
	sudo restorecon -vrF "${PWD}/${shared}"
