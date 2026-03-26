# code-interpreter

This is E2B code-interpreter template for Agent-Sandbox. It is based on the E2B official image `e2bdev/code-interpreter:latest` and adds the E2B official `agent-runtime envd`  to provide compatibility with E2B SDKs.

## Some Changes
1. main.py improve startup speed
`main.py`, was automatically create python and typescript kernels for Jupyter,
Change `main.py`, which only creates a python kernel **when first execute code**, and the typescript kernel is removed. This is because the typescript kernel is not commonly used in the code interpreter scenario, and removing it can reduce the image size and improve startup speed.

2. warmup
Change [start-up.sh](start-up.sh) to [warmup.sh](warmup.sh) and [startup.sh](startup.sh), the `warmup.sh` is used for pre-warming the container, which will be executed during the pool sandbox created, and the `startup.sh` is used for Acquire sandbox from pool, which will be executed when the user creating sandbox. This change can further improve the startup speed and reduce the resource consumption of the pool sandbox.

## build
docker build -t ghcr.io/agent-sandbox/code-interpreter:0.4.0 .

docker push ghcr.io/agent-sandbox/code-interpreter:0.4.0

