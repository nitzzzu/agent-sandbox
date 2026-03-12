---
name: e2b-sandbox
description: >
  Create, manage, and use E2B sandboxes — run commands, manage files, use git,
  persist state, and configure networking. Use when building agent workflows
  that need isolated execution environments.
---

# E2B Sandbox — Lifecycle & Operations

## Quick Start

Install the SDK and set your API key:

```bash
# JavaScript/TypeScript
npm install @e2b/code-interpreter

# Python
pip install e2b-code-interpreter
```

Set `E2B_API_KEY` in your environment (get one at https://e2b.dev/dashboard).

### Create a sandbox and run a command

```typescript
import { Sandbox } from '@e2b/code-interpreter'

const sandbox = await Sandbox.create()
const result = await sandbox.commands.run('echo "Hello from E2B"')
console.log(result.stdout) // Hello from E2B
await sandbox.kill()
```

```python
from e2b_code_interpreter import Sandbox

sandbox = Sandbox.create()
result = sandbox.commands.run('echo "Hello from E2B"')
print(result.stdout)  # Hello from E2B
sandbox.kill()
```

## Package Choice

| Need | Package | JS import | Python import |
|------|---------|-----------|---------------|
| Code execution (`runCode`) | Code Interpreter | `import { Sandbox } from '@e2b/code-interpreter'` | `from e2b_code_interpreter import Sandbox` |
| MCP gateway | Base SDK | `import Sandbox from 'e2b'` | `from e2b import Sandbox` |
| Templates only | Base SDK | `import { Template } from 'e2b'` | `from e2b import Template` |

**Recommendation:** Use Code Interpreter by default — it includes all base features plus `runCode()`/`run_code()`.

## Sandbox Creation Options

```typescript
const sandbox = await Sandbox.create({
  template: 'my-template',       // Custom template name or ID
  timeoutMs: 5 * 60_000,         // Sandbox lifetime (default: 5 min)
  metadata: { user: 'alice' },   // Custom metadata for filtering
  envs: { API_KEY: 'secret' },   // Runtime environment variables
  secure: true,                   // Enable HTTPS-only mode
  allowInternetAccess: false,     // Disable outbound network (default: true)
  network: {                      // Fine-grained network control
    denyOut: [ALL_TRAFFIC],
    allowOut: ['api.example.com'],
  },
})
```

```python
sandbox = Sandbox.create(
    template='my-template',
    timeout=300,                    # Sandbox lifetime in seconds
    metadata={'user': 'alice'},
    envs={'API_KEY': 'secret'},
    secure=True,
    allow_internet_access=False,
    network={
        'deny_out': [ALL_TRAFFIC],
        'allow_out': ['api.example.com'],
    },
)
```

## Running Commands

```typescript
const repoPath = '/home/user/repo'

// Simple command
const result = await sandbox.commands.run('ls -la')
console.log(result.stdout)

// With options
const result = await sandbox.commands.run('npm start', {
  cwd: repoPath,
  envs: { NODE_ENV: 'production' },
  user: 'root',
  timeoutMs: 30_000,
  onStdout: (data) => console.log(data),
  onStderr: (data) => console.error(data),
})

// Background process
const proc = await sandbox.commands.run('npm start', { background: true })
// Later...
await proc.kill()

// List running processes
const processes = await sandbox.commands.list()

// Send input to a process
await proc.sendStdin('yes\n')
```

```python
repo_path = '/home/user/repo'

# Simple command
result = sandbox.commands.run('ls -la')
print(result.stdout)

# With options
result = sandbox.commands.run(
    'npm start',
    cwd=repo_path,
    envs={'NODE_ENV': 'production'},
    user='root',
    timeout=30,
    on_stdout=lambda data: print(data),
    on_stderr=lambda data: print(data, file=sys.stderr),
)

# Background process
proc = sandbox.commands.run('npm start', background=True)
# Later...
proc.kill()

# List running processes
processes = sandbox.commands.list()

# Send input to a process
proc.send_stdin('yes\n')
```

## File Operations

```typescript
// Read a file
const content = await sandbox.files.read('/home/user/config.json')

// Write a single file
await sandbox.files.write('/home/user/hello.txt', 'Hello, World!')

// Write multiple files at once
await sandbox.files.write([
  { path: '/home/user/file1.txt', data: 'Content 1' },
  { path: '/home/user/file2.txt', data: 'Content 2' },
])

// List directory contents
const entries = await sandbox.files.list('/home/user')

// Create a directory
await sandbox.files.makeDir('/home/user/new-dir')

// Remove a file or directory
await sandbox.files.remove('/home/user/old-file.txt')

// Rename / move
await sandbox.files.rename('/home/user/old.txt', '/home/user/new.txt')

// Watch for changes
const watcher = await sandbox.files.watch('/home/user', (event) => {
  console.log(event)
})
```

```python
# Read a file
content = sandbox.files.read('/home/user/config.json')

# Write a single file
sandbox.files.write('/home/user/hello.txt', 'Hello, World!')

# Write multiple files at once
sandbox.files.write_files([
    {"path": "/home/user/file1.txt", "data": "Content 1"},
    {"path": "/home/user/file2.txt", "data": "Content 2"},
])

# List directory contents
entries = sandbox.files.list('/home/user')

# Create a directory
sandbox.files.make_dir('/home/user/new-dir')

# Remove a file or directory
sandbox.files.remove('/home/user/old-file.txt')

# Rename / move
sandbox.files.rename('/home/user/old.txt', '/home/user/new.txt')

# Watch for changes
watcher = sandbox.files.watch('/home/user', lambda event: print(event))
```

## Git Operations

### Authentication

Three options for private repositories:

**1. Inline credentials** — pass with each command:
```typescript
await sandbox.git.push(repoPath, {
  username: process.env.GIT_USERNAME,
  password: process.env.GIT_TOKEN,
})
```

```python
sandbox.git.push(
    repo_path,
    username=os.environ.get("GIT_USERNAME"),
    password=os.environ.get("GIT_TOKEN"),
)
```

**2. Credential helper** — authenticate once, reuse across commands:
```typescript
await sandbox.git.dangerouslyAuthenticate({
  username: process.env.GIT_USERNAME,
  password: process.env.GIT_TOKEN,
})
// Now HTTPS git operations use stored credentials
```

```python
sandbox.git.dangerously_authenticate(
    username=os.environ.get("GIT_USERNAME"),
    password=os.environ.get("GIT_TOKEN"),
)
```

**3. Keep credentials in remote URL** — pass `dangerouslyStoreCredentials: true` / `dangerously_store_credentials=True` to `clone()`.

### Common git workflows

```typescript
const repoPath = '/home/user/repo'

// Clone
await sandbox.git.clone('https://github.com/org/repo.git', {
  path: repoPath, branch: 'main', depth: 1,
})

// Configure identity
await sandbox.git.configureUser('Bot', 'bot@example.com')

// Branch workflow
await sandbox.git.createBranch(repoPath, 'feature/my-change')
// ... make changes ...
await sandbox.git.add(repoPath, { files: ['README.md', 'src/index.ts'] })
await sandbox.git.commit(repoPath, 'Add feature')
await sandbox.git.push(repoPath, { remote: 'origin', branch: 'feature/my-change', setUpstream: true })

// Check status
const status = await sandbox.git.status(repoPath)
console.log(status.currentBranch, status.ahead, status.behind, status.fileStatus)

// List branches
const branches = await sandbox.git.branches(repoPath)
```

```python
repo_path = '/home/user/repo'

# Clone
sandbox.git.clone(
    'https://github.com/org/repo.git',
    path=repo_path, branch='main', depth=1,
)

# Configure identity
sandbox.git.configure_user('Bot', 'bot@example.com')

# Branch workflow
sandbox.git.create_branch(repo_path, 'feature/my-change')
# ... make changes ...
sandbox.git.add(repo_path, files=['README.md', 'src/index.ts'])
sandbox.git.commit(repo_path, 'Add feature')
sandbox.git.push(repo_path, remote='origin', branch='feature/my-change', set_upstream=True)

# Check status
status = sandbox.git.status(repo_path)
print(status.current_branch, status.ahead, status.behind, status.file_status)

# List branches
branches = sandbox.git.branches(repo_path)
```

## Persistence (Pause & Resume)

Sandboxes can be paused and resumed later, preserving full state (RAM, disk, running processes).

```typescript
import { Sandbox } from '@e2b/code-interpreter'

// Create and do work
const sandbox = await Sandbox.create()
await sandbox.commands.run('echo "state to preserve" > /tmp/data.txt')
const sandboxId = sandbox.sandboxId

// Pause
await sandbox.betaPause()

// Later: resume by connecting (auto-resumes paused sandboxes)
const resumed = await Sandbox.connect(sandboxId)
const result = await resumed.commands.run('cat /tmp/data.txt')
console.log(result.stdout) // state to preserve
```

```python
from e2b_code_interpreter import Sandbox

# Create and do work
sandbox = Sandbox.create()
sandbox.commands.run('echo "state to preserve" > /tmp/data.txt')
sandbox_id = sandbox.sandbox_id

# Pause
sandbox.beta_pause()

# Later: resume by connecting (auto-resumes paused sandboxes)
resumed = Sandbox.connect(sandbox_id)
result = resumed.commands.run('cat /tmp/data.txt')
print(result.stdout)  # state to preserve
```

### Auto-pause (pause when idle instead of killing)

```typescript
const sandbox = await Sandbox.betaCreate({
  autoPause: true,
  timeoutMs: 5 * 60_000,
})
// Sandbox will pause (not kill) after 5 minutes of inactivity
```

```python
sandbox = Sandbox.beta_create(
    auto_pause=True,
    timeout=300,
)
```

### List paused sandboxes

```typescript
const paginator = await Sandbox.list({
  query: { state: ['paused'] },
})
const sandboxes = await paginator.nextItems()
```

```python
from e2b_code_interpreter import Sandbox, SandboxQuery, SandboxState

paginator = Sandbox.list(SandboxQuery(state=[SandboxState.PAUSED]))
sandboxes = paginator.next_items()
```

## Reconnecting to a Running Sandbox

```typescript
// Save the sandbox ID
const sandboxId = sandbox.sandboxId

// Later: reconnect
const sandbox = await Sandbox.connect(sandboxId)

// List all running sandboxes
const paginator = await Sandbox.list({
  query: { state: ['running'] },
})
const sandboxes = await paginator.nextItems()
```

```python
sandbox_id = sandbox.sandbox_id

# Reconnect
sandbox = Sandbox.connect(sandbox_id)

# List all running sandboxes
paginator = Sandbox.list()
sandboxes = paginator.next_items()
```

## Networking

See the [networking reference](references/networking.md) for port exposure, public URLs, traffic access tokens, domain-based filtering, and host masking.

### Quick: Get a public URL for a port

```typescript
const host = sandbox.getHost(3000)
const url = `https://${host}` // e.g., https://3000-i62mff4aht.e2b.app
```

```python
host = sandbox.get_host(3000)
url = f'https://{host}'
```
