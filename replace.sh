#!/bin/sh

# Check if two arguments are provided
if [ $# -ne 2 ]; then
    echo "Usage: $0 <string1> <string2>"
    exit 1
fi

# Assign arguments to variables
string1="$1"
string2="$2"

# Find files containing string1 and replace with string2
find . -type f -exec sed -i "s/$string1/$string2/g" {} +
