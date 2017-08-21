fmt:
	find . -name '*.go' | xargs gofmt -s -w
	find . -name '*.go' | xargs goimports -w

connect:
	mysql -ugolang_int_user -ptestpass golang_intensive -h 195.211.23.152 -P 5435
