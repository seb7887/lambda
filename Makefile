clean:
	@rm -rf $(BIN)

zip: build
	@zip -j $(FUNCTION).zip bin/bootstrap

.PHONY: test
test:
	@go test -v -short -coverpkg=./... -coverprofile=cov.out ./...

cov:test
	@go tool cover -html=./cov.out -o ./cov.html