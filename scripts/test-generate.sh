#!/usr/bin/env bash
set -euo pipefail

ROOT="$(pwd)"
TMPDIR=$(mktemp -d /tmp/spring-cli-test-XXXX)
echo "Using tmp: $TMPDIR"
cd "$TMPDIR"

# build binary
GO_BINARY="$ROOT/spring-cli"
if [ ! -x "$GO_BINARY" ]; then
  (cd "$ROOT" && go build -o spring-cli .)
fi

# create project
"$ROOT/spring-cli" install:project maven -n testproj -p com.example.test
cd testproj

# run generators
"$ROOT/spring-cli" make entity User --fields "name:String,age:int"
"$ROOT/spring-cli" make repository User
"$ROOT/spring-cli" make service User --entity User
"$ROOT/spring-cli" make controller User --crud --entity User

# set version
"$ROOT/spring-cli" version --set 0.1.0 --commit || true

# checks
if [ ! -f src/main/java/com/example/test/entity/User.java ]; then
  echo "Missing entity"; exit 2
fi
if [ ! -f src/main/java/com/example/test/repository/UserRepository.java ]; then
  echo "Missing repository"; exit 2
fi
if [ ! -f src/main/java/com/example/test/service/UserService.java ]; then
  echo "Missing service"; exit 2
fi
if [ ! -f src/main/java/com/example/test/controller/UserController.java ]; then
  echo "Missing controller"; exit 2
fi

# check version in pom.xml
if ! grep -q "<version>0.1.0</version>" pom.xml; then
  echo "Version not updated in pom.xml"; exit 2
fi

echo "All checks passed in $TMPDIR"
