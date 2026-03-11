# Moltbot Sandbox 集成指南

本文档说明如何为 Moltbot 提供 Sandbox 能力，包括两种方案：
1. 使用 Moltbot 内置的 Docker Sandbox 功能
2. 使用 Agent-Sandbox 作为后端（推荐，企业级方案）

## 方案一：Moltbot 内置 Docker Sandbox

Moltbot 支持在 Docker 容器中运行工具来提供沙箱隔离。这是 Moltbot 的原生功能。

### 1. 构建 Sandbox 镜像

Moltbot 默认使用 `moltbot-sandbox:bookworm-slim` 镜像。你需要先构建这个镜像：

```bash
# 克隆 Moltbot 仓库（如果还没有）
git clone https://github.com/moltbot/moltbot
cd moltbot

# 运行构建脚本
scripts/sandbox-setup.sh
```

如果需要浏览器支持，还需要构建浏览器镜像：

```bash
scripts/sandbox-browser-setup.sh
```

### 2. 配置 Moltbot

在 Moltbot 配置文件中启用 sandbox：

```json5
{
  agents: {
    defaults: {
      sandbox: {
        mode: "non-main",  // 或 "all" 来沙箱化所有会话
        scope: "session",  // 每个会话一个容器
        workspaceAccess: "none",  // 或 "ro"/"rw"
        docker: {
          network: "none",  // 默认无网络，需要时改为 "bridge"
          readOnlyRoot: false,  // 如果需要写入
          user: "0:0",  // root 用户（用于安装包）
          env: {
            // 环境变量，如 API keys
            "OPENAI_API_KEY": "your-key"
          },
          setupCommand: "apt-get update && apt-get install -y python3-pip",  // 一次性设置命令
          binds: [
            "/home/user/source:/source:ro",  // 只读挂载
            "/var/run/docker.sock:/var/run/docker.sock:ro"  // Docker socket（只读）
          ]
        },
        browser: {
          autoStart: true,
          autoStartTimeoutMs: 30000,
          allowHostControl: false
        }
      }
    }
  }
}
```

### 3. 配置说明

- **mode**: 
  - `"off"`: 不启用沙箱
  - `"non-main"`: 只沙箱化非主会话（默认推荐）
  - `"all"`: 所有会话都沙箱化

- **scope**:
  - `"session"`: 每个会话一个容器（推荐，完全隔离）
  - `"agent"`: 每个 agent 一个容器
  - `"shared"`: 所有会话共享一个容器

- **workspaceAccess**:
  - `"none"`: 沙箱使用独立工作空间（默认）
  - `"ro"`: 只读挂载 agent 工作空间到 `/agent`
  - `"rw"`: 读写挂载到 `/workspace`

### 4. 限制和注意事项

- Moltbot 的 sandbox 是**工具执行级别的隔离**，不是完整的沙箱服务
- 需要手动管理 Docker 镜像和容器
- 不支持多租户、状态持久化等企业级特性
- 配置相对复杂，需要了解 Docker 和容器网络

---

## 方案二：使用 Agent-Sandbox 作为后端（推荐）

Agent-Sandbox 提供了企业级的沙箱服务，可以作为 Moltbot 的后端，提供更强大的能力。

### 1. 部署 Agent-Sandbox

```bash
# 创建命名空间
kubectl create namespace agent-sandbox

# 部署 Agent-Sandbox
kubectl apply -n agent-sandbox -f install.yaml

# 配置 Ingress（可选）
kubectl apply -n agent-sandbox -f dev/ingress.yaml
```

### 2. 配置 Moltbot 使用 Agent-Sandbox

有两种集成方式：

#### 方式 A：通过 MCP Server 集成

Agent-Sandbox 提供 MCP Server，Moltbot 可以通过 MCP 协议使用：

```json5
{
  agents: {
    defaults: {
      // 禁用 Moltbot 内置 sandbox
      sandbox: {
        mode: "off"
      },
      // 配置 MCP 服务器连接到 Agent-Sandbox
      mcp: {
        servers: {
          "agent-sandbox": {
            url: "http://agent-sandbox.your-host.com/mcp",
            transport: "sse"  // 或 "http"
          }
        }
      }
    }
  }
}
```

