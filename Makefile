
.PHONY: test
test:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

.PHONY: coverage-html
coverage-html:
	go tool cover -html=coverage.out
