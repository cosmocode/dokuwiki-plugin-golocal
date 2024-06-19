
NAME := golocal
OUT := ../dl

clean:
	rm -f $(OUT)/*

windows64:
	CGO_ENABLED=1 \
	GOOS=windows \
	GOARCH=amd64 \
	CC=x86_64-w64-mingw32-gcc \
	go build -o $(OUT)/$(NAME)-windows64.exe -ldflags -H=windowsgui .

windows32:
	CGO_ENABLED=1 \
	GOOS=windows \
	GOARCH=386 \
	CC=i686-w64-mingw32-gcc \
	go build -o $(OUT)/$(NAME)-windows32.exe -ldflags -H=windowsgui .

windows: windows64 windows32

linux64:
	GOOS=linux
	GOARCH=amd64
	go build -o $(OUT)/$(NAME)-linux64 .

linux32:
	GOOS=linux
	GOARCH=i386
	go build -o $(OUT)/$(NAME)-linux32 .


linux: linux64 linux32

all: windows linux
