.PHONY: test

test:
	go test -bench . -benchmem ./...