precommit: fmt license test

test:
	@go test -race ./...
fmt:
	@go fmt ./...
license:
	@go run $(GOPATH)/src/github.com/limetext/tasks/gen_license.go
fast_test:
	@go test ./...

check_fmt:
ifneq ($(shell gofmt -l ./ | grep -v testdata),)
	$(error code not fmted, run make fmt. $(shell gofmt -l ./ | grep -v testdata))
endif

check_license:
	@go run $(GOPATH)/src/github.com/limetext/tasks/gen_license.go -check

tasks:
	go get -d -u github.com/limetext/tasks

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
	for pkg in $$(go list "./..."); do \
		go test -covermode=count -coverprofile=tmp.cov $$pkg; \
		sed 1d tmp.cov >> coverage.cov; \
		rm tmp.cov; \
	done

report_cover:
ifeq ($(REPORT_COVERAGE),true)
	$$(go env GOPATH | awk 'BEGIN{FS=":"} {print $1}')/bin/goveralls -coverprofile=coverage.cov -service=travis-ci
endif
	rm coverage.cov
