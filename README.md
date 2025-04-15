# webautoma

`webautoma` is a powerful Go-based WebDriver automation tool designed to execute Selenium IDE (`.side`) test files. It acts as a runner that interacts with browser drivers (like ChromeDriver) to perform web actions defined in your `.side` projects, while offering extended capabilities beyond the standard Selenium IDE feature set.

Use `webautoma` to automate browser tasks, run regression tests, monitor web application performance, and leverage unique features like integrated email OTP retrieval and an interactive debugging console.

## Features

* **Selenium IDE (`.side`) Execution:** Faithfully runs tests exported from Selenium IDE (v3 JSON format).
* **WebDriver Integration:** Connects to and controls standard WebDriver services (e.g., ChromeDriver, GeckoDriver via Selenium Server).
* **Service Management:** Can optionally manage the lifecycle of the WebDriver service executable (e.g., start ChromeDriver via `-s` flag).
* **Extended Command Set:** Implements many standard Selenium IDE commands and adds numerous custom commands for:
    * **Performance Timing:** Create, start, stop, and finalize timers (`timerCreate`, `timerStart`, `timerStop`, `timerFinalize`) to measure specific parts of your tests.
    * **Variable Management:** Store dynamic values from page elements (`stackAdd`) and use them later via templating. Reset (`stackReset`) and inspect (`stackPrint`) the variable stack.
    * **Advanced Interactions:** Perform right-clicks, double-clicks, mouse-overs, complex drag-and-drop sequences (`mouseDownAt`, `mouseMoveAt`, `mouseUpAt`, `mouseMultipleMoveAt`), and intricate keyboard actions (`keys` command for press/release/pause sequences, `actionsSendKeys` for special keys like Enter/Tab/Arrows).
    * **File Handling:** Download files directly from URLs or by clicking links, using current browser cookies (`download`, `clickDownload`).
    * **Scrolling Control:** Scroll the page by specific amounts or scroll elements into view with fine-grained control (`scroll`, `scrollTo`).
    * **Email OTP Retrieval:** Connect to IMAP email accounts, find emails matching subject/body patterns, and extract One-Time Passwords for use in login flows (`otp` command, result available as `{{.otp}}`).
    * **Conditional Logic & Flow:** Verify element existence/readiness (`exists`), assert text content (`assert`), wait until an element disappears (`until`), and jump between test steps (`jump`).
    * **Window/Frame/Alert Management:** Switch between windows/tabs by handle, title, or index; manage frames; accept/dismiss JavaScript alerts (`selectWindow*`, `selectFrame`, `selectAlert`, `storeWindowHandle`, etc.).
    * **Debug & Meta Controls:** Enable/disable error halting or debug logs (`enable/disableError`, `enable/disableDebug`), set timeouts (`setTimeout`), inspect state (`status`, `activeElement`, `pageSource`), set run identifiers (`execId`, `probe`).
* **Detailed JSON Logging:** Creates structured logs (`log.json` by default) recording each command, execution time, status (success/error), network activity summary, and optional screenshot references.
* **Screenshot & Image Metadata:** Captures screenshots on error (or optionally on demand) and saves image metadata (ID, timestamp, associated log event) to a separate JSON file (`images.json` by default). Can optionally dump PNG images to disk (`-a` flag).
* **Network Monitoring:** Leverages WebDriver's performance logging capabilities to capture network requests/responses. Can extract specific headers (`-c` flag) for use in logs or analysis.
* **Variable Templating:** Inject dynamic data into `.side` file commands (in `target` and `value` fields) using Go template syntax (`{{.variable_name}}`). Variables are provided via the `-z` command-line flag (inline or from a file).
* **Interactive Debug Console (SAM):** The unique "Step-Aside-Mode" (`-e` flag) provides a terminal UI to step through `.side` commands, inspect the current command, view browser state (e.g., page source, network logs), edit command parameters on-the-fly, and jump between steps. Ideal for debugging complex tests.
* **Server Mode:** Can operate as an HTTP server (`-y` flag), allowing tests to be triggered remotely or integrated into larger systems. *(Note: API and specific functionality TBD based on code implementation).*
* **Cross-Platform:** Written in Go, enabling cross-compilation for various operating systems.

## Requirements

