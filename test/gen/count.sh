#!/bin/sh

FILE=$1

exec grep -v "^[[:space:]]*$" $1 \
    | sort \
    | uniq \
    | wc -l
