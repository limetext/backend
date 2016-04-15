default: test

test:
	@go test -race ./lib/...
fmt:
	@go fmt ./lib/...
license:
	@go run license.go lib

checkfmt:
ifneq ($(shell gofmt -l lib),)
	$(error code not fmted, run make fmt. $(shell gofmt -l lib))
endif

checklisence:
ifneq ($(shell go run license.go lib),)
	$(error license is not added to all files, run make license)
endif

glide:
	go get -v -u github.com/Masterminds/glide
	glide install
cover_dep:
	go get -v -u github.com/mattn/goveralls
	go get -v -u github.com/axw/gocov/gocov

travis:
ifeq ($(TRAVIS_OS_NAME),osx)
	brew update
	brew install oniguruma python3
endif

travis_test: test cover report_cover

cover:
	@echo "mode: count" > coverage.cov; \
	for pkg in $(shell go list "./lib/..." | grep -v /vendor/); do \
		go test -covermode=count -coverprofile=tmp.cov $$pkg; \
		sed 1d tmp.cov >> coverage.cov; \
		rm tmp.cov; \
	done

report_cover:
ifeq ($(REPORT_COVERAGE),true)
	$(shell go env GOPATH | awk 'BEGIN{FS=":"} {print $1}')/bin/goveralls -coverprofile=coverage.cov -service=travis-ci
endif
	rm coverage.cov
