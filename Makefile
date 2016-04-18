default: test

test:
	@go test -race $(shell go list ./... | grep -v vendor)
fmt:
	@go fmt $(shell go list ./... | grep -v vendor)
license:
	@go run gen_license.go ./

check_fmt:
ifneq ($(shell gofmt -l ./ | grep -v vendor | grep -v testdata),)
	$(error code not fmted, run make fmt. $(shell gofmt -l ./ | grep -v vendor | grep -v testdata))
endif

check_license:
ifneq ($(shell go run gen_license.go ./),)
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
	brew install oniguruma
endif

travis_test: export PKG_CONFIG_PATH += $(PWD)/vendor/github.com/limetext/rubex:$(GOPATH)/src/github.com/limetext/rubex
travis_test: test cover report_cover

cover:
	@echo "mode: count" > coverage.cov; \
	for pkg in $(shell go list "./..." | grep -v /vendor/); do \
		go test -covermode=count -coverprofile=tmp.cov $$pkg; \
		sed 1d tmp.cov >> coverage.cov; \
		rm tmp.cov; \
	done

report_cover:
ifeq ($(REPORT_COVERAGE),true)
	$(shell go env GOPATH | awk 'BEGIN{FS=":"} {print $1}')/bin/goveralls -coverprofile=coverage.cov -service=travis-ci
endif
	rm coverage.cov
