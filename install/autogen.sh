#!/bin/sh

set -e

automake --add-missing -c
autoreconf -i

./configure

set +e
