#!/bin/sh

set -e
cd ..
make
cd install
aclocal
automake --add-missing -c
autoreconf -fi

./configure "$@"

set +e
