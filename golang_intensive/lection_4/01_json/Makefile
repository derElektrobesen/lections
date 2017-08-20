fmt:
	find . -name '*.go' | xargs gofmt -s -w
	find . -name '*.go' | xargs goimports -w

easy:
	easyjson -all -snake_case test/05_easyjson.go
