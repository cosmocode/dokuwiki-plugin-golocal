
NAME := golocal
OUT := ../dl

clean:
	rm -f $(OUT)/*

windows:
	CGO_ENABLED=1 \
	GOOS=windows \
	GOARCH=amd64 \
	CC=x86_64-w64-mingw32-gcc \
	go build -o $(OUT)/$(NAME)-windows.exe -ldflags -H=windowsgui .


linux:
	GOOS=linux \
	GOARCH=amd64 \
	go build -o $(OUT)/$(NAME)-linux .

macos:
	CGO_ENABLED=1 \
	GOOS=darwin \
	GOARCH=amd64 \
	CC=o64-clang \
	go build -o $(OUT)/$(NAME)-macos .

all: clean linux windows macos
