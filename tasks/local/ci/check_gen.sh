#!/usr/bin/env bash

source "$(dirname -- "$0")/../../general/ci/warn.sh"
source "$(dirname -- "$0")/../../general/ci/setup.sh"

ret=0

fold_start "gen.python" "check python api"
diff_test "go run $(dirname -- "$0")/../build/gen_python_api.go"
let ret=$ret+$test_result
fold_end "gen.python"

fold_start "gen.loaders" "check loaders"
diff_test "go run $(dirname -- "$0")/../build/gen_loaders.go"
let ret=$ret+$test_result
fold_end "gen.loaders"

exit $ret
