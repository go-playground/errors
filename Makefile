GOCMD=go

linters-install:
	@gometalinter --version >/dev/null 2>&1 || { \
		echo "installing linting tools..."; \
		$(GOCMD) get github.com/alecthomas/gometalinter; \
		gometalinter --install; \
	}

lint: linters-install
	gometalinter --vendor --disable-all --enable=vet --enable=vetshadow --enable=golint --enable=megacheck --enable=ineffassign --enable=misspell --enable=errcheck --enable=goconst ./...

test:
	$(GOCMD) test -cover -race ./...

bench:
	$(GOCMD) test -bench=. -benchmem ./...

.PHONY: test lint linters-install