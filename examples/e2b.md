# E2B Examples
In this directory, you can find examples demonstrating how to use Agent-Sandbox with E2B protocol compatibility, including code interpreter and desktop environments.

## Before You Start
### 1, In Development
When developing and testing locally, you can use the following environment variable configurations.

E2B default sandbox URL format is as follows:  
`https://{port}-{sandboxID}.your-domain.com`   
e.g.  
`https://6080-294bef011f1e4567b4c5d02593e2e90e.example.com`

‼️ If you don't have `https` or Wildcard Domain(*.example.com) support, please config and hack by the following two functions in `local()` and `dev()` to adapt to your environment.

```python
import os

# for no https and no wildcard domain support
def local():
    os.environ['E2B_DEBUG'] = "true"
    os.environ['E2B_API_URL'] = 'http://localhost:10000/e2b/v1'
    os.environ['E2B_API_KEY'] = 'testuser-aef134ef-7aa1-945e-9399-7df9a4ad0c3f'
    os.environ['E2B_DOMAIN'] = 'localhost:10000'

    def __connection_config_get_host(_, sandbox_id: str, sandbox_domain: str, port: int) -> str:
        return f"{sandbox_domain}/sandboxes/router/{sandbox_id}/{port}"
    from e2b import ConnectionConfig
    ConnectionConfig.get_host = __connection_config_get_host


# for no https, but wildcard domain support
def dev():
    os.environ['E2B_DEBUG'] = "true"
    os.environ['E2B_API_KEY'] = 'testuser-aef134ef-7aa1-945e-9399-7df9a4ad0c3f'
    os.environ['E2B_DOMAIN'] = 'example.domain.com'
    os.environ['E2B_API_URL'] = 'http://example.domain.com/e2b/v1'

    def __connection_config_get_host(_, sandbox_id: str, sandbox_domain: str, port: int) -> str:
        return f"{port}-{sandbox_id}.{sandbox_domain}"
    from e2b import ConnectionConfig
    ConnectionConfig.get_host = __connection_config_get_host
```
---
### 2, In Production
Normally in production, you can set the following environment variables to configure E2B SDK:

```python
import os
os.environ['E2B_DOMAIN'] = 'example.domain.com'
os.environ['E2B_API_URL'] = 'https://example.domain.com/e2b/v1'
os.environ['E2B_API_KEY'] = 'testuser-aef134ef-7aa1-945e-9399-7df9a4ad0c3f'
```

Agent-Sandbox ingress should also be configured to support wildcard domain and https.
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: agent-sandbox
  annotations:
    nginx.ingress.kubernetes.io/proxy-body-size: 1024M
    nginx.ingress.kubernetes.io/proxy-connect-timeout: "300"
    nginx.ingress.kubernetes.io/proxy-read-timeout: "300"
    nginx.ingress.kubernetes.io/proxy-send-timeout: "300"
spec:
  ingressClassName: external-ingress-controller
  rules:
    - host: "*.example.domain.com"
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


## More Information

- [E2B Official Documentation](https://docs.e2b.dev/)
- [e2b-code-interpreter SDK](https://github.com/e2b-dev/e2b-code-interpreter)
- [E2B Desktop SDK](https://github.com/e2b-dev/desktop)