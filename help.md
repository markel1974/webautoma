# webautoma User Manual

## Introduction

(Placeholder: General description of webautoma, purposes, strengths, requirements)

## Installation and Quick Start

(Placeholder: Instructions on how to compile/install, WebDriver setup, basic startup example)

## Basic Usage (.side Execution)

(Placeholder: Explanation of main flags, output format, interpreting logs/images)

## .side File Format

It is based on the standard Selenium IDE v3 JSON format, with some additional commands and parameters specific to webautoma.

### General Structure

(Placeholder: Description of the main JSON structure of the .side file: version, name, url, tests, suites...)

### Command Structure

Each command within a test follows this basic JSON structure:

```json
{
"id": "unique-string",
"command": "commandName",
"target": "selector_or_value",
"value": "value_or_parameter",
// Specific webautoma fields may be present here
"until": "wait_condition", // Ex: "1"
"windowHandleName": "handle_name",
"windowTimeout": milliseconds,
"opensWindow": true/false,
"x": coordinate_x,
"y": coordinate_y,
"offsetX": offset_x,
"offsetY": offset_y
}
```

**Main Fields:**

- `id`: Unique identifier of the test step (standard Selenium IDE).
- `command`: The name of the webautoma command to execute (see below).
- `target`: The object of the action (e.g., CSS/XPath selector, URL, timer ID, variable name). Supports templating `{{.variable}}`.
- `value`: The value associated with the command (e.g., text to type, expected value, specific options). Supports templating `{{.variable}}`.

**Specific webautoma Fields (Optional):**

- `until`: Used by `assert`, `exists`, `until`. Indicates a wait condition (e.g., `1` to wait for the element to be ready - visible and enabled).
- `windowHandleName`, `windowTimeout`, `opensWindow`: Used for advanced window/tab management (see specific commands).
- `x`, `y`, `offsetX`, `offsetY`: Used for commands requiring coordinates or offsets (e.g., `scroll`, `scrollTo`).

## Supported Commands

Here is the list of commands recognized by webautoma:

### Category: Navigation

