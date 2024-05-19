#!/bin/sh

cd ..
make
cd install
zip install.zip Makefile.am bazos admin.txt.in ../images/base_pfp.png configure.ac NEWS README AUTHORS ChangeLog autogen.sh
