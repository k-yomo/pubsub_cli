
.PHONY: test
test:
	go test -coverprofile=cover.out ./...
	go tool cover -func=cover.out

.PHONY: cover-html
cover-html:
	go tool cover -html=cover.out
