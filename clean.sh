#!/bin/bash

set -e

rm -rf ./sign-*/example/example
rm -rf ./sign-*/example/data
rm -rf ./sign-*/example/udb
rm -rf ./reset-pwd/example/example
rm -rf ./reset-pwd/example/data
rm -rf ./reset-pwd/example/udb
rm -rf ./udb/data
rm -rf ./data

# delete all binary files
find . -type f -executable -exec sh -c "file -i '{}' | grep -q 'x-executable; charset=binary'" \; -print | xargs rm -f
for f in $(find ./ -name '*.log' -or -name '*.doc'); do rm $f; done