test:
	go test -v ./...

generate:
	go generate ./...

deploy:
	cd lambda && $(MAKE) deploy