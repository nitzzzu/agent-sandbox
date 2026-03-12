# E2B SDK API Reference

Complete method listing for all SDK namespaces. JS methods use camelCase, Python uses snake_case.

## Sandbox Lifecycle

| Operation | JavaScript/TypeScript | Python |
|-----------|----------------------|--------|
| Create | `await Sandbox.create({ timeoutMs?, template?, metadata?, envs?, network?, secure? })` | `Sandbox.create(timeout=, template=, metadata=, envs=, network=, secure=)` |
| Connect | `await Sandbox.connect(sandboxId, { timeoutMs? })` | `Sandbox.connect(sandbox_id, timeout=)` |
| Kill | `await sandbox.kill()` | `sandbox.kill()` |
| Set timeout | `await sandbox.setTimeout(ms)` | `sandbox.set_timeout(seconds)` |
| Get info | `await sandbox.getInfo()` | `sandbox.get_info()` |
| Is running | `sandbox.isRunning()` | `sandbox.is_running()` |
| List | `await Sandbox.list({ query? })` | `Sandbox.list(query=)` |
| Pause (beta) | `await sandbox.betaPause()` | `sandbox.beta_pause()` |
| Create paused (beta) | `await Sandbox.betaCreate({ autoPause?, timeoutMs? })` | `Sandbox.beta_create(auto_pause=, timeout=)` |
| Sandbox ID | `sandbox.sandboxId` | `sandbox.sandbox_id` |

## Commands (`sandbox.commands`)

| Operation | JavaScript/TypeScript | Python |
|-----------|----------------------|--------|
| Run | `await sandbox.commands.run('cmd', { background?, cwd?, envs?, user?, onStdout?, onStderr?, timeoutMs? })` | `sandbox.commands.run('cmd', background=, cwd=, envs=, user=, on_stdout=, on_stderr=, timeout=)` |
| List | `await sandbox.commands.list()` | `sandbox.commands.list()` |
| Kill | `await sandbox.commands.kill(pid)` | `sandbox.commands.kill(pid)` |
| Send stdin | `await process.sendStdin('input')` | `process.send_stdin('input')` |
| Connect | `await sandbox.commands.connect(pid, { onStdout?, onStderr? })` | `sandbox.commands.connect(pid, on_stdout=, on_stderr=)` |

## Files (`sandbox.files`)

| Operation | JavaScript/TypeScript | Python |
|-----------|----------------------|--------|
| Read | `await sandbox.files.read('/path')` | `sandbox.files.read('/path')` |
| Write (single) | `await sandbox.files.write('/path', 'content')` | `sandbox.files.write('/path', 'content')` |
| Write (batch) | `await sandbox.files.write([{ path, data }])` | `sandbox.files.write_files([{"path": ..., "data": ...}])` |
| List | `await sandbox.files.list('/dir')` | `sandbox.files.list('/dir')` |
| Make dir | `await sandbox.files.makeDir('/path')` | `sandbox.files.make_dir('/path')` |
| Remove | `await sandbox.files.remove('/path')` | `sandbox.files.remove('/path')` |
| Rename | `await sandbox.files.rename('/old', '/new')` | `sandbox.files.rename('/old', '/new')` |
| Watch | `await sandbox.files.watch('/dir', callback)` | `sandbox.files.watch('/dir', callback)` |
| Get info | `await sandbox.files.getInfo('/path')` | `sandbox.files.get_info('/path')` |

## Git (`sandbox.git`)

| Operation | JavaScript/TypeScript | Python |
|-----------|----------------------|--------|
| Clone | `await sandbox.git.clone(url, { path, branch?, depth?, username?, password? })` | `sandbox.git.clone(url, path=, branch=, depth=, username=, password=)` |
| Init | `await sandbox.git.init(path)` | `sandbox.git.init(path)` |
| Add | `await sandbox.git.add(path, { files? })` | `sandbox.git.add(path, files=)` |
| Commit | `await sandbox.git.commit(path, message, { authorName?, authorEmail?, allowEmpty? })` | `sandbox.git.commit(path, message, author_name=, author_email=, allow_empty=)` |
| Push | `await sandbox.git.push(path, { remote?, branch?, setUpstream?, username?, password? })` | `sandbox.git.push(path, remote=, branch=, set_upstream=, username=, password=)` |
| Pull | `await sandbox.git.pull(path, { remote?, branch?, username?, password? })` | `sandbox.git.pull(path, remote=, branch=, username=, password=)` |
| Status | `await sandbox.git.status(path)` | `sandbox.git.status(path)` |
| Branches | `await sandbox.git.branches(path)` | `sandbox.git.branches(path)` |
| Create branch | `await sandbox.git.createBranch(path, name)` | `sandbox.git.create_branch(path, name)` |
| Checkout | `await sandbox.git.checkoutBranch(path, name)` | `sandbox.git.checkout_branch(path, name)` |
| Delete branch | `await sandbox.git.deleteBranch(path, name, { force? })` | `sandbox.git.delete_branch(path, name, force=)` |
| Configure user | `await sandbox.git.configureUser(name, email, { scope?, path? })` | `sandbox.git.configure_user(name, email, scope=, path=)` |
| Auth helper | `await sandbox.git.dangerouslyAuthenticate({ username, password, host?, protocol? })` | `sandbox.git.dangerously_authenticate(username=, password=, host=, protocol=)` |
| Set config | `await sandbox.git.setConfig(key, value, { scope?, path? })` | `sandbox.git.set_config(key, value, scope=, path=)` |
| Get config | `await sandbox.git.getConfig(key, { scope?, path? })` | `sandbox.git.get_config(key, scope=, path=)` |
| Add remote | `await sandbox.git.remoteAdd(path, name, url, { fetch?, overwrite? })` | `sandbox.git.remote_add(path, name, url, fetch=, overwrite=)` |

## Networking

| Operation | JavaScript/TypeScript | Python |
|-----------|----------------------|--------|
| Get host | `sandbox.getHost(port)` | `sandbox.get_host(port)` |
| Traffic token | `sandbox.trafficAccessToken` | `sandbox.traffic_access_token` |
| Download URL | `sandbox.downloadUrl` | `sandbox.download_url` |
| Upload URL | `sandbox.uploadUrl` | `sandbox.upload_url` |

## Code Interpreter (requires `@e2b/code-interpreter` / `e2b_code_interpreter`)

| Operation | JavaScript/TypeScript | Python |
|-----------|----------------------|--------|
| Run code | `await sandbox.runCode(code, { language?, onStdout?, onStderr? })` | `sandbox.run_code(code, language=, on_stdout=, on_stderr=)` |
| Create context | `await sandbox.createCodeContext()` | `sandbox.create_code_context()` |
| Run in context | `await sandbox.runCode(code, { context })` | `sandbox.run_code(code, context=)` |

## Templates

| Operation | JavaScript/TypeScript | Python |
|-----------|----------------------|--------|
| Build | `await Template.build(template, 'name', { cpuCount?, memoryMB?, tags?, skipCache?, onBuildLogs? })` | `Template.build(template, 'name', cpu_count=, memory_mb=, tags=, skip_cache=, on_build_logs=)` |
| Build (background) | `await Template.buildInBackground(template, 'name', opts)` | `Template.build_in_background(template, 'name', ...)` |
| Build status | `await Template.getBuildStatus(buildInfo, { logsOffset? })` | `Template.get_build_status(build_info, logs_offset=)` |
| Assign tags | `await Template.assignTags('name:tag', 'newTag')` | `Template.assign_tags('name:tag', 'newTag')` |
| Remove tags | `await Template.removeTags('name', 'tag')` | `Template.remove_tags('name', 'tag')` |
| Check exists | `await Template.exists('name')` | `Template.exists('name')` |

## MCP Gateway (requires base `e2b` package)

| Operation | JavaScript/TypeScript | Python |
|-----------|----------------------|--------|
| Get MCP URL | `sandbox.getMcpUrl()` | `sandbox.get_mcp_url()` |
| Get MCP token | `sandbox.getMcpToken()` | `sandbox.get_mcp_token()` |
