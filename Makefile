all: 
	go run main.go

build:
	go build main.go -o simple_dns

debug:
	go build -gcflags="all=-N -l" main.go
	gdlv debug


test: 
	./test.sh
