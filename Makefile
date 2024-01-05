clean:
	@rm -rf $(BIN)

zip: build
	@zip -j $(FUNCTION).zip bin/bootstrap

.PHONY: test
test:
	@go test -v -cover ./...