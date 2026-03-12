# E2B Common Mistakes

Two-column reference: wrong patterns on the left, correct patterns on the right.

## Construction & Imports

| ❌ Wrong | ✅ Correct | Why |
|----------|-----------|-----|
| `new Sandbox()` / `Sandbox()` | `await Sandbox.create()` / `Sandbox.create()` | `create()` is a static factory method, not a constructor |
| `import Sandbox from '@e2b/code-interpreter'` | `import { Sandbox } from '@e2b/code-interpreter'` | Named export, not default (JS) |
| `from e2b import Sandbox` (for code exec) | `from e2b_code_interpreter import Sandbox` | Code interpreter is a separate package |
| `import { Sandbox } from 'e2b'` (for code exec) | `import { Sandbox } from '@e2b/code-interpreter'` | Code interpreter is a separate package |

## Namespacing

| ❌ Wrong | ✅ Correct | Why |
|----------|-----------|-----|
| `sandbox.run('ls')` | `sandbox.commands.run('ls')` | Commands are under `sandbox.commands` namespace |
| `sandbox.exec('ls')` | `sandbox.commands.run('ls')` | Method is `run`, not `exec` |
| `sandbox.read('/file')` | `sandbox.files.read('/file')` | Files are under `sandbox.files` namespace |
| `sandbox.write('/file', data)` | `sandbox.files.write('/file', data)` | Files are under `sandbox.files` namespace |
| `sandbox.clone(url)` | `sandbox.git.clone(url, { path })` | Git is under `sandbox.git` namespace |
| `sandbox.execute(code)` | `sandbox.runCode(code)` | Code execution is `runCode`/`run_code` |

## Timeout Units

| ❌ Wrong | ✅ Correct | Why |
|----------|-----------|-----|
| `Sandbox.create({ timeoutMs: 60 })` (JS) | `Sandbox.create({ timeoutMs: 60_000 })` | JS uses **milliseconds** — 60 = 60ms, not 60s |
| `Sandbox.create(timeout=60_000)` (Python) | `Sandbox.create(timeout=60)` | Python uses **seconds** — 60_000 = ~16 hours |
| `sandbox.setTimeout(300)` (JS) | `sandbox.setTimeout(300_000)` | 300ms vs 5 minutes |
| `sandbox.set_timeout(300_000)` (Python) | `sandbox.set_timeout(300)` | 300_000 seconds = 3.5 days |

## Sandbox Lifecycle

| ❌ Wrong | ✅ Correct | Why |
|----------|-----------|-----|
| `sandbox.pause()` | `sandbox.betaPause()` / `sandbox.beta_pause()` | Pause is a beta feature, method has `beta` prefix |
| `Sandbox.resume(id)` | `Sandbox.connect(id)` | `connect()` auto-resumes paused sandboxes |
| `sandbox.id` | `sandbox.sandboxId` / `sandbox.sandbox_id` | Property is `sandboxId`/`sandbox_id` |
| `Sandbox.get(id)` | `Sandbox.connect(id)` | Method is `connect`, not `get` |
| `sandbox.destroy()` | `sandbox.kill()` | Method is `kill`, not `destroy` |

## File Operations

| ❌ Wrong | ✅ Correct | Why |
|----------|-----------|-----|
| `sandbox.files.write([items])` (Python batch) | `sandbox.files.write_files([items])` | Python batch write is `write_files`, not `write` |
| `sandbox.files.mkdir('/path')` | `sandbox.files.makeDir('/path')` / `sandbox.files.make_dir('/path')` | Method is `makeDir`/`make_dir` |
| `sandbox.files.delete('/path')` | `sandbox.files.remove('/path')` | Method is `remove`, not `delete` |

## Git Operations

| ❌ Wrong | ✅ Correct | Why |
|----------|-----------|-----|
| `sandbox.git.clone(url, path)` | `sandbox.git.clone(url, { path })` (JS) / `sandbox.git.clone(url, path=path)` (Python) | Path is a named option, not positional |
| `sandbox.git.checkout(branch)` | `sandbox.git.checkoutBranch(path, branch)` | Requires repo path; method is `checkoutBranch`/`checkout_branch` |
| `sandbox.git.authenticate(...)` | `sandbox.git.dangerouslyAuthenticate(...)` / `sandbox.git.dangerously_authenticate(...)` | Method includes `dangerously` prefix as a security reminder |

## Templates

| ❌ Wrong | ✅ Correct | Why |
|----------|-----------|-----|
| `template.build()` | `Template.build(template, 'name')` | `build` is a static method on `Template`, takes the template object and a name |
| `Template.build(template)` | `Template.build(template, 'name')` | Name is required |
| `template.install('numpy')` | `template.pipInstall(['numpy'])` / `template.pip_install(['numpy'])` | Use specific installer method; takes array |
| `template.run('cmd')` | `template.runCmd('cmd')` / `template.run_cmd('cmd')` | Method is `runCmd`/`run_cmd` |
| `template.setEnvs({})` for runtime | `Sandbox.create({ envs: {} })` | Template envs are build-time only; use `Sandbox.create({ envs })` for runtime |

## Code Interpreter

| ❌ Wrong | ✅ Correct | Why |
|----------|-----------|-----|
| `sandbox.execute(code)` | `sandbox.runCode(code)` / `sandbox.run_code(code)` | Method is `runCode`/`run_code` |
| `result.output` | `execution.text` or `execution.results` | Output is in `.text` (text) or `.results` (rich: png, html, svg) |
| `result.stdout` | `execution.logs.stdout` | Logs are nested under `.logs` |
| `result.error` | `execution.error` (with `.name`, `.value`, `.traceback`) | Error is an object with structured fields |

## MCP

| ❌ Wrong | ✅ Correct | Why |
|----------|-----------|-----|
| `import { Sandbox } from '@e2b/code-interpreter'` (for MCP) | `import Sandbox from 'e2b'` (JS) / `from e2b import AsyncSandbox` (Python) | MCP uses the base `e2b` package |
| `sandbox.mcpUrl` | `sandbox.getMcpUrl()` / `sandbox.get_mcp_url()` | It's a method, not a property |
| `sandbox.mcpToken` | `sandbox.getMcpToken()` / `sandbox.get_mcp_token()` | It's a method, not a property |
