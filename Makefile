all: oar-build

oar-build:
	export GOPROXY=https://goproxy.cn,direct
	mkdir -p bin
	go build -ldflags "-X main.GitCommitId=$(shell if which git > /dev/null; then git describe --always --dirty --long --tags --abbrev=40; else echo $(CI_COMMIT_SHA); fi)" -x -o bin/oar github.com/forestsay/oar-dns

clean:
	rm -r bin
