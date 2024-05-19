#!/bin/sh

cd ..
make
cd install
zip install.zip Makefile.am bazos admin.txt.in base_pfp.png configure.ac NEWS README AUTHORS ChangeLog autogen.sh
