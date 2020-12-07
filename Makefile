.PHONY: run test
run:
	go run cmd/server/main.go -p 8080

test:
	bash ./test.sh
