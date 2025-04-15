Basic Usage (`.side` Execution)

The primary way to use `webautoma` is to execute a Selenium IDE `.side` file from your command line. This section covers the fundamental command structure and the essential flags required for running a test.

**Command Structure**

The basic syntax for running a test is:

webautoma [flags] -x <your_test.side>

You replace `[flags]` with options to configure the execution and `<your_test.side>` with the path to your test definition file.

**Essential Command-Line Flags**

While `webautoma` offers various flags (see Command-Line Arguments section or use `webautoma -h`), these are the most commonly used for basic test execution:

* `-x <file.side>`: (Required) Specifies the path to the `.side` JSON file containing the test(s) you want to run.

* `-s <path>`: Specifies the full path to the WebDriver executable (e.g., `/usr/local/bin/chromedriver`).
    * If provided, `webautoma` will attempt to start this WebDriver service before the test and stop it afterwards.
    * This is often necessary if the WebDriver executable is not in your system's PATH.
    * If you have already started the WebDriver service manually (e.g., running `chromedriver` in another terminal), you typically *omit* the `-s` flag and use `-b`/`-p` if needed.

* `-b <url>`: The base URL of the running WebDriver service (default: `http://127.0.0.1`). Only needed if the service is not running on the default URL and you are *not* using the `-s` flag.

* `-p <port>`: The port number of the running WebDriver service (default: `9515`). Only needed if the service is not running on the default port and you are *not* using the `-s` flag.

* `-u <prefix>`: The URL prefix for the WebDriver service, if required (e.g., `/wd/hub` for Selenium Server Grid).

**Output Files**

By default, `webautoma` generates two main output files in the directory where it's executed:

* `-l <logfile.json>`: Specifies the path for the detailed execution log (default: `log.json`). This JSON file contains step-by-step results, including command details, timings, success/failure status, error messages, network summaries, and screenshot references.
* `-i <imagefile.json>`: Specifies the path for the image metadata log (default: `images.json`). This JSON file contains information about each screenshot taken during the run (e.g., associated step ID, timestamp).
* `-a`: If this flag is present, `webautoma` will also save the actual screenshot images as PNG files (named with timestamps) in the execution directory, in addition to logging their metadata in the image file.

**Passing Variables**

* `-z <vars>`: Allows you to pass dynamic data into your `.side` test. This is useful for parameterizing tests (e.g., different usernames, environments).
    * Format 1 (Inline): `key1=value1;key2=value2`
    * Format 2 (File): `@<filename.ndjson>` (Loads variables from a file containing newline-delimited JSON objects).
    * Usage in `.side`: Inside the `target` and `value` fields of your `.side` commands, use Go template syntax: `{{.key1}}`, `{{.key2}}`. (More details in the `.side` File Format section).

**Passing Arguments to WebDriver**

* `-d <args>`: Allows passing specific command-line arguments directly to the WebDriver service executable when started via the `-s` flag. Separate arguments with a semicolon. Useful for enabling features like headless mode.
    * Example: `-d "--headless;--window-size=1920,1080"`

**Example Command**

Here's a more complete example command line:

```bash
./webautoma \
    -s /usr/local/bin/chromedriver \
    -x /path/to/my_complex_test.side \
    -l test_run_log.json \
    -i test_run_images.json \
    -a \
    -z "baseURL=[https://staging.example.com](https://staging.example.com);user=testuser" \
    -d "--headless;--window-size=1920,1080"