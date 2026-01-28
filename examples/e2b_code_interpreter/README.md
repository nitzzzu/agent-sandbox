# E2B Code Interpreter Usage Examples

This project demonstrates how to use the `e2b-code-interpreter` SDK to create and manage code execution sandbox environments.

## Features

- ✅ Create new code execution sandboxes
- ✅ Connect to existing sandboxes
- ✅ Execute Python code in sandboxes
- ✅ File upload and management
- ✅ Port forwarding and server access
- ✅ Real-time logging and error handling

## Environment Configuration

Before use, you need to set the following environment variables:

```python
import os

# Set E2B service domain
os.environ['E2B_DOMAIN'] = 'your.domain.com'

# Set API URL
os.environ['E2B_API_URL'] = 'http://localhost:10000/e2b/v1'

# Set API key, e.g. testuser-aef134ef-7aa1-945e-9399-7df9a4ad0c3f
os.environ['E2B_API_KEY'] = 'your-api-key-here'

# Optional: Enable debug mode
os.environ['E2B_DEBUG'] = "true"
```

## Installation

### config agent-sandbox server ingress endpoint
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: agent-sandbox
  namespace: agent-sandbox
spec:
  rules:
    - host: "*.your.domain.com"
      http:
        paths:
          - backend:
              service:
                name: agent-sandbox
                port:
                  number: 80
            path: /
            pathType: ImplementationSpecific
```

## Usage Examples

```bash
pip install -r requirements.txt
```

### 1. Create New Sandbox

Create a new code execution sandbox:

```python
from e2b_code_interpreter import Sandbox

# Create sandbox with default configuration
# Default template: "code-interpreter-v1"
# Default timeout: 300 seconds
sbx = Sandbox.create()

# Or specify template and timeout
sbx = Sandbox.create(template="code-interpreter", timeout=300)

print(f"Sandbox ID: {sbx.sandbox_id}")
```

### 2. Connect to Existing Sandbox

Connect to an existing sandbox:

```python
from e2b_code_interpreter import Sandbox

# Connect by sandbox ID
sbx = Sandbox.connect("sandbox-id-here")
print(f"Sandbox ID: {sbx.sandbox_id}")
```

### 3. Get Sandbox Information

Retrieve detailed information about a sandbox:

```python
info = sbx.get_info()
print(info)

# Example output:
# SandboxInfo(
#   sandbox_id='ig6f1yt6idvxkxl562scj-419ff533',
#   template_id='u7nqkmpn3jjf1tvftlsu',
#   name='base',
#   metadata={},
#   started_at=datetime.datetime(2025, 3, 24, 15, 42, 59, 255612, tzinfo=tzutc()),
#   end_at=datetime.datetime(2025, 3, 24, 15, 47, 59, 255612, tzinfo=tzutc())
# )
```

### 4. Execute Code

Execute Python code in the sandbox:

```python
# Simple code execution
execution = sbx.run_code("print('hello world')")
print(execution.logs)

# Execute mathematical operations
execution = sbx.run_code("print(1+1)")
print(execution.logs)
```

### 5. Code Execution with Callbacks

Use callback functions to handle output and errors during execution:

```python
code_to_run = """
import time
import sys
print("This goes first to stdout")
time.sleep(3)
print("This goes later to stderr", file=sys.stderr)
time.sleep(3)
print("This goes last")
"""

sbx.run_code(
    code_to_run,
    # Handle runtime errors
    on_error=lambda error: print('error:', error),
    # Handle standard output
    on_stdout=lambda data: print('stdout:', data),
    # Handle standard error output
    on_stderr=lambda data: print('stderr:', data),
)
```

### 6. File Operations

#### List Files

```python
# List files in root directory
files = sbx.files.list("/")
print(files)

# List files in specified directory
files = sbx.files.list("/home/user")
print(files)
```

#### Upload Files

```python
# Read from local filesystem and upload to sandbox
with open("test.md", "rb") as file:
    sbx.files.write("/home/user/test2.md", file)
```

### 7. Port Forwarding and Server Access

Start a server in the sandbox and access it through port forwarding:

```python
import requests
from e2b_code_interpreter import Sandbox

sandbox = Sandbox.connect("sandbox-id")

# Start HTTP server in sandbox (run in background)
sandbox.commands.run("python -m http.server 8080", background=True)

# Get host address for port forwarding
host = sandbox.get_host(8080)
url = f"http://{host}"

# Access the server
response = requests.get(url)
print(response.status_code)  # Note: Requests without token may return 403
```

### 8. Close Sandbox

Close the sandbox when done:

```python
sbx.kill()
```

## Development Mode Configuration

In development environments, you can customize connection configuration:

```python
def dev():
    def __connection_config_get_host(_, sandbox_id: str, sandbox_domain: str, port: int) -> str:
        print(f"host request params sandbox_id={sandbox_id}, sandbox_domain={sandbox_domain}, port={port}")
        # Return custom host address
        return f"10.100.21.221:{port}"
    
    from e2b import ConnectionConfig
    ConnectionConfig.get_host = __connection_config_get_host

if os.environ.get('E2B_DEBUG'):
    dev()
```

## Test Functions

The code includes the following test functions:

- `create_test()`: Demonstrates how to create a new sandbox
- `get_test()`: Demonstrates how to connect to an existing sandbox and get information
- `port_test()`: Demonstrates port forwarding and server access
- `connect()`: Demonstrates complete workflow of connection, code execution, and file operations

## Running Examples

Run the main program:

```bash
python main.py
```

Or call specific test functions:

```python
if __name__ == "__main__":
    # Choose which test function to run
    get_test()      # Connect and get sandbox info
    # main()        # Create new sandbox
    # port_test()   # Test port forwarding
    # connect()     # Full feature test
```

## Notes

1. **Sandbox Timeout**: Default timeout is 300 seconds, adjust as needed
1. **Port Access**: When accessing services in sandbox through port forwarding, authentication token may be required
1. **Resource Cleanup**: Remember to call `kill()` method to close the sandbox after use to avoid resource waste

## More Information

- [E2B Official Documentation](https://docs.e2b.dev/)
- [e2b-code-interpreter SDK](https://github.com/e2b-dev/e2b-code-interpreter)
