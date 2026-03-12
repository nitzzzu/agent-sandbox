---
name: e2b-code-interpreter
description: >
  Execute code in E2B sandboxes and integrate with LLMs for tool calling.
  Use when building AI agents that need to run Python/JS code, analyze data,
  generate charts, or use LLM function calling with E2B.
---

# E2B Code Interpreter — Code Execution & LLM Integration

## Setup

```bash
# JavaScript/TypeScript
npm install @e2b/code-interpreter

# Python
pip install e2b-code-interpreter
```

```typescript
import { Sandbox } from '@e2b/code-interpreter'

const sandbox = await Sandbox.create()
```

```python
from e2b_code_interpreter import Sandbox

sandbox = Sandbox.create()
```

## Running Code

```typescript
const execution = await sandbox.runCode('print("Hello, World!")')

console.log(execution.text)         // "Hello, World!"
console.log(execution.logs.stdout)  // ["Hello, World!\n"]
console.log(execution.logs.stderr)  // []
console.log(execution.error)        // null (or { name, value, traceback })
```

```python
execution = sandbox.run_code('print("Hello, World!")')

print(execution.text)          # "Hello, World!"
print(execution.logs.stdout)   # ["Hello, World!\n"]
print(execution.logs.stderr)   # []
print(execution.error)         # None (or object with name, value, traceback)
```

### Execution Result Structure

| Field | Type | Description |
|-------|------|-------------|
| `execution.text` | `string` | Last text output |
| `execution.results` | `array` | Rich outputs (`.png`, `.html`, `.svg`, `.json`, `.text`) |
| `execution.logs.stdout` | `string[]` | Stdout lines |
| `execution.logs.stderr` | `string[]` | Stderr lines |
| `execution.error` | `object \| null` | Error with `.name`, `.value`, `.traceback` |

### Streaming Output

```typescript
const execution = await sandbox.runCode('for i in range(5): print(i)', {
  onStdout: (line) => console.log('stdout:', line),
  onStderr: (line) => console.error('stderr:', line),
})
```

```python
execution = sandbox.run_code(
    "for i in range(5): print(i)",
    on_stdout=lambda line: print("stdout:", line),
    on_stderr=lambda line: print("stderr:", line),
)
```

### Language Support

Default is Python. Specify other languages with the `language` option:

```typescript
// Python (default)
await sandbox.runCode('print("Python")')

// JavaScript
await sandbox.runCode('console.log("JavaScript")', { language: 'javascript' })

// R
await sandbox.runCode('cat("R language")', { language: 'r' })

// Java
await sandbox.runCode('System.out.println("Java")', { language: 'java' })

// Bash
await sandbox.runCode('echo "Bash"', { language: 'bash' })
```

```python
sandbox.run_code('print("Python")')
sandbox.run_code('console.log("JavaScript")', language="javascript")
sandbox.run_code('cat("R language")', language="r")
```

## Execution Contexts (Shared State)

By default each `runCode` call is independent. Use contexts to share state across calls:

```typescript
const context = await sandbox.createCodeContext()

await sandbox.runCode('x = 42', { context })
const execution = await sandbox.runCode('print(x)', { context })
console.log(execution.text) // "42"
```

```python
context = sandbox.create_code_context()

sandbox.run_code("x = 42", context=context)
execution = sandbox.run_code("print(x)", context=context)
print(execution.text)  # "42"
```

## Charts & Visualizations

For matplotlib charts, the code must call `display()` on the figure:

```python
# Code to execute in sandbox
code = """
import matplotlib.pyplot as plt
import numpy as np

x = np.linspace(0, 10, 100)
plt.figure(figsize=(10, 6))
plt.plot(x, np.sin(x))
plt.title('Sine Wave')
display(plt.gcf())
"""
```

```typescript
const execution = await sandbox.runCode(code)
// Chart is in execution.results[0].png (base64-encoded)
const chartBase64 = execution.results[0].png
```

```python
execution = sandbox.run_code(code)
chart_base64 = execution.results[0].png
```

## LLM Tool Calling Pattern

The standard pattern for connecting LLMs to E2B:

1. Define a tool/function that executes code in an E2B sandbox
2. Send the tool definition to the LLM
3. When the LLM calls the tool, execute the code in the sandbox
4. Return execution results to the LLM for interpretation

### OpenAI Function Calling

