all: 
	go run main.go

debug:
	go build -gcflags="all=-N -l" main.go
	gdlv debug


test: 
	./test.sh