1. **open**
   - **Description:** Opens a URL in the browser or navigates to a path relative to the current base URL.
   - **Parameters:**
     - `target`: The full URL (e.g., https://google.com) or a relative path (e.g., /page). If relative, it is appended to the base url defined in the .side file or the last known base URL.
     - `value`: Not used.
   - **Example:**
     ```json
     {
       "id": "...",
       "command": "open",
       "target": "https://www.google.com",
       "value": ""
     }
     ```

### Category: Element Interaction

2. **click**
   - **Description:** Simulates a mouse click on the specified element. Waits for the element to be ready (visible and enabled) before clicking.
   - **Parameters:**
     - `target`: The selector of the element to click (e.g., id=myButton, css=.submit-btn).
     - `value`: Not used.
     - `until`: If set to `1` or `2`, waits for the element to be ready before clicking.
   - **Example:**
     ```json
     {
       "id": "...",
       "command": "click",
       "target": "id=loginButton",
       "value": "",
       "until": "1"
     }
     ```

3. **type**
   - **Description:** Enters text into an input or textarea field. Simulates character-by-character typing with a small delay (see humanWait).
   - **Parameters:**
     - `target`: The selector of the element to type into (e.g., id=username, name=password).
     - `value`: The text to insert. Supports variables (e.g., `{{.myUsername}}`).
     - `until`: If set to `1` or `2`, waits for the element to be ready before typing.
   - **Example:**
     ```json
     {
       "id": "...",
       "command": "type",
       "target": "id=searchField",
       "value": "Text to search"
     }
     ```

4. **select**
   - **Description:** Selects an option from a `<select>` (dropdown) element based on the option's visible text (label).
   - **Parameters:**
     - `target`: The selector of the `<select>` element.
     - `value`: The string `label=Option Text` that identifies the option to select.
     - `until`: If set to `1` or `2`, waits for the `<select>` element to be ready.
   - **Example:**
     ```json
     {
       "id": "...",
       "command": "select",
       "target": "id=countryDropdown",
       "value": "label=Italy"
     }
     ```

### Category: Timer Management (Custom webautoma)

5. **timerCreate**
   - **Description:** Creates and initializes a new timer, without starting it. Useful for measuring times spanning multiple actions.
   - **Parameters:**
     - `target`: The unique ID to assign to the timer (e.g., loginTime).
     - `value`: An optional description for the timer (reported in logs).

6. **timerStart**
   - **Description:** Starts (or restarts) a timer previously created with `timerCreate`. Records the start time.
   - **Parameters:**
     - `target`: The ID of the timer to start.
     - `value`: Not used.

7. **timerStop**
   - **Description:** Stops a previously started timer. Records the interval elapsed since the last `timerStart` or `timerStop`. If value is "finalize", it finalizes the timer and writes the complete event to the JSON log; otherwise, it only records the partial interval.
   - **Parameters:**
     - `target`: The ID of the timer to stop.
     - `value`: If set to `finalize` (case-insensitive), finalizes the timer. Otherwise, does nothing special besides stopping the current interval.

8. **timerFinalize** (Alternative to timerStop with value=finalize)
   - **Description:** Finalizes a timer, calculating the total elapsed time by summing all intervals recorded with `timerStop`. Writes the complete event to the JSON log. The timer can no longer be used after finalization.
   - **Parameters:**
     - `target`: The ID of the timer to finalize.
     - `value`: Not used.

### Category: Variable Stack Management (Custom webautoma)

9. **stackAdd**
   - **Description:** Finds an element, extracts some of its properties (text, tag, displayed/enabled status) and saves them in an internal map ("stack") associating them with the ID provided in the `id` field of the command. Useful for storing intermediate states or dynamic values.
   - **Parameters:**
     - `target`: The selector of the element to extract information from.
     - `value`: Not used directly.
     - `id`: (Standard command field) Important: This id is used as the key to store information in the stack.
     - `until`: If set to `1`, waits for the element to be ready.

10. **stackReset**
    - **Description:** Completely clears the internal variable stack.
    - **Parameters:** None used.

11. **stackPrint**
    - **Description:** Prints the current contents of the stack to the console (standard output) in indented JSON format. Useful for debugging.
    - **Parameters:** None used.

### Category: Flow Control & Utility

12. **jump** (Custom webautoma)
    - **Description:** Jumps execution to another command within the same test, identified by its `id`. Note: This is not a standard Selenium IDE command and makes the test flow harder to follow in the IDE itself.
    - **Parameters:**
       - `target`: The id of the command to jump to.
       - `value`: Not used.

13. **pause**
    - **Description:** Suspends execution for a specified number of milliseconds.
    - **Parameters:**
       - `target`: The number of milliseconds to suspend execution for.
       - `value`: Not used.

14. **humanWait** (Custom webautoma)
    - **Description:** Introduces a "human" pause, i.e., a random variable duration pause based on a configurable base value (`-humanWait` parameter or `executor.humanWaitBase` in code). If a value is provided in the command's `value`, it uses that as a fixed duration in milliseconds.
    - **Parameters:**
       - `target`: Not used.
       - `value`: Fixed duration of the pause in milliseconds (optional). If omitted, uses the random pause based on humanWaitBase.

### Category: Assertions and Verifications

15. **assert**
    - **Description:** Verifies that the visible text of an element exactly matches the provided value. Fails if the text is different.
    - **Parameters:**
       - `target`: The selector of the element whose text needs to be verified.
       - `value`: The exact text expected to be found in the element. Supports variables `{{.variable}}`.
       - `until`: (Optional) If set to `1` or `2`, waits for the element to be ready before verifying.

16. **exists**
    - **Description:** Verifies that a specified element exists in the DOM and is ready (visible and enabled) within the configured timeout. Fails if the element is not found or does not become ready.
    - **Parameters:**
       - `target`: The selector of the element to search for.
       - `value`: Not used.
       - `until`: (Optional) If set to `1` or `2`, actively waits for the element to exist and be ready for the duration of the timeout.

17. **until** (Custom webautoma)
    - **Description:** Waits until a specified element is *no longer* present or *no longer* ready (visible/enabled) on the page, or until the timeout expires. Useful for waiting for temporary elements (e.g., loading messages) to disappear. Fails if the element remains present and ready at the end of the timeout.
    - **Parameters:**
       - `target`: The selector of the element to monitor.
       - `value`: Not used.
       - `until`: (Optional, but **recommended to set to 1** for this command) If `1`, waits for the element to *not* be ready/present.

### Category: Window/Frame/Alert Management

18. **setWindowSize**
    - **Description:** Resizes the current window to the specified dimensions.
    - **Parameters:**
       - `target`: The desired size in "WidthxHeight" format (e.g., "1280x720").

19. **selectWindow** (Custom webautoma)
    - **Description:** Moves the driver focus to a specific window or tab. Can use the direct WebDriver handle, a name previously assigned with `storeWindowHandle` (using `${handleName}` syntax), or a numeric index.
    - **Parameters:**
       - `target`: The window identifier (`handle=...`, `${savedName}`).

20. **selectWindowMain** (Custom webautoma)
    - **Description:** Brings focus back to the main/initial window (the one opened at startup or set with `setWindowMain`).

21. **selectWindowTitle** (Custom webautoma)
    - **Description:** Moves focus to the first window/tab whose title matches (partially, exactly, or via regex) the specified `value`.
    - **Parameters:**
       - `target`: Title matching mode (`contains`, `exact`, `regexp`).
       - `value`: The title text (or regular expression) to search for.

22. **closeWindow** (Custom webautoma)
    - **Description:** Closes the *currently* focused window or tab. Cannot close the main/initial window.

23. **storeWindowHandle** (Custom webautoma)
    - **Description:** Saves the currently focused window handle, associating it with a symbolic name. This name can later be used in `selectWindow` or `close`.

24. **windowHandles** (Custom webautoma, Debug)
    - **Description:** Prints the list of handles for all currently open windows to the console.

25. **close** (Custom webautoma)
    - **Description:** Closes a specific window identified by its handle or a previously saved name.
    - **Parameters:**
       - `target`: The identifier of the window to close (e.g., `${popup}`).

26. **setWindowMain** (Custom webautoma)
    - **Description:** Sets the *currently* focused window handle as the new "main" or "root" window for `webautoma`.

### Category: Frame Management

27. **selectFrame**
    - **Description:** Moves focus to a frame (or iframe) within the current page.
    - **Parameters:**
       - `target`: Frame identifier (`index=N`, `relative=parent`, `relative=top`, or element Name/ID).

28. **selectParentFrame**
    - **Description:** Moves focus from the current frame to its direct parent frame. Equivalent to `selectFrame` with `target=relative=parent`.

### Category: Alert Management

29. **selectAlert** (Custom webautoma)
    - **Description:** Handles a JavaScript alert (`alert()`, `confirm()`, `prompt()`) appearing on the page. Reads the alert text (for logging) and then accepts or dismisses it.
    - **Parameters:**
       - `value`: `1` to accept, `0` or omitted to dismiss/cancel.

### Category: Advanced Interactions (Mouse)

30. **rightClick**
    - **Description:** Simulates a right mouse click on the specified element. Waits for the element to be ready.
    - **Parameters:**
       - `target`: The selector of the element to right-click on.
       - `until`: (Optional) If `1` or `2`, waits for the element to be ready.

31. **doubleClick**
    - **Description:** Simulates a double left mouse click on the specified element. Waits for the element to be ready.

32. **mouseOver**
    - **Description:** Moves the mouse cursor over the specified element, potentially triggering hover effects or tooltips. Waits for the element to be ready.

33. **mouseOut**
    - **Description:** Simulates the mouse leaving the area of the specified element. **Note:** In the current code (`doMouse`), this command doesn't seem to execute specific actions on the WebDriver.

34. **mouseDownAt**
    - **Description:** Simulates pressing (without releasing) the left mouse button on the specified element, potentially at relative coordinates. Starts a drag operation.
    - **Parameters:**
       - `target`: The selector of the element to press the mouse on.
       - `value`: (Optional) Relative coordinates to the top-left corner of the element, format "X,Y". If omitted, presses in the center.

35. **mouseMoveAt**
    - **Description:** Moves the mouse while the left button is held down (started with `mouseDownAt`). Used for dragging. **Requires** a previous `mouseDownAt`.
    - **Parameters:**
       - `value`: Relative coordinates to the top-left corner of the *original* element from `mouseDownAt`, format "X,Y". Indicates the position *to which* to move the mouse.

36. **mouseMultipleMoveAt** (Custom webautoma)
    - **Description:** Similar to `mouseMoveAt`, but executes a sequence of consecutive relative movements while the button is pressed. Useful for simulating dragging along a path. **Requires** a previous `mouseDownAt`.
    - **Parameters:**
       - `value`: Sequence of *incremental* relative coordinates, separated by `|`. Each "X,Y" pair is relative to the *previous position*. Format: "dX1,dY1|dX2,dY2|...".

37. **mouseUpAt**
    - **Description:** Simulates releasing the left mouse button, completing a drag and drop operation. **Requires** a previous `mouseDownAt`.
    - **Parameters:**
       - `value`: (Optional) Relative coordinates to the top-left corner of the *original* element from `mouseDownAt`, indicating the *final* release position.

### Category: Advanced Interactions (Keyboard - Custom webautoma)

38. **actionsSendKeys** (Custom webautoma)
    - **Description:** Sends a special key press (non-alphanumeric) or a simple sequence to the active element on the page. Useful for simulating Enter, Tab, Arrow keys, etc.
    - **Parameters:**
       - `value`: The name of the special key (e.g., `Enter`, `Tab`, `ArrowDown`, `Control`, `Alt`, `Shift`, `F5`, etc.) or a sequence of simple characters.

39. **keys** (Custom webautoma)
    - **Description:** Allows defining complex sequences of keyboard interactions, including press (`press`), release (`release`), and pauses (`pause`). Useful for simulating key combinations (e.g., Ctrl+C).
    - **Parameters:**
       - `target`: (Optional) Selector of the element to send events to. If omitted, sends to the active element.
       - `value`: Formatted string describing the sequence, comma-separated. Each item is `action:value`. (e.g., `press:Control,press:c,release:c,release:Control`).

### Category: File & Download (Custom webautoma)

40. **download**
    - **Description:** Downloads a file directly from a specified URL and saves it to the indicated local path. Uses current browser session cookies to handle any server-required authentication.
    - **Parameters:**
       - `target`: The full URL to download the file from. Supports `{{.variable}}` variables.
       - `value`: The full path, including filename, where to save the downloaded file locally. Supports `{{.variable}}` variables.

41. **clickDownload**
    - **Description:** Finds an element on the page (usually an `<a>` link), extracts its URL from the `href` attribute, then downloads the file from that URL to the specified local path. Useful when the download URL is dynamic.
    - **Parameters:**
       - `target`: The selector of the element containing the `href` attribute with the file URL.
       - `value`: The full local save path.

### Category: Scrolling (Custom webautoma)

42. **scroll**
    - **Description:** Scrolls the page (viewport) or from a specific element, by a given amount (delta X, delta Y). Can simulate smooth scrolling by splitting it into steps with intermediate pauses. Uses the WebDriver Actions API.
    - **Parameters:**
       - `target`: (Optional) The selector of the element to calculate the starting point of the scroll from.
       - `value`: Specifies the movement (delta) and optionally steps and interval. Format: "DeltaX,DeltaY[,NumberSteps,IntervalMs]".
       - `x`, `y`: (Optional) Absolute X, Y coordinates in the viewport to start scrolling from.
       - `offsetX`, `offsetY`: (Optional) Offset in pixels to add to the starting coordinates.

43. **scrollTo**
    - **Description:** Scrolls the page so that a specific point *within* a target element becomes visible in the viewport. Can simulate smooth scrolling. Uses the WebDriver Actions API.
    - **Parameters:**
       - `target`: (Optional) The selector of the target element to scroll towards.
       - `value`: Specifies coordinates *relative to the top-left corner of the target element* and optionally steps and interval. Format: "TargetX,TargetY[,NumberSteps,IntervalMs]".

### Category: Advanced Custom (webautoma)

44. **otp**
    - **Description:** Retrieves a One-Time Password (OTP) from an email account. Connects to the specified IMAP server, searches for the most recent email matching criteria, extracts the OTP via a regular expression, and stores it internally. Accessible as `{{.otp}}`.
    - **Parameters:**
       - `target`: IMAP connection string. Format: `protocol[authMode]://user:password@server:port`
       - `value`: Configuration string for search and extraction, parameters separated by `|||`. Format: `"SubjectRegExp|||BodyRegExpWithOTPCaptureGroup[|||CheckIntervalSec[|||MaxEmailAgeMin]]"`

### Category: Debug and Meta-Commands (Custom webautoma)

45. **disableError**
    - **Description:** Temporarily disables execution interruption in case of errors in subsequent commands. Error will still be logged. Use `enableError` to restore normal behavior.

46. **enableError**
    - **Description:** Re-enables execution interruption in case of an error, canceling the effect of a previous `disableError`. This is the default behavior.

47. **disableDebug**
    - **Description:** Disables detailed log output (debug level) generated by the `webautoma` adapter during execution.

48. **enableDebug**
    - **Description:** Re-enables detailed log output (debug level).

49. **status**
    - **Description:** Requests and prints WebDriver server status information (version, OS, availability) to the console.

50. **execId**
    - **Description:** Sets a global identifier for the entire current execution. This ID will be included in generated JSON logs.
    - **Parameters:**
       - `target`: The string to use as the execution ID.

51. **probe**
    - **Description:** Sets an identifier for the "probe" or specific test instance currently running. Included in JSON logs.
    - **Parameters:**
       - `target`: The string to use as the probe/instance ID.

52. **setTimeout**
    - **Description:** Sets the maximum time (implicit timeout) in milliseconds the WebDriver will wait when searching for an element before returning an error.
    - **Parameters:**
       - `target`: The wait time in milliseconds.

53. **activeElement** (Debug)
    - **Description:** Identifies the currently active element (the one with focus) and prints its details to the console.

54. **pageSource** (Debug)
    - **Description:** Retrieves the entire HTML source of the currently displayed page and prints it to the console.

55. **noop**
    - **Description:** "No Operation" command. Performs no action. Can be used as a placeholder or to insert comments in the `.side` file.

## Target/Selectors Syntax

(Placeholder: Detailed description of supported target formats: id=, css=, xpath=, linkText=, name=, class=)

## Value/Variables Syntax

(Placeholder: Explanation of using value and `{{.variableName}}` templating to insert dynamic values)

## Interactive Console (SAM)

(Placeholder: Description of the SAM console, how to start it, available commands, usage for debugging)

## Server Mode

(Placeholder: Description of server mode, how to start it, APIs, purpose)

## Advanced Features

(Placeholder: Details on email OTP, advanced configuration, capabilities management, network event capture, etc.)

## Troubleshooting

(Placeholder: Common errors, interpreting logs, known issues)
