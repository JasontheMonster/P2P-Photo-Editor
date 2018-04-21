all:
	go build *.go
run:
	go run *.go -id=1 -addr=127.0.0.1:8888
clear:
	rm const