* **Go:** A recent version of the Go compiler (e.g., 1.18 or later, check `go.mod`). Ensure Go is installed and configured correctly in your environment.
* **WebDriver Service:** A compatible WebDriver executable for the browser you intend to automate.
    * **ChromeDriver:** Download from [https://chromedriver.chromium.org/downloads](https://chromedriver.chromium.org/downloads)
    * *(Optional) GeckoDriver for Firefox, etc.*
      The WebDriver executable should either be in your system's PATH or its location must be specified using the `-s` flag when running `webautoma`.

## Installation

1.  **Clone the Repository:**
    ```bash
    # git clone <repository_url>
    # cd webautoma
    ```
2.  **Build the Executable:**
    ```bash
    go build -o webautoma main.go
    ```
    This will create the `webautoma` executable in the current directory. You might want to move it to a location in your system's PATH for easier access.

## Getting Started

1.  **Download ChromeDriver:** Get the appropriate ChromeDriver version for your installed Chrome browser.
2.  **Start ChromeDriver (Option 1 - Manual):**
    Run ChromeDriver in a separate terminal window. It will typically listen on `http://127.0.0.1:9515`.
    ```bash
    /path/to/chromedriver --port=9515
    ```
3.  **Prepare a `.side` File:**
    Use Selenium IDE in your browser to record a simple test (e.g., opening a site, clicking a link). Save the test as `test.side`.
4.  **Run `webautoma`:**
    * If you started ChromeDriver manually:
        ```bash
        ./webautoma -x test.side
        ```
    * If you want `webautoma` to start ChromeDriver (provide the path):
        ```bash
        ./webautoma -s /path/to/chromedriver -x test.side
        ```
5.  **Review Output:** Check the console output and the generated `log.json` and `images.json` files for execution details.

## Command-Line Arguments

`webautoma [flags]`

**Core Execution:**

* `-x <file.side>`: **(Required)** Path to the Selenium IDE `.side` file to execute.
* `-z <vars>`: Define variables for the `.side` file.
    * Inline format: `key1=value1;key2=value2`
    * File format: `@<filename.ndjson>` (Loads variables from a file containing newline-delimited JSON objects, e.g., `{"key1":"value1"}\n{"key2":"value2"}`). Use `{{.key}}` within the `.side` file's `target` and `value` fields.

**WebDriver Configuration:**

* `-s <path>`: Path to the WebDriver executable (e.g., ChromeDriver). If provided, `webautoma` attempts to start this service.
* `-b <url>`: WebDriver base URL (default: `http://127.0.0.1`).
* `-p <port>`: WebDriver port (default: `9515`).
* `-u <prefix>`: WebDriver URL prefix if needed (e.g., `/wd/hub` for Selenium Server).
* `-d <args>`: Pass arguments directly to the WebDriver service, separated by semicolons (e.g., `-d "--headless;--disable-gpu"`).

**Output & Logging:**

* `-l <logfile.json>`: Path for the detailed JSON results log (default: `log.json`). Contains timing, status, network info for each step.
* `-i <imagefile.json>`: Path for the JSON image metadata log (default: `images.json`). Contains references to screenshots.
* `-a`: Enable image dumping. Saves screenshots (on error or explicit command) as timestamped PNG files in the execution directory.

**Debugging & Advanced Modes:**

* `-e`: Enable interactive SAM (Step-Aside-Mode) console for debugging execution step-by-step.
* `-c <headers>`: Comma-separated list of *request* or *response* headers to capture from WebDriver performance logs and include in the output logs (e.g., `-c "Authorization,X-Request-ID,Content-Type"`). Requires browser/driver support for performance logging.
* `-y <listen_addr>`: Start `webautoma` in server mode, listening on the specified address (e.g., `:8080`). Disables direct `.side` execution.
* `-w`: WebDriver service mode. Starts the WebDriver specified by `-s` and exits. Useful for preparing the environment.

**Other:**

* `-h`: Show the help message.
* `-v`: Show `webautoma` version information.
* `-t <package:name:dir>`: (For development) Build Go asset file from a directory.

## `.side` File Execution Details

`webautoma` processes `.side` files based on the Selenium IDE v3 JSON format.

* **Standard Commands:** Many common commands like `open`, `click`, `type`, `select`, `pause`, `assertText`, `verifyText`, `executeScript`, etc., are supported.
* **Custom Commands:** Offers unique commands like `timer*`, `stack*`, `otp`, `download`, `clickDownload`, `keys`, `actionsSendKeys`, `scroll*`, `jump`, `exists`, `until`, `*Window*`, `*Frame*`, `selectAlert`, `enable/disableError`, `enable/disableDebug`, `status`, `probe`, `execId`, `humanWait`, `noop`. Refer to `manual.md` for details on each custom command and its parameters.
* **Selectors (`target`):** Supports standard WebDriver locator strategies prefixed in the `target` field:
    * `id=<element_id>`
    * `name=<element_name>`
    * `linkText=<visible_link_text>`
    * `partialLinkText=<partial_visible_text>`
    * `css=<css_selector>`
    * `xpath=<xpath_expression>`
    * `class=<class_name>` (Convenience, often equivalent to `css=.class_name`)
* **Variables (`value`, `target`):** Use `{{.variableName}}` syntax within the `target` and `value` fields of your commands to insert values provided via the `-z` flag.

## SAM Console (Interactive Debugging)

Start `webautoma` with the `-e` flag to enter the Step-Aside-Mode console.

* **Purpose:** Allows interactive execution and debugging of `.side` files.
* **Key Commands:**
    * `step [N]`: Execute the next N command(s) (default 1).
    * `next`: Move the cursor to the next command without executing.
    * `prev`: Move the cursor to the previous command.
    * `curr`: Display the current command JSON.
    * `list`: Display the entire list of commands in the current test.
    * `jump <index>`: Set the execution pointer to the command at the specified index.
    * `edit <field> <new_value>`: Modify a field (`command`, `target`, `value`, etc.) of the *current* command before stepping.
    * `redo`: Execute the *current* command again without advancing the cursor.
    * `print html [file]`: Print current page source (optionally save to file).
    * `print net [file]`: Print captured network headers (optionally save to file).
    * `windows`: List currently open window handles and titles.
    * `run`: Resume normal execution until the end or next breakpoint (TBD).
    * `stop`: Halt execution (TBD).
    * `exit`: Quit the SAM console and `webautoma`.
    * Use `help` within the SAM console for a command list.

## Logging and Output

`webautoma` generates two primary JSON log files:

1.  **Results Log (`-l`, default `log.json`):** Contains a sequence of JSON objects, one for each executed command or timer finalization. Includes timestamps, command details, execution time, success/failure status, error messages, network error counts, and captured headers/variables.
2.  **Image Log (`-i`, default `images.json`):** Contains a sequence of JSON objects, one for each screenshot taken. Includes a unique ID (referenced in the results log), timestamp, and potentially the base64-encoded image data or file path if `-a` is used.

## Architecture Overview (Key Packages)

* `main.go`: Entry point, command-line flag parsing.
* `executor/`: Core logic for parsing `.side` files, executing commands, managing state (timers, stack), handling events, logging, and the SAM console.
    * `executor/adapter.go`: Wraps the WebDriver client, adding retry logic, logging, screenshots, network handling, human waits.
* `wd/`: WebDriver client implementation (seems to be a custom wrapper or implementation). Handles communication protocol.
* `shell/`: Implements the terminal UI framework used by the SAM console.
* `service/`: Code for managing external processes like ChromeDriver and Xvfb.
* `email/`: IMAP client logic for the `otp` command.
* `server/`: HTTP server implementation for the `-y` mode.

## Troubleshooting

* **WebDriver Connection Issues:** Ensure the WebDriver service (e.g., ChromeDriver) is running and accessible at the specified URL/port (`-b`, `-p`, `-u`). Check for version compatibility between the browser, driver, and `webautoma` (if relevant). Verify firewall settings.
* **Element Not Found:** Double-check your selectors (`target` values). Use browser developer tools to verify IDs, CSS paths, or XPaths. Increase implicit wait time (`setTimeout` command) if elements load slowly. Use `exists` with `until=1` to wait for elements.
* **Command Failures:** Examine the `log.json` file for detailed error messages associated with the failing command ID. Enable image dumping (`-a`) to see the state of the page when an error occurs. Use the SAM console (`-e`) to step through the test and isolate the issue.

## Full Documentation

For the most detailed information, including explanations for *all* custom commands and advanced configurations, please refer to the `manual.md` file (currently available in Italian).

*(Placeholder: Add sections for Contributing, Development Setup, etc. if needed)*

## License

This project is licensed under the Apache License 2.0. See file headers for details.