#!/bin/bash
IFS=$'\n'

pkgs=`go list ./ && go list ./*/`
pkgs=($pkgs); unset IFS;
pkgs=`printf ',%s' "${pkgs[@]}"`
pkgs=${pkgs:1}

go test -v -race -coverprofile=coverage.txt -covermode=atomic \
    -coverpkg=$pkgs
