#!/usr/bin/env bash
set -euo pipefail
unset CDPATH; cd "$( dirname "${BASH_SOURCE[0]}" )"; cd "`pwd -P`"

# Clean
rm -f dist/bridge.go dist/python/flywheel.py dist/c/flywheel.h

# Ensure the SDK is ready
../make.sh

# Load the go environment
eval $(../make.sh env)

# Generate the go bridge and clients
go run generator/*.go

# Ensure the go bridge is valid
# Only necessary when testing changes to the Go template.
# echo
# go install -v flywheel.io/sdk/bridge/dist

# Generate the C bridge
echo
go build -v -buildmode=c-shared -o dist/c/flywheel.so flywheel.io/sdk/bridge/dist
