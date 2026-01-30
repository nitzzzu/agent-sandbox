# 🚀 Manus的Sandbox我们可以轻松拥有，Agent-Sandbox 完全支持 E2B 协议：轻松本地部署

![e2bdesktop.gif](e2bdesktop.gif)


对于 AI Agent 开发者来说，E2B（[https://e2b.dev/](https://e2b.dev/)）是一个非常强大的沙箱执行环境。它提供了代码执行、浏览器自动化、桌面应用控制等功能，让 AI Agent 能够安全地执行各种任务。然而，E2B 作为云服务，对于需要本地化部署、数据安全控制或成本优化的企业来说，一直是一个痛点，Manus底层能力就是靠它。

**现在，这个痛点彻底解决了！** Agent-Sandbox 实现了对 E2B 协议的**完全兼容**，让你可以：

- ✅ **零代码改动**：直接使用 E2B 官方 SDK（`e2b-code-interpreter`、`e2b-desktop` 等）
- ✅ **本地化部署**：在自己的 Kubernetes 集群中运行，完全掌控数据和安全，部署非常简单，无需复杂配置

## 什么是 E2B？

E2B 是一个为 AI Agent 设计的云原生执行环境，提供了：

- **代码执行沙箱**：安全执行 Python、Node.js 等代码
- **浏览器自动化**：完整的浏览器控制能力
- **桌面应用控制**：启动和控制桌面应用程序
- **文件系统操作**：完整的文件读写能力
- **端口转发**：访问沙箱内运行的服务

E2B 通过统一的 API 和 SDK，让 AI Agent 能够轻松创建、管理和使用这些执行环境。而且E2B的生态很好，有很多文章和资源可以参考。

## Agent-Sandbox 的 E2B 协议兼容实现

[Agent-Sandbox](https://github.com/agent-sandbox/agent-sandbox) 完全实现了 E2B 的 API 协议，包括：

### 核心 API 端点

- `POST /e2b/v1/sandboxes` - 创建沙箱
- `GET /e2b/v1/sandboxes/{sandboxID}` - 获取沙箱信息
- `GET /e2b/v1/v2/sandboxes` - 列出所有沙箱
- `DELETE /e2b/v1/sandboxes/{sandboxID}` - 删除沙箱
- `POST /e2b/v1/sandboxes/{sandboxID}/connect` - 连接沙箱

### 支持本地路由

Agent-Sandbox 实现了两种路由方式，完美兼容 E2B SDK 的连接需求：

1. **路径路由模式**（适用于无 HTTPS 和通配符域名的环境）：
   ```
   http://your-domain.com/sandboxes/router/{sandboxID}/{port}/
   ```

2. **原生路由模式**（适用于有通配符域名的生产环境）：
   ```
   https://{port}-{sandboxID}.your-domain.com/
   ```

这种设计让 Agent-Sandbox 能够适应各种部署环境，从本地开发到生产环境都能无缝工作。

## 快速开始：3 步实现本地化部署

### 步骤 1：部署 Agent-Sandbox

在你的 Kubernetes 集群中部署 Agent-Sandbox（需要 Kubernetes 1.26+）：

```bash
kubectl create namespace agent-sandbox
kubectl apply -n agent-sandbox -f install.yaml
```

### 步骤 2：配置 Ingress（可选）

如果需要从外部访问，配置 Ingress：

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: agent-sandbox
  namespace: agent-sandbox
  annotations:
    nginx.ingress.kubernetes.io/proxy-body-size: 1024M
spec:
  ingressClassName: nginx
  rules:
    - host: "*.your-domain.com"  # 通配符域名支持
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

### 步骤 3：使用 E2B SDK

现在你可以直接使用 E2B 的官方 SDK，无需任何修改！

#### 示例 1：代码执行沙箱

```python
from e2b_code_interpreter import Sandbox
import os

# 配置 Agent-Sandbox 端点
os.environ['E2B_API_URL'] = 'http://your-domain.com/e2b/v1'
os.environ['E2B_API_KEY'] = 'testuser-aef134ef-7aa1-945e-9399-7df9a4ad0c3f'
os.environ['E2B_DOMAIN'] = 'your-domain.com'

# 创建沙箱（与 E2B 完全相同的 API）
sbx = Sandbox.create()

# 执行代码
execution = sbx.run_code("print('Hello from Agent-Sandbox!')")
print(execution.logs)

# 文件操作
files = sbx.files.list("/")
print(files)

# 上传文件
with open("local_file.txt", "rb") as file:
    sbx.files.write("/home/user/remote_file.txt", file)
```

#### 示例 2：桌面应用控制

```python
from e2b_desktop import Sandbox
import os

# 配置 Agent-Sandbox 端点
os.environ['E2B_API_URL'] = 'http://your-domain.com/e2b/v1'
os.environ['E2B_API_KEY'] = 'testuser-aef134ef-7aa1-945e-9399-7df9a4ad0c3f'
os.environ['E2B_DOMAIN'] = 'your-domain.com'

# 创建桌面沙箱
desktop = Sandbox.create()

# 启动应用
desktop.launch('google-chrome')

# 控制应用
desktop.write('https://google.com')
desktop.press('enter')

# 截图
image = desktop.screenshot(format="bytes")
```

## 实际应用场景

### 场景 1：AI Agent 代码执行

```python
from e2b_code_interpreter import Sandbox

sbx = Sandbox.create()

# Agent 生成的代码可以直接执行
code = """
import pandas as pd
import matplotlib.pyplot as plt

data = pd.DataFrame({'x': [1, 2, 3], 'y': [4, 5, 6]})
plt.plot(data['x'], data['y'])
plt.savefig('/home/user/plot.png')
"""

execution = sbx.run_code(code)
print(execution.logs)
```

### 场景 2：Web 自动化

```python
from e2b_desktop import Sandbox

desktop = Sandbox.create()
desktop.launch('google-chrome')
desktop.write('https://example.com')
desktop.press('enter')
desktop.wait(5000)

# 截图保存结果
screenshot = desktop.screenshot(format="bytes")
with open("result.png", "wb") as f:
    f.write(screenshot)
```

### 场景 3：数据分析工作流

```python
from e2b_code_interpreter import Sandbox

sbx = Sandbox.create()

# 上传数据文件
with open("data.csv", "rb") as file:
    sbx.files.write("/home/user/data.csv", file)

# 执行分析
analysis_code = """
import pandas as pd
df = pd.read_csv('/home/user/data.csv')
print(df.describe())
print(df.head())
"""

execution = sbx.run_code(analysis_code)
print(execution.logs)
```


## 总结

Agent-Sandbox 对 E2B 协议的完全支持，为 AI Agent 开发者带来了简单灵活的本地部署方案：

- 🎯 **零成本迁移**：直接使用 E2B SDK，无需修改代码
- 🏢 **企业级部署**：在自己的基础设施上运行，完全掌控
- 💰 **成本优化**：无需按使用量付费，充分利用现有资源
- 🔒 **数据安全**：数据完全在本地，满足合规要求
- 🚀 **开箱即用**：简单配置即可使用，降低运维成本



## 相关地址

- [Agent-Sandbox GitHub](https://github.com/agent-sandbox/agent-sandbox)
- [E2B 官方文档](https://docs.e2b.dev/)
- [e2b-code-interpreter SDK](https://github.com/e2b-dev/e2b-code-interpreter)
- [E2B Desktop SDK](https://github.com/e2b-dev/desktop)
- [E2B示例代码](https://github.com/agent-sandbox/agent-sandbox/tree/main/examples)


