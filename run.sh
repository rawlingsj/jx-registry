#!/bin/bash -e

if [ -f ".jx/variables.sh" ]; then
  echo "sourcing .jx/variables.sh"
  source .jx/variables.sh
else
  echo "file does not exist: .jx/variables.sh"
fi

echo "verifying the container registry is setup"
jx-registry create