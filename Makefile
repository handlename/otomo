test:
	go tool gotestsum --format testdox

ci:
	go tool gotestsum --format github-actions --junitfile test-report.xml

generate:
	go generate ./...

deploy:
	cd lambda && $(MAKE) deploy
