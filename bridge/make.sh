#!/usr/bin/env bash
set -euo pipefail
unset CDPATH; cd "$( dirname "${BASH_SOURCE[0]}" )"; cd "`pwd -P`"

# Get system info
localOs=$( uname -s | tr '[:upper:]' '[:lower:]' )

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
echo "Building the C bridge..."
if [[ "$localOs" == "darwin" ]]; then
    ext="dylib"
else
    ext="so"
fi
go build -buildmode=c-shared -o dist/c/flywheelBridge.$ext flywheel.io/sdk/bridge/dist

# Matlab wants a simpler copy of the header file
cp dist/c/flywheelBridge.* dist/matlab/
# Remove typedef and line precompiler directive, as they confuse matlab
sed -i '/^typedef /d; /^\#line /d;' dist/matlab/flywheelBridge.h
# Rename file
mv dist/matlab/flywheelBridge.h dist/matlab/flywheelBridgeSimple.h
