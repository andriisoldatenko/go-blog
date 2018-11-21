run:
	go run *.go

build:
	go build -o blog *.go

debug:
	dlv debug github.com/andriisoldatenko/go-blog