这样 Moltbot 可以通过 MCP 工具调用 Agent-Sandbox 的功能：
- 创建沙箱
- 执行代码
- 浏览器自动化
- 文件操作
- 删除沙箱

#### 方式 B：通过 E2B 协议集成

Agent-Sandbox 完全兼容 E2B 协议，如果你的 Moltbot 支持 E2B，可以直接配置：

```json5
{
  agents: {
    defaults: {
      sandbox: {
        mode: "off"  // 禁用内置 sandbox
      },
      // 如果 Moltbot 支持 E2B 后端配置
      e2b: {
        apiUrl: "http://agent-sandbox.your-host.com/e2b/v1",
        apiKey: "your-api-key",
        domain: "your-host.com"
      }
    }
  }
}
```

### 3. 使用 Agent-Sandbox 的优势

相比 Moltbot 内置 sandbox，Agent-Sandbox 提供：

✅ **企业级特性**：
- 多租户支持（每个 API Key 一个租户）
- 多会话管理（同时运行多个沙箱）
- 状态持久化（沙箱状态可持久化存储）
- 资源控制（精确控制 CPU、内存、磁盘）

✅ **更好的隔离**：
- 每个会话完全隔离的容器
- 支持自定义镜像和模板
- 支持网络隔离和端口转发

✅ **易于管理**：
- RESTful API 和 MCP 接口
- 自动生命周期管理
- 支持 TTL 和空闲回收

✅ **生产就绪**：
- 基于 Kubernetes，支持自动扩缩容
- 高可用和容错
- 完整的监控和日志

### 4. 示例：Moltbot + Agent-Sandbox 工作流

```python
# Moltbot 通过 MCP 调用 Agent-Sandbox
# 1. Agent 需要执行代码时，自动创建沙箱
# 2. 在沙箱中执行代码
# 3. 获取结果
# 4. 自动清理沙箱

# 或者直接使用 E2B SDK（如果 Moltbot 支持）
from e2b_code_interpreter import Sandbox
import os

os.environ['E2B_API_URL'] = 'http://agent-sandbox.your-host.com/e2b/v1'
os.environ['E2B_API_KEY'] = 'your-api-key'
os.environ['E2B_DOMAIN'] = 'your-host.com'

sbx = Sandbox.create()
execution = sbx.run_code("print('Hello from Agent-Sandbox!')")
print(execution.logs)
sbx.close()
```

---

## 方案对比

| 特性 | Moltbot 内置 Sandbox | Agent-Sandbox |
|------|---------------------|---------------|
| 部署复杂度 | 中等（需要构建镜像） | 低（Kubernetes 一键部署） |
| 隔离级别 | 工具执行级别 | 完整容器隔离 |
| 多租户 | ❌ | ✅ |
| 状态持久化 | ❌ | ✅ |
| 资源控制 | 基础 | 完整（CPU/内存/磁盘） |
| 生命周期管理 | 手动 | 自动（TTL、空闲回收） |
| 扩展性 | 有限 | 高（Kubernetes 自动扩缩容） |
| 企业级特性 | ❌ | ✅ |
| 生产就绪 | ⚠️ | ✅ |

---

## 推荐方案

**对于生产环境**：推荐使用 **Agent-Sandbox**，因为它提供了：
- 更好的隔离和安全性
- 企业级的多租户和资源管理
- 自动化的生命周期管理
- 生产级的可靠性和扩展性

**对于简单场景**：如果只是本地开发或测试，可以使用 Moltbot 内置的 Docker Sandbox。

---

## 相关资源

- [Moltbot Sandbox 文档](https://docs.molt.bot/gateway/sandboxing)
- [Agent-Sandbox GitHub](https://github.com/agent-sandbox/agent-sandbox)
- [Agent-Sandbox E2B 协议支持](./e2b-protocol-support.md)
- [Agent-Sandbox 使用示例](../examples/)

---

## 下一步

1. 根据你的需求选择合适的方案
2. 部署 Agent-Sandbox（如果选择方案二）
3. 配置 Moltbot 连接到 Agent-Sandbox
4. 测试沙箱功能
5. 根据实际使用情况调整配置

如有问题，请参考 [Agent-Sandbox 文档](../README.md) 或提交 Issue。

