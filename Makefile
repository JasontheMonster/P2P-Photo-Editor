all:
	go build *.go 
run1:
	python app/client.py 127.0.0.1:5050 127.0.0.1:5051&
	go run *.go -id=1 -addr=127.0.0.1:8080 -listenAddr=127.0.0.1:5051 -sendAddr=127.0.0.1:5050

run2:
	python app/client.py 127.0.0.1:5052 127.0.0.1:5053&
	go run *.go -id=2 -addr=127.0.0.1:8081 -listenAddr=127.0.0.1:5053 -sendAddr=127.0.0.1:5052

run3:
	python app/client.py 127.0.0.1:5054 127.0.0.1:5055&
	go run *.go -id=3 -addr=136.167.207.7:8082 -listenAddr=127.0.0.1:5055 -sendAddr=127.0.0.1:5054
clear:
	rm const
