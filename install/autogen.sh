#!/bin/sh

set -e

aclocal
automake --add-missing -c
autoreconf -fi

./configure "$@"

set +e
