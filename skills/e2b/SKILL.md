---
name: e2b
description: >
  Core E2B SDK knowledge — correct imports, namespacing, terminology, and API patterns.
  Activates automatically when working with E2B sandboxes, templates, or code execution.
user-invocable: false
---

# E2B SDK — Core Knowledge

E2B provides **sandboxes** — secure, isolated Linux environments purpose-built as an AI agent runtime. Agents use sandboxes to safely execute code, run commands, manage files, and interact with git repositories.

## Two SDK Packages

| Package | JS/TS import | Python import | Use when |
|---------|-------------|---------------|----------|
| Base SDK | `import Sandbox from 'e2b'` | `from e2b import Sandbox` | Templates, MCP gateway, sandbox lifecycle only |
| Code Interpreter | `import { Sandbox } from '@e2b/code-interpreter'` | `from e2b_code_interpreter import Sandbox` | Code execution with `runCode()`/`run_code()` + all base features |

**Default choice:** Use the Code Interpreter package — it includes everything from the base SDK plus code execution. Only use the base `e2b` package when you specifically need MCP gateway or don't need `runCode()`.

## Authentication

Set `E2B_API_KEY` environment variable or pass `apiKey`/`api_key` to SDK methods.

## Critical: Correct Namespacing

All sandbox operations use **namespaced methods**:

| Namespace | Example (JS) | Example (Python) |
|-----------|-------------|------------------|
| `sandbox.commands` | `sandbox.commands.run('ls')` | `sandbox.commands.run('ls')` |
| `sandbox.files` | `sandbox.files.read('/path')` | `sandbox.files.read('/path')` |
| `sandbox.git` | `sandbox.git.clone(url, { path })` | `sandbox.git.clone(url, path=path)` |

**WRONG:** `sandbox.run()`, `sandbox.read()`, `sandbox.clone()` — these do NOT exist.

## Critical: Sandbox Creation

Use `Sandbox.create()` — a static factory method. **NEVER** use `new Sandbox()` or `Sandbox()`.

```typescript
// ✅ Correct
const sandbox = await Sandbox.create()

// ❌ Wrong — constructor does not create a sandbox
const sandbox = new Sandbox()
```

```python
# ✅ Correct
sandbox = Sandbox.create()

# ❌ Wrong — constructor does not create a sandbox
sandbox = Sandbox()
```

## Critical: Timeout Units

| Language | Unit | Parameter name | Example |
|----------|------|---------------|---------|
| JavaScript/TypeScript | **Milliseconds** | `timeoutMs` | `Sandbox.create({ timeoutMs: 60_000 })` |
| Python | **Seconds** | `timeout` | `Sandbox.create(timeout=60)` |

## Terminology

- Use **"sandbox"** or **"E2B sandbox"**
- **NEVER** say "code interpreter sandbox", "ephemeral", "short-lived", or "cloud VM"
- E2B is an **"AI agent runtime"**, not a "cloud provider" or "VM service"
- The main entity is `Sandbox`, templates are `Template`

## Reference Files

- [API Reference](references/api-reference.md) — Complete method listing for all namespaces
- [Common Mistakes](references/common-mistakes.md) — Anti-patterns with corrections
