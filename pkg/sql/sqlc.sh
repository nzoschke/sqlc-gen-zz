#!/bin/bash

rm -rf c
sqlc generate
go fmt ./c
