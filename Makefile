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
	chmod 777 "./${SHARED_DIR}"
	sudo semanage fcontext "${PWD}/${SHARED_DIR}(/.*)?" --deleteall
	sudo semanage fcontext -a -t svirt_image_t "${PWD}/${SHARED_DIR}(/.*)?"
	sudo restorecon -vrF "${PWD}/${SHARED_DIR}"
