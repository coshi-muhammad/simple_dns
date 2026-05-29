all: 
	go run main.go
test: 
	dig @127.0.0.1 -p 8080 hello.com
