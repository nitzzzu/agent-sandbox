from e2b_desktop import Sandbox

import e2b_hack
e2b_hack.local()

test_id = "e7060cf98e68427faaf32649a60ab5ad"

def create_test():
    # Create a new desktop sandbox
    desktop = Sandbox.create()
    print('Desktop sandbox created', desktop.sandbox_id)

    # Stream the application's window
    # Note: there can be only one stream at a time
    # You need to stop the current stream before streaming another application
    print('Starting to stream Google Chrome')
    desktop.stream.start()

    # Print the stream URL
    print('Stream URL:', desktop.stream.get_url().replace('https:', 'http:'))


    # Launch an application
    print('Launching Google Chrome')
    desktop.wait(10000)
    desktop.launch('google-chrome')  # or vscode, firefox, etc.

    # Wait 15s for the application to open
    desktop.wait(10000)

    # Do some actions in the application
    print('Writing to Google Chrome')
    desktop.press('esc')
    desktop.wait(2000)
    desktop.write('https://infoq.com')

    print('Pressing Enter')
    desktop.press('enter')

    # wait 15s for page to load
    print('Waiting 15s')
    desktop.wait(15000)

    # Stop the stream
    # print('Stopping the stream')
    # desktop.stream.stop()

    # Open another application
    print('Launching VS Code')
    desktop.launch('code')

    # Wait 15s for the application to open
    desktop.wait(5000)

    # Kill the sandbox after the tasks are finished
    # desktop.kill()


def launch_test():
    desktop=Sandbox.connect(test_id)

    # original_run = with_display_env(desktop.commands.run)
    # desktop.commands.run = original_run

    print(f"sandbox id: {desktop.sandbox_id}")

    # Launch an application
    desktop.launch('google-chrome')  # or vscode, firefox, etc.

    # Wait 10s for the application to open
    desktop.wait(5000)

    # Do some actions in the application
    desktop.press('esc') # close update prompt

    print('Writing to Google Chrome')
    desktop.write('https://baidu.com')
    w = desktop.get_application_windows('google-chrome')
    print("windows", w)

    print('Pressing Enter')
    desktop.press('enter')


    # Open another application
    print('Launching VS Code')
    desktop.launch('code')

    # Wait 15s for the application to open
    desktop.wait(15000)



def get_test():
    sandbox=Sandbox.connect(test_id)
    wid=sandbox.get_current_window_id()
    print(f"sandbox id: {wid}")

    wtitle=sandbox.get_window_title(wid)
    print(f"sandbox title: {wtitle}")

    info = sandbox.get_info()
    print(f"sandbox info: {info}")


    image_byte_array = sandbox.screenshot(format= "bytes")
    # save bytearray to image
    with open("desktop.png", "wb") as f:
        f.write(image_byte_array)



def port_test():
    import requests

    sandbox=Sandbox.connect(test_id)

    # Start a server inside the sandbox
    sandbox.commands.run("python3 -m http.server 8080", background=True)

    host = sandbox.get_host(8080)
    url = f"http://{host}"

    # Request without token will fail with 403
    response1 = requests.get(url)
    print(response1.text)  # 403

def upload_file():
    sbx=Sandbox.connect(test_id)
    print(f"sandbox id: {sbx.sandbox_id}")

    # Read file from local filesystem
    with open("desktop.py", "rb") as file:
        # Upload file to sandbox
        sbx.files.write("/home/user/desktop.py", file)


if __name__ == "__main__":
    create_test()