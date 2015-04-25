#!/usr/bin/env bash

source "$(dirname -- "$0")/../../general/ci/setup.sh"

fold_start "get" "get dependencies"
go get github.com/limetext/lime-backend/lib/...
fold_end "get"
