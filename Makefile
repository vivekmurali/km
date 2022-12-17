serve:
	go build -o server/httpserver ./server/ 
	./server/httpserver

cli:
	go build -o km.out .
	./km.out

build-server:
	go build -o server/httpserver ./server/
	service km restart
