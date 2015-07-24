build:
	go build
	go install
test:
	go test -v
	rm -rf ~/test/
	rm -rf ~/.rdb/test-prog/*
	rm -rf ~/test2/
