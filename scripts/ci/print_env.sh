#!/usr/bin/env bash

echo "working directory"
pwd

echo -e "\ngo env"
go version
go env

echo -e "\ngit config"
git config -l
