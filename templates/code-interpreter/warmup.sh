#!/bin/bash

echo "Starting Code Interpreter server..."
MATPLOTLIBRC=/root/.config/matplotlib/.matplotlibrc jupyter server --ip=0.0.0.0 --no-browser --IdentityProvider.token=""