Installation and Quick Start

This section guides you through installing `webautoma` and running your first simple test.

**Prerequisites**

Before you begin, ensure you have the following installed:

1.  **Go:** A recent version of the Go programming language (e.g., 1.18 or later is recommended). You can download it from [https://go.dev/dl/](https://go.dev/dl/). Verify your installation by running `go version` in your terminal.
2.  **WebDriver Executable:** You need the specific WebDriver executable for the browser you want to automate.
    * **ChromeDriver:** (Recommended for initial use with `webautoma`) Download the version that matches your installed Google Chrome browser from [https://chromedriver.chromium.org/downloads](https://chromedriver.chromium.org/downloads).
    * *(Optional: GeckoDriver for Firefox, etc.)*
      Make sure the downloaded WebDriver executable is either placed in a directory included in your system's PATH environment variable or that you know its full path.

**Installation**

1.  **Clone the Repository:** (Replace with the actual URL if you host it)
    ```bash
    git clone <repository_url>
    cd webautoma
    ```

2.  **Build the Executable:**
    ```bash
    go build -o webautoma main.go
    ```
    This command compiles the source code and creates an executable file named `webautoma` (or `webautoma.exe` on Windows) in the current directory (`webautoma`).

3.  **(Optional) Move to PATH:** For easier access, you can move the `webautoma` executable to a directory listed in your system's PATH, like `/usr/local/bin` or similar.

**Quick Start Guide**

Let's run a simple test that opens Google and performs a search.

1.  **Download/Locate ChromeDriver:** Ensure you have the correct ChromeDriver executable for your Chrome version. Note its full path if it's not in your system's PATH.

2.  **Start ChromeDriver (Manual Method):**
    Open a separate terminal window and run ChromeDriver. By default, it listens on port 9515.
    ```bash
    /path/to/your/chromedriver --port=9515
    ```
    Leave this terminal window open while you run `webautoma`.

3.  **Create a Simple `.side` File:**
    Create a text file named `google_test.side` and paste the following JSON content into it:

    ```json
    {
      "id": "google-test-uuid",
      "version": "3.0",
      "name": "Simple Google Search",
      "url": "[https://www.google.com](https://www.google.com)",
      "tests": [{
        "id": "test-uuid",
        "name": "Search for WebAutoma",
        "commands": [{
          "id": "cmd-1",
          "command": "open",
          "target": "/",
          "value": ""
        }, {
          "id": "cmd-2",
          "command": "type",
          "target": "name=q",
          "value": "WebAutoma",
          "until": "1"
        }, {
          "id": "cmd-3",
          "command": "actionsSendKeys",
          "target": "",
          "value": "Enter"
        }, {
          "id": "cmd-4",
          "command": "exists",
          "target": "id=search",
          "value": "",
          "until": "1"
        }]
      }],
      "suites": [{
        "id": "suite-uuid",
        "name": "Default Suite",
        "tests": ["test-uuid"]
      }]
    }
    ```
    This test opens Google, types "WebAutoma" into the search bar (identified by `name=q`), presses Enter (using the custom `actionsSendKeys` command), and waits for the search results area (`id=search`) to exist.

4.  **Run `webautoma` (Connecting to Manual ChromeDriver):**
    Open another terminal window, navigate to the directory where you built `webautoma` and saved `google_test.side`, and run:
    ```bash
    ./webautoma -x google_test.side
    ```
    (Use `webautoma.exe` on Windows).

5.  **Observe:** You should see a Chrome browser window open, navigate to Google, perform the search, and then close. Check the console output and the generated `log.json` and `images.json` files in the same directory.

6.  **Run `webautoma` (Managing ChromeDriver - Alternative):**
    If you didn't start ChromeDriver manually in step 2, you can ask `webautoma` to start it by providing the path using the `-s` flag:
    ```bash
    ./webautoma -s /path/to/your/chromedriver -x google_test.side
    ```
    `webautoma` will attempt to start the ChromeDriver service, run the test, and then stop the service.

You have now successfully installed `webautoma` and executed your first test!