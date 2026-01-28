import os

# for no https and forward to sandbox by path vars
def local():
    os.environ['E2B_DEBUG'] = "true"
    os.environ['E2B_API_URL'] = 'http://localhost:10000/e2b/v1'
    os.environ['E2B_API_KEY'] = 'testuser-aef134ef-7aa1-945e-9399-7df9a4ad0c3f'
    os.environ['E2B_DOMAIN'] = 'localhost:10000'

    def __connection_config_get_host(_, sandbox_id: str, sandbox_domain: str, port: int) -> str:
        return f"{sandbox_domain}/sandboxes/router/{sandbox_id}/{port}"
    from e2b import ConnectionConfig
    ConnectionConfig.get_host = __connection_config_get_host

# for no https
def dev():
    os.environ['E2B_DEBUG'] = "true"
    os.environ['E2B_API_KEY'] = 'testuser-aef134ef-7aa1-945e-9399-7df9a4ad0c3f'
    os.environ['E2B_DOMAIN'] = 'example.domain.com'
    os.environ['E2B_API_URL'] = 'http://example.domain.com/e2b/v1'

    def __connection_config_get_host(_, sandbox_id: str, sandbox_domain: str, port: int) -> str:
        return f"{port}-{sandbox_id}.{sandbox_domain}"
    from e2b import ConnectionConfig
    ConnectionConfig.get_host = __connection_config_get_host

