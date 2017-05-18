#!/usr/bin/env bash
set -euo pipefail
unset CDPATH; cd "$( dirname "${BASH_SOURCE[0]}" )"; cd "`pwd -P`"

# Clean
rm -rf compiled
mkdir -p compiled compiled/c

# Ensure the SDK is ready
../make.sh

# Load the go environment
eval $(../make.sh env)

# Generate the go bridge and clients
go run generate/*.go

# Ensure the go bridge is valid
echo
go install -v flywheel.io/sdk/bridge

# Generate the C bridge
echo
go build -v -buildmode=c-shared -o compiled/c/flywheel.so flywheel.io/sdk/bridge

# Group artifacts together
cp -r python/ compiled/python
