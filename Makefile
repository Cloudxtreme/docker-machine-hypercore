default: build

build: 
	GOGC=off go build -i -o docker-machine-hypercore ./bin

clean:
	$(RM) docker-machine-hypercore

install: build
	cp ./docker-machine-hypercore /usr/local/bin/docker-machine-hypercore

.PHONY: install