```typescript
import OpenAI from 'openai'
import { Sandbox } from '@e2b/code-interpreter'

const openai = new OpenAI()
const sandbox = await Sandbox.create()

const tools = [{
  type: 'function',
  function: {
    name: 'execute_python',
    description: 'Execute Python code in a sandbox',
    parameters: {
      type: 'object',
      properties: {
        code: { type: 'string', description: 'Python code to execute' },
      },
      required: ['code'],
    },
  },
}]

const response = await openai.chat.completions.create({
  model: 'gpt-4o',
  messages: [{ role: 'user', content: 'Calculate fibonacci of 10' }],
  tools,
})

// Handle tool call
const toolCall = response.choices[0].message.tool_calls?.[0]
if (toolCall) {
  const { code } = JSON.parse(toolCall.function.arguments)
  const execution = await sandbox.runCode(code)
  console.log(execution.text)
}
```

```python
from openai import OpenAI
from e2b_code_interpreter import Sandbox

client = OpenAI()
sandbox = Sandbox.create()

tools = [{
    "type": "function",
    "function": {
        "name": "execute_python",
        "description": "Execute Python code in a sandbox",
        "parameters": {
            "type": "object",
            "properties": {
                "code": {"type": "string", "description": "Python code to execute"},
            },
            "required": ["code"],
        },
    },
}]

response = client.chat.completions.create(
    model="gpt-4o",
    messages=[{"role": "user", "content": "Calculate fibonacci of 10"}],
    tools=tools,
)

tool_call = response.choices[0].message.tool_calls[0]
if tool_call:
    code = json.loads(tool_call.function.arguments)["code"]
    execution = sandbox.run_code(code)
    print(execution.text)
```

### Anthropic Tool Use

```typescript
import Anthropic from '@anthropic-ai/sdk'
import { Sandbox } from '@e2b/code-interpreter'

const anthropic = new Anthropic()
const sandbox = await Sandbox.create()

const response = await anthropic.messages.create({
  model: 'claude-sonnet-4-5-20250929',
  max_tokens: 1024,
  tools: [{
    name: 'execute_python',
    description: 'Execute Python code in a sandbox',
    input_schema: {
      type: 'object',
      properties: {
        code: { type: 'string', description: 'Python code to execute' },
      },
      required: ['code'],
    },
  }],
  messages: [{ role: 'user', content: 'Calculate the first 20 primes' }],
})

// Handle tool use
for (const block of response.content) {
  if (block.type === 'tool_use') {
    const execution = await sandbox.runCode(block.input.code)
    console.log(execution.text)
  }
}
```

```python
import anthropic
from e2b_code_interpreter import Sandbox

client = anthropic.Anthropic()
sandbox = Sandbox.create()

response = client.messages.create(
    model="claude-sonnet-4-5-20250929",
    max_tokens=1024,
    tools=[{
        "name": "execute_python",
        "description": "Execute Python code in a sandbox",
        "input_schema": {
            "type": "object",
            "properties": {
                "code": {"type": "string", "description": "Python code to execute"},
            },
            "required": ["code"],
        },
    }],
    messages=[{"role": "user", "content": "Calculate the first 20 primes"}],
)

for block in response.content:
    if block.type == "tool_use":
        execution = sandbox.run_code(block.input["code"])
        print(execution.text)
```

## Data Analysis Workflow

Upload data, prompt the LLM to generate analysis code, execute in sandbox, extract results.

```typescript
import { Sandbox } from '@e2b/code-interpreter'

const sandbox = await Sandbox.create()

// 1. Upload data
await sandbox.files.write('/home/user/data.csv', csvContent)

// 2. Run analysis code (generated by LLM or hand-written)
const execution = await sandbox.runCode(`
import pandas as pd
import matplotlib.pyplot as plt

df = pd.read_csv('/home/user/data.csv')
print(df.describe())

plt.figure(figsize=(10, 6))
df.plot(kind='bar')
plt.tight_layout()
display(plt.gcf())
`)

// 3. Get results
console.log(execution.text)           // Statistical summary
const chart = execution.results[0].png // Base64 chart image
```

```python
from e2b_code_interpreter import Sandbox

sandbox = Sandbox.create()

# 1. Upload data
sandbox.files.write("/home/user/data.csv", csv_content)

# 2. Run analysis
execution = sandbox.run_code("""
import pandas as pd
import matplotlib.pyplot as plt

df = pd.read_csv('/home/user/data.csv')
print(df.describe())

plt.figure(figsize=(10, 6))
df.plot(kind='bar')
plt.tight_layout()
display(plt.gcf())
""")

# 3. Get results
print(execution.text)
chart = execution.results[0].png
```

## Python Context Manager

```python
from e2b_code_interpreter import Sandbox

# Sandbox is automatically killed when the block exits
with Sandbox.create() as sandbox:
    execution = sandbox.run_code("print('Hello')")
    print(execution.text)
```