run:
	@mkdir -p tmp
	@go build -o tmp/main main.go
	@./tmp/main
