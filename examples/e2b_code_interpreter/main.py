from e2b_code_interpreter import Sandbox

import e2b_hack
e2b_hack.local()


# sandbox id for tests
test_id = "4943125c14da4de98f423c052c2012ec"

def create_test():
    # Default:
    # TemplateID = "code-interpreter-v1"
    # Timeout = {int} 300
    sbx = Sandbox.create(template="template-demo")

    # Retrieve sandbox information.
    info = sbx.get_info()
    print(info)


    execution = sbx.run_code("print('hello world')")  # Execute Python inside the sandbox
    print(execution.logs)

    files = sbx.files.list("/")
    print(files)


def del_test():
    sbx=Sandbox.connect(test_id)
    # e2b kill skip delete action if debug=true, ref:e2b.sandbox_sync.sandbox_api.SandboxApi._cls_kill
    import os
    os.environ['E2B_DEBUG'] = "false"
    sbx.kill(debug=0)


def list_test():
    paginator = Sandbox.list()

    # Get the first page of sandboxes (running and paused)
    firstPage = paginator.next_items()

    running_sandbox = firstPage[0]

    print('Running sandbox metadata:', running_sandbox.metadata)
    print('Running sandbox id:', running_sandbox.sandbox_id)
    print('Running sandbox started at:', running_sandbox.started_at)
    print('Running sandbox template id:', running_sandbox.template_id)

def get_test():
    sbx=Sandbox.connect(test_id)
    print(f"sandbox id: {sbx.sandbox_id}")

    # Retrieve sandbox information.
    info = sbx.get_info()
    print(info)


def port_test():
    import requests

    sandbox=Sandbox.connect(test_id)

    # Start a server inside the sandbox
    sandbox.commands.run("python -m http.server 8080", background=True)

    host = sandbox.get_host(8080)
    url = f"http://{host}"

    # Request without token will fail with 403
    response1 = requests.get(url)
    print(response1.text)  # 403

def run_code():
    sbx=Sandbox.connect(test_id)
    print(f"sandbox id: {sbx.sandbox_id}")

    # files = sbx.files.list("/home/user")
    # print(files)

    execution = sbx.run_code("print(1+1)")  # Execute Python inside the sandbox
    print(execution.logs)

def connect():
    sbx=Sandbox.connect(test_id)
    print(f"sandbox id: {sbx.sandbox_id}")

    execution = sbx.run_code("print(1+1)")  # Execute Python inside the sandbox
    print(execution.logs)

    files = sbx.files.list("/home/user")
    print(files)

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
        # Use `on_error` to handle runtime code errors
        on_error=lambda error: print('error:', error),
        on_stdout=lambda data: print('stdout:', data),
        on_stderr=lambda data: print('stderr:', data),
    )

    # Read file from local filesystem
    with open("test.md", "rb") as file:
        # Upload file to sandbox
        sbx.files.write("/home/user/test2.md", file)


def upload_file():
    sbx=Sandbox.connect(test_id)
    print(f"sandbox id: {sbx.sandbox_id}")

    # Read file from local filesystem
    with open("test.md", "rb") as file:
        # Upload file to sandbox
        sbx.files.write("/home/user/test2.md", file)





if __name__ == "__main__":
    create_test()