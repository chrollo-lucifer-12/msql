#!/usr/bin/env bash

APP_NAME="msql"
MAIN="."

run () {
    go run "$MAIN"
}

build () {
    go build -o "bin/$APP_NAME" "$MAIN"
}

test () {
    read -rp "enter package name: " pkg
    cd "$pkg"
    go test -v
}


case "$1" in
    run)
        run
        ;;
    build)
        build
        ;;
    test)
        test
        ;;
esac
