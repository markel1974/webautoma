package help

import (
	"fmt"
	"sort"
	"strings"
)

// Provider provides direct access to predefined help strings for commands.
type Provider struct {
	// Map containing help text for each command.
	// Key: command name (e.g. "open", "selectWindow")
	// Value: multi-line string with description, parameters and example.
	commandHelp map[string]string
}

// NewHelp creates and initializes a new Provider.
// Help strings are defined directly here.
func NewHelp() *Provider {
	helpMap := make(map[string]string)

	// Populate the map with descriptions
	// (Using backticks ` for multi-line strings)

	// --- Category: Navigation ---
	helpMap["open"] = `
        - Description: Opens a URL in the browser or navigates to a path relative to the current base URL.
        - Parameters:
          - target: The full URL (e.g. https://google.com) or a relative path (e.g. /page). If relative, it is appended to the base URL defined in the .side file or to the last known base URL.
          - value: Not used.
        - Example:
          {
            "id": "...",
            "command": "open",
            "target": "https://www.google.com",
            "value": ""
          }
    `

	// --- Category: Element Interaction ---
	helpMap["click"] = `
        - Description: Simulates a mouse click on the specified element. Waits for the element to be ready (visible and enabled) before clicking.
        - Parameters:
          - target: The selector of the element to click (e.g. id=myButton, css=.submit-btn).
          - value: Not used.
          - until: If set to 1 or 2, waits for the element to be ready before clicking.
        - Example:
          {
            "id": "...",
            "command": "click",
            "target": "id=loginButton",
            "value": "",
            "until": "1"
          }
    `
	helpMap["type"] = `
       - Description: Enters text into an input or textarea field. Simulates character-by-character typing with a small delay (see humanWait).
       - Parameters:
         - target: The selector of the element to type into (e.g. id=username, name=password).
         - value: The text to insert. Supports variables (e.g. {{.myUsername}}).
         - until: If set to 1 or 2, waits for the element to be ready before typing.
       - Example:
         {
           "id": "...",
           "command": "type",
           "target": "id=searchField",
           "value": "Text to search"
         }
         {
           "id": "...",
           "command": "type",
           "target": "name=password",
           "value": "{{.userPassword}}"
         }
    `
	helpMap["select"] = `
       - Description: Selects an option from a <select> (dropdown) element based on the option's visible text (label).
       - Parameters:
         - target: The selector of the <select> element.
         - value: The string label=Option Text that identifies the option to select.
         - until: If set to 1 or 2, waits for the <select> element to be ready.
       - Example:
         {
           "id": "...",
           "command": "select",
           "target": "id=countryDropdown",
           "value": "label=Italy"
         }
    `

	// --- Category: Timer Management (Custom webautoma) ---
	helpMap["timerCreate"] = `
       - Description: Creates and initializes a new timer, without starting it. Useful for measuring times spanning multiple actions.
       - Parameters:
         - target: The unique ID to assign to the timer (e.g. loginTime).
         - value: An optional description for the timer (reported in logs).
       - Example:
         {
           "id": "...",
           "command": "timerCreate",
           "target": "pageLoadTimer",
           "value": "Initial page load time"
         }
    `
	helpMap["timerStart"] = `
       - Description: Starts (or restarts) a timer previously created with timerCreate. Records the start time.
       - Parameters:
         - target: The ID of the timer to start.
         - value: Not used.
       - Example:
         {
           "id": "...",
           "command": "timerStart",
           "target": "pageLoadTimer",
           "value": ""
         }
    `
	helpMap["timerStop"] = `
       - Description: Stops a previously started timer. Records the interval elapsed since the last timerStart or timerStop. If value is "finalize", it finalizes the timer and writes the complete event to the JSON log; otherwise, it only records the partial interval.
       - Parameters:
         - target: The ID of the timer to stop.
         - value: If set to finalize (case-insensitive), finalizes the timer. Otherwise, does nothing special besides stopping the current interval.
       - Example (Partial stop):
         {
           "id": "...",
           "command": "timerStop",
           "target": "userActionTimer",
           "value": ""
         }
       - Example (Stop and Finalize):
         {
           "id": "...",
           "command": "timerStop",
           "target": "totalTestTimer",
           "value": "finalize"
         }
    `
	helpMap["timerFinalize"] = `
       - Description: Finalizes a timer, calculating the total elapsed time by summing all intervals recorded with timerStop. Writes the complete event to the JSON log. The timer can no longer be used after finalization.
       - Parameters:
         - target: The ID of the timer to finalize.
         - value: Not used.
       - Example:
         {
           "id": "...",
           "command": "timerFinalize",
           "target": "loginProcessTimer",
           "value": ""
         }
    `

	// --- Category: Variable Stack Management (Custom webautoma) ---
	helpMap["stackAdd"] = `
       - Description: Finds an element, extracts some of its properties (text, tag, displayed/enabled status) and saves them in an internal map ("stack") associating them with the ID provided in the command's id field. Useful for storing intermediate states or dynamic values.
       - Parameters:
         - target: The selector of the element to extract information from.
         - value: Not used directly.
         - id: (Standard command field) Important: This id is used as the key to store information in the stack.
         - until: If set to 1, waits for the element to be ready.
       - Example:
         {
           "id": "userInfo", // This ID will be the key in the stack
           "command": "stackAdd",
           "target": "id=userDetails",
           "value": "",
           "until": "1"
         }
    `
	helpMap["stackReset"] = `
        - Description: Completely clears the internal variable stack.
        - Parameters: None used.
        - Example:
          {
            "id": "...",
            "command": "stackReset",
            "target": "",
            "value": ""
          }
    `
	helpMap["stackPrint"] = `
        - Description: Prints the current contents of the stack to the console (standard output) in indented JSON format. Useful for debugging.
        - Parameters: None used.
        - Example:
          {
            "id": "...",
            "command": "stackPrint",
            "target": "",
            "value": ""
          }
    `

	// --- Category: Flow Control & Utility ---
	helpMap["jump"] = `
        - Description: Jumps execution to another command within the same test, identified by its id. Note: This is not a standard Selenium IDE command and makes the test flow harder to follow in the IDE itself.
        - Parameters:
          - target: The id of the command to jump to.
          - value: Not used.
        - Example:
          {
            "id": "...",
            "command": "jump",
            "target": "startLoop", // Jumps to command with id "startLoop"
            "value": ""
          }
    `
	helpMap["pause"] = `
        - Description: Suspends execution for a specified number of milliseconds.
        - Parameters:
          - target: The number of milliseconds to suspend execution for.
          - value: Not used.
        - Example:
          {
            "id": "...",
            "command": "pause",
            "target": "5000", // Pause for 5 seconds
            "value": ""
          }
    `
	helpMap["humanWait"] = `
        - Description: Introduces a "human" pause, i.e., a random variable duration pause based on a configurable base value (-humanWait parameter or executor.humanWaitBase in code). If a value is provided in the command's value, it uses that as a fixed duration in milliseconds.
        - Parameters:
          - target: Not used.
          - value: Fixed duration of the pause in milliseconds (optional). If omitted, uses the random pause based on humanWaitBase.
        - Example (Random pause):
          {
            "id": "...",
            "command": "humanWait",
            "target": "",
            "value": ""
          }
        - Example (Fixed pause):
          {
            "id": "...",
            "command": "humanWait",
            "target": "",
            "value": "1500" // Fixed pause of 1.5 seconds
          }
    `

	// --- Category: Assertions and Verifications ---
	helpMap["assert"] = `
        - Description: Verifies that the visible text of an element exactly matches the provided value. Fails if the text is different.
        - Parameters:
          - target: The selector of the element whose text needs to be verified (e.g. id=errorMessage, css=.result).
          - value: The exact text expected to be found in the element. Supports variables {{.variable}}.
          - until: (Optional) If set to 1 or 2, waits for the element to be ready (visible and enabled) before verifying.
        - Example:
          {
            "id": "...",
            "command": "assert",
            "target": "id=statusMessage",
            "value": "Operation completed.",
            "until": "1"
          }
    `
	helpMap["exists"] = `
        - Description: Verifies that a specified element exists in the DOM and is ready (visible and enabled) within the configured timeout. Fails if the element is not found or does not become ready.
        - Parameters:
          - target: The selector of the element to search for (e.g. id=popupConfirm, css=button.primary).
          - value: Not used.
          - until: (Optional) If set to 1 or 2, actively waits for the element to exist and be ready for the duration of the timeout. If omitted (or 0), only checks if the element is present and ready in the current page state (less common for robust checks).
        - Example:
          {
            "id": "...",
            "command": "exists",
            "target": "css=.loading-spinner",
            "value": "",
            "until": "1" // Waits for spinner to be visible/enabled
          }
    `
	helpMap["until"] = `
        - Description: Waits until a specified element is *no longer* present or *no longer* ready (visible/enabled) on the page, or until the timeout expires. Useful for waiting for temporary elements (e.g. loading messages) to disappear. Fails if the element remains present and ready at the end of the timeout.
        - Parameters:
          - target: The selector of the element to monitor for its disappearance or inactivity.
          - value: Not used.
          - until: (Optional, but **recommended to set to 1** for this command) If 1, waits for the element to *not* be ready/present. If 0 or 2, waits for it to be ready, but the command will fail if the element *is* actually ready (less intuitive usage).
        - Example:
          {
            "id": "...",
            "command": "until",
            "target": "css=.loading-indicator",
            "value": "",
            "until": "1" // Waits for loading indicator to disappear or become not ready
          }
    `

	// --- Category: Window/Frame/Alert Management ---
	helpMap["setWindowSize"] = `
        - Description: Resizes the current window to the specified dimensions.
        - Parameters:
          - target: The desired size in "WidthxHeight" format (e.g. "1280x720").
          - value: Not used.
        - Example:
          {
            "id": "...",
            "command": "setWindowSize",
            "target": "1920x1080",
            "value": ""
          }
    `
	helpMap["selectWindow"] = `
        - Description: Moves the driver focus to a specific window or tab. Can use the direct WebDriver handle, a name previously assigned with storeWindowHandle (using the ${handleName} syntax), or a numeric index (less common/reliable).
        - Parameters:
          - target: The window identifier. Possible formats:
              - handle=WEBDRIVER_HANDLE_NAME (rarely used manually)
              - ${savedHandleName} (name assigned with storeWindowHandle)
              - Potentially a numeric index (to verify in WebDriver code)
          - value: Not used.
        - Example (using a saved handle):
          {
            "id": "...",
            "command": "selectWindow",
            "target": "${mainWindow}",
            "value": ""
          }
    `
	helpMap["selectWindowMain"] = `
        - Description: Brings focus back to the main/initial window (the one opened at startup or set with setWindowMain).
        - Parameters: None used.
        - Example:
          {
            "id": "...",
            "command": "selectWindowMain",
            "target": "",
            "value": ""
          }
    `
	helpMap["selectWindowTitle"] = `
        - Description: Moves focus to the first window/tab whose title matches (partially, exactly, or via regex) the specified value.
        - Parameters:
          - target: Title matching mode (Optional, default: contains):
              - contains: The window title contains value (case-insensitive).
              - exact: The window title is exactly value.
              - regexp: The window title matches the regular expression in value.
          - value: The title text (or regular expression) to search for.
        - Example (Contains):
          {
            "id": "...",
            "command": "selectWindowTitle",
            "target": "contains", // or omitted
            "value": "Results Page"
          }
        - Example (Regexp):
          {
            "id": "...",
            "command": "selectWindowTitle",
            "target": "regexp",
            "value": "^Cart \\(\\d+\\)$" // E.g.: Title "Cart (3)"
          }
    `
	helpMap["closeWindow"] = `
        - Description: Closes the window or tab *currently* in focus. Cannot close the main/initial window. After closing, focus is *not* automatically moved; use selectWindowMain or another selectWindow* command if necessary.
        - Parameters: None used.
        - Example:
          {
            "id": "...",
            "command": "closeWindow",
            "target": "",
            "value": ""
          }
    `
	helpMap["storeWindowHandle"] = `
        - Description: Saves the handle of the currently focused window, associating it with a symbolic name. This name can be used subsequently in selectWindow or close with the ${handleName} syntax.
        - Parameters:
          - target: The name to assign to the current window handle (e.g. loginWindow, detailPopup).
          - value: Not used.
        - Specific Fields: Can use windowHandleName (redundant?), windowTimeout, opensWindow to handle waiting for a new window to open before saving its handle (if opensWindow is true).
        - Example (Saving current handle):
          {
            "id": "...",
            "command": "storeWindowHandle",
            "target": "mainWindow", // Saves current handle as "mainWindow"
            "value": ""
          }
        - Example (Waiting for and saving popup):
          {
            "id": "...",
            "command": "storeWindowHandle",
            "target": "popupWindow",
            "value": "",
            "opensWindow": true,
            "windowTimeout": 5000 // Waits up to 5s for a new window to appear
          }
    `
	helpMap["windowHandles"] = `
        - Description: Prints the list of handles for all currently open windows to the console (standard output). Useful mainly for debugging.
        - Parameters: None used.
        - Example:
          {
            "id": "...",
            "command": "windowHandles",
            "target": "",
            "value": ""
          }
    `
	helpMap["close"] = `
        - Description: Closes a specific window identified by its handle or by a name previously saved with storeWindowHandle. Unlike closeWindow, this command requires the identifier of the window to close.
        - Parameters:
          - target: The identifier of the window to close (e.g. ${popup}, handle=HANDLE_ID).
          - value: Not used.
        - Example:
          {
            "id": "...",
            "command": "close",
            "target": "${windowToClose}",
            "value": ""
          }
    `
	helpMap["setWindowMain"] = `
        - Description: Sets the handle of the window *currently* in focus as the new "main" or "root" window for webautoma. Useful if the workflow permanently shifts to a new main window.
        - Parameters: None used.
        - Example:
          {
            "id": "...",
            "command": "setWindowMain",
            "target": "",
            "value": ""
          }
    `
	helpMap["selectFrame"] = `
        - Description: Moves focus to a frame (or iframe) within the current page.
        - Parameters:
          - target: Frame identifier:
              - index=N: Numeric frame index (0-based).
              - relative=parent: Switch to parent frame.
              - relative=top: Switch to top-level page context (outside all frames).
              - String: Name or ID of the (i)frame element.
              - Empty: Reset to top-level page context (equivalent to relative=top).
          - value: Not used.
        - Example (By index):
          {
            "id": "...",
            "command": "selectFrame",
            "target": "index=0",
            "value": ""
          }
        - Example (By ID/Name):
          {
            "id": "...",
            "command": "selectFrame",
            "target": "contentFrame",
            "value": ""
          }
    `
	helpMap["selectParentFrame"] = `
        - Description: Moves focus from the current frame to its direct parent frame. Equivalent to selectFrame with target=relative=parent.
        - Parameters: None used.
        - Example:
          {
            "id": "...",
            "command": "selectParentFrame",
            "target": "",
            "value": ""
          }
    `
	helpMap["selectAlert"] = `
        - Description: Handles a JavaScript alert (alert(), confirm(), prompt() popup) that appears on the page. Reads the alert text (for logging) and then accepts or dismisses it.
        - Parameters:
          - target: Not used.
          - value: Determines the action:
              - 1: Accept alert (e.g. click "OK" in confirm).
              - 0 or omitted: Dismiss/Cancel alert (e.g. click "Cancel" in confirm).
        - Example (Accept):
          {
            "id": "...",
            "command": "selectAlert",
            "target": "",
            "value": "1"
          }
        - Example (Dismiss):
          {
            "id": "...",
            "command": "selectAlert",
            "target": "",
            "value": "0" // or omitted
          }
    `

	// --- Category: Advanced Interactions (Mouse) ---
	helpMap["rightClick"] = `
        - Description: Simulates a right mouse click on the specified element. Waits for the element to be ready.
        - Parameters:
          - target: The selector of the element to right-click on.
          - value: Not used.
          - until: (Optional) If 1 or 2, waits for the element to be ready.
        - Example:
          {
            "id": "...",
            "command": "rightClick",
            "target": "id=contextMenuArea",
            "value": "",
            "until": "1"
          }
    `
	helpMap["doubleClick"] = `
        - Description: Simulates a double left mouse click on the specified element. Waits for the element to be ready.
        - Parameters:
          - target: The selector of the element to double-click on.
          - value: Not used.
          - until: (Optional) If 1 or 2, waits for the element to be ready.
        - Example:
          {
            "id": "...",
            "command": "doubleClick",
            "target": "css=div.editable",
            "value": "",
            "until": "1"
          }
    `
	helpMap["mouseOver"] = `
        - Description: Moves the mouse cursor over the specified element, potentially triggering hover effects or tooltips. Waits for the element to be ready.
        - Parameters:
          - target: The selector of the element to move the mouse over.
          - value: Not used.
          - until: (Optional) If 1 or 2, waits for the element to be ready.
        - Example:
          {
            "id": "...",
            "command": "mouseOver",
            "target": "id=menuItem",
            "value": "",
            "until": "1"
          }
    `
	helpMap["mouseOut"] = `
        - Description: Simulates the mouse leaving the area of the specified element. Note: In the current code (doMouse), this command does not seem to execute specific actions on the WebDriver.
        - Parameters:
          - target: The selector of the element the mouse "leaves".
          - value: Not used.
        - Example:
          {
            "id": "...",
            "command": "mouseOut",
            "target": "id=menuItem",
            "value": ""
          }
    `
	helpMap["mouseDownAt"] = `
        - Description: Simulates pressing (without releasing) the left mouse button on the specified element, potentially at relative coordinates. Starts a drag operation.
        - Parameters:
          - target: The selector of the element to press the mouse on.
          - value: (Optional) Relative coordinates to the top-left corner of the element, format "X,Y" (e.g. "10,15"). If omitted, presses at the center (or driver default).
          - until: (Optional) If 1 or 2, waits for the element to be ready.
        - Example:
          {
            "id": "...",
            "command": "mouseDownAt",
            "target": "id=draggableElement",
            "value": "5,5", // Press near top-left corner
            "until": "1"
          }
    `
	helpMap["mouseMoveAt"] = `
        - Description: Moves the mouse while the left button is held down (started with mouseDownAt). Used for dragging. Requires a previous mouseDownAt on the same implicit element.
        - Parameters:
          - target: (Generally ignored, acts on the mouseDownAt element).
          - value: Relative coordinates to the top-left corner of the *original* element from mouseDownAt, format "X,Y". Indicates the position *to which* to move the mouse.
        - Example:
          {
            "id": "...",
            "command": "mouseMoveAt",
            "target": "", // Target not needed here
            "value": "100,50" // Move mouse 100px right, 50px down from drag origin
          }
    `
	helpMap["mouseMultipleMoveAt"] = `
        - Description: Similar to mouseMoveAt, but executes a sequence of consecutive relative movements while the button is pressed. Useful for simulating dragging along a path. Requires a previous mouseDownAt.
        - Parameters:
          - target: (Generally ignored).
          - value: Sequence of *incremental* relative coordinates, separated by |. Each "X,Y" pair is relative to the *previous position*. Format: "dX1,dY1|dX2,dY2|...".
        - Example:
          {
            "id": "...",
            "command": "mouseMultipleMoveAt",
            "target": "",
            "value": "50,0|0,50|-50,0" // Move 50px right, then 50px down, then 50px left
          }
    `
	helpMap["mouseUpAt"] = `
        - Description: Simulates releasing the left mouse button, completing a drag and drop operation. Requires a previous mouseDownAt.
        - Parameters:
          - target: (Generally ignored, acts on the mouseDownAt element).
          - value: (Optional) Relative coordinates to the top-left corner of the *original* element from mouseDownAt, format "X,Y". Indicates the *final* position where to release the mouse. If omitted, releases at the current position.
        - Example:
          {
            "id": "...",
            "command": "mouseUpAt",
            "target": "",
            "value": "200,100" // Release mouse 200px right, 100px down from drag origin
          }
    `

	// --- Category: Advanced Interactions (Keyboard - Custom webautoma) ---
	helpMap["actionsSendKeys"] = `
        - Description: Sends a special key press (non-alphanumeric) or a simple sequence to the active element on the page. Useful for simulating Enter, Tab, Arrow keys, etc.
        - Parameters:
          - target: (Generally not used, acts on the active element).
          - value: The name of the special key (e.g. Enter, Tab, ArrowDown, Control, Alt, Shift, F5, etc. - see wd/base/keys.go for the full list of mapped names) or a sequence of simple characters. The mapping searches KeyFromMapping in base/keys.go for special names.
        - Example (Enter):
          {
            "id": "...",
            "command": "actionsSendKeys",
            "target": "",
            "value": "Enter"
          }
        - Example (Down Arrow):
          {
            "id": "...",
            "command": "actionsSendKeys",
            "target": "",
            "value": "ArrowDown"
          }
    `
	helpMap["keys"] = `
        - Description: Allows defining complex sequences of keyboard interactions, including press (press), release (release), and pauses (pause). Useful for simulating key combinations (e.g. Ctrl+C) or specific behaviors.
        - Parameters:
          - target: (Optional) Selector of the element to send events to. If omitted, sends to the active element.
          - value: Formatted string describing the sequence, comma-separated. Each item is action:value.
              - press:KEY: Simulates pressing a key (e.g. press:Control, press:c). Uses KeyFromMapping for special keys.
              - release:KEY: Simulates releasing a key (e.g. release:Control, release:c).
              - pause:MS: Inserts a pause in milliseconds (e.g. pause:100).
        - Example (Ctrl+A, Ctrl+C):
          {
            "id": "...",
            "command": "keys",
            "target": "id=myTextArea",
            "value": "press:Control,press:a,release:a,pause:50,press:c,release:c,release:Control"
          }
        - Example (Type "test" while holding Shift):
          {
            "id": "...",
            "command": "keys",
            "target": "id=myInput",
            "value": "press:Shift,press:t,release:t,press:e,release:e,press:s,release:s,press:t,release:t,release:Shift"
          }
    `

	// --- Category: File & Download (Custom webautoma) ---
	helpMap["download"] = `
        - Description: Downloads a file directly from a specified URL and saves it to the indicated local path. Uses current browser session cookies to handle any authentication required by the server for the download.
        - Parameters:
          - target: The full URL to download the file from. If the URL starts with '/', it is considered relative to the current page base URL (obtained from the last 'open' command or navigation). Supports variables {{.variable}}.
          - value: The full path, including filename, where to save the downloaded file on the local system (e.g. /local/path/filename.pdf or C:\Download\report.xlsx). Supports variables {{.variable}}.
        - Example:
          {
            "id": "...",
            "command": "download",
            "target": "https://example.com/resources/document.zip",
            "value": "/home/user/download/archive.zip"
          }
          {
            "id": "...",
            "command": "download",
            "target": "/api/export?id={{.reportId}}", // Relative to base URL
            "value": "report_{{.reportId}}.csv"
          }
    `
	helpMap["clickDownload"] = `
        - Description: Finds an element on the page (usually an <a> link), extracts the URL from its href attribute, then downloads the file from that URL to the specified local path. Useful when the download URL is dynamic or not known in advance. Uses current session cookies.
        - Parameters:
          - target: The selector of the element (e.g. <a> link) containing the href attribute with the URL of the file to download.
          - value: The full path, including filename, where to save the downloaded file on the local system. Supports variables {{.variable}}.
          - until: (Optional) If set to 1 or 2, waits for the element to be ready before attempting to read its href attribute.
        - Example:
          {
            "id": "...",
            "command": "clickDownload",
            "target": "css=a.download-link[data-file='report']",
            "value": "/tmp/downloaded_report.pdf",
            "until": "1"
          }
    `

	// --- Category: Scrolling (Custom webautoma) ---
	helpMap["scroll"] = `
        - Description: Scrolls the page (viewport) or from a specific element by a given amount (delta X, delta Y). Useful for shifting view by a fixed amount. Can simulate smooth scrolling by splitting it into steps with intermediate pauses. Uses the WebDriver Actions API (wheel/touchpad simulation).
        - Parameters:
          - target: (Optional) The selector of the element from which to calculate the scroll starting point. If omitted, scroll starts from the coordinates specified in the command's x/y fields (default 0,0) plus any offsetX/offsetY.
          - value: Specifies the displacement (delta) and optionally steps and interval. Format: "DeltaX,DeltaY[,NumberSteps,IntervalMs]".
              - DeltaX, DeltaY: Horizontal and vertical shift in pixels (negative for left/up, positive for right/down).
              - NumberSteps: (Optional) Number of steps to divide the total scroll into.
              - IntervalMs: (Optional, requires NumberSteps) Milliseconds pause between steps.
          - x, y: (Optional, used only if target is omitted) Absolute X, Y coordinates in the viewport where the scroll action begins. Default: 0,0.
          - offsetX, offsetY: (Optional) Offset in pixels to add to starting coordinates (both x/y and those calculated from target element).
        - Example (Scroll viewport 500px down):
          {
            "id": "...",
            "command": "scroll",
            "target": "", // Viewport scroll
            "value": "0,500"
          }
        - Example (Smooth scroll viewport 1000px down in 10 steps):
          {
            "id": "...",
            "command": "scroll",
            "target": "",
            "value": "0,1000,10,50" // 10 steps, 50ms pause between steps
          }
        - Example (Scroll starting from bottom of a header element):
          {
            "id": "...",
            "command": "scroll",
            "target": "id=mainHeader",
            "value": "0,300", // Scroll 300px down
            "offsetY": 50 // Starting 50px below header
          }
    `
	helpMap["scrollTo"] = `
        - Description: Scrolls the page so that a specific point *within* a target element becomes visible in the viewport. If the element is not initially visible, it first tries to bring it into view. Can simulate smooth scrolling. Uses the WebDriver Actions API.
        - Parameters:
          - target: (Optional) The selector of the target element to scroll towards. If omitted, uses the currently active element (the one with focus).
          - value: Specifies coordinates *relative to the top-left corner of the target element* and optionally steps and interval. Format: "TargetX,TargetY[,NumberSteps,IntervalMs]".
              - TargetX, TargetY: X, Y coordinates *inside* the target element to bring into view. "0,0" corresponds to the top-left corner of the element.
              - NumberSteps: (Optional) Number of steps to reach the position.
              - IntervalMs: (Optional, requires NumberSteps) Milliseconds pause between steps.
          - offsetX, offsetY: (Optional) Offset in pixels to add to TargetX, TargetY coordinates specified in value.
          - until: (Optional) If set to 1 or 2, waits for the target element to be ready before attempting the scroll.
        - Example (Scroll to top of footer):
          {
            "id": "...",
            "command": "scrollTo",
            "target": "id=pageFooter",
            "value": "0,0", // Bring corner 0,0 of footer into view
            "until": "1"
          }
        - Example (Smoothly scroll to a specific point inside a div):
          {
            "id": "...",
            "command": "scrollTo",
            "target": "css=div.scrollable-content",
            "value": "0,500,10,50" // Bring point Y=500 inside div into view, in 10 steps
          }
    `

	// --- Category: Advanced Custom (webautoma) ---
	helpMap["otp"] = `
        - Description: Retrieves a One-Time Password (OTP) from an email account. Connects to the specified IMAP server, searches for the most recent email matching the criteria (subject, maximum age), extracts the OTP code from the email body using a regular expression, and stores it internally. The retrieved OTP can be used in subsequent commands (e.g. 'type') using the variable {{.otp}}.
        - Parameters:
          - target: Connection string to IMAP server. Format: protocol[authMode]://user:password@server:port
              - protocol: Can be tls (recommended), starttls, or insecure.
              - [authMode]: (Optional) Specify [oauth] or [oauth2] if using OAuth/OAuth2 authentication instead of direct password.
              - user: Username for IMAP login.
              - password: Password or OAuth token.
              - server: IMAP server address.
              - port: IMAP server port (e.g. 993 for TLS, 143 for StartTLS/Insecure).
              - TLS Example: tls://my.email@example.com:MyPassword@imap.example.com:993
              - OAuth2 Example: tls[oauth2]://user@gmail.com:OAuth2AccessToken@imap.gmail.com:993
          - value: Configuration string for search and extraction, with parameters separated by |||. Format: "SubjectRegExp|||BodyRegExpWithOTPCaptureGroup[|||CheckIntervalSec[|||MaxEmailAgeMin]]"
              - SubjectRegExp: Regular expression (standard Go) to identify the email subject containing the OTP (e.g. ^Verification code.*$).
              - BodyRegExpWithOTPCaptureGroup: Regular expression (standard Go) applied to the email body to extract the OTP. **Must** contain a capture group in parentheses () that isolates the OTP code exactly (e.g. Your OTP code is ([0-9]{6})\., captures 6 digits).
              - CheckIntervalSec: (Optional, default: 60) Maximum number of seconds during which webautoma will attempt to retrieve the email (periodically checking the mailbox).
              - MaxEmailAgeMin: (Optional, default: 5) Maximum age in minutes that the email can have to be considered valid.
        - Result: The extracted OTP is saved in the internal variable otp, accessible as {{.otp}} in the value fields of subsequent commands.
        - Example:
          {
            "id": "retrieveOTP",
            "command": "otp",
            "target": "tls://user@example.com:secretPassword@imap.example.com:993",
            "value": "Your one-time code Example Corp|||verification code: ([A-Z0-9]+)|||90|||3"
            // Searches email with subject "Your one-time code Example Corp"
            // Extracts alphanumeric code from body (e.g. "verification code: XY78Z1")
            // Attempts for 90 seconds, valid emails if newer than 3 minutes
          }
          {
            "id": "enterOTP",
            "command": "type",
            "target": "id=otpField",
            "value": "{{.otp}}" // Uses the OTP retrieved in previous step
          }
    `

	// --- Category: Debug and Meta-Commands (Custom webautoma) ---
	helpMap["disableError"] = `
        - Description: Temporarily disables execution interruption in case of errors in subsequent commands. Useful for attempting actions that might fail without stopping the entire test. The error will still be logged (unless debug is also disabled). Use enableError to re-enable normal behavior.
        - Parameters: None used.
        - Example:
          {
            "id": "...",
            "command": "disableError",
            "target": "",
            "value": ""
          }
    `
	helpMap["enableError"] = `
        - Description: Re-enables execution interruption in case of an error, canceling the effect of a previous disableError. This is the default behavior.
        - Parameters: None used.
        - Example:
          {
            "id": "...",
            "command": "enableError",
            "target": "",
            "value": ""
          }
    `
	helpMap["disableDebug"] = `
        - Description: Disables detailed log output (debug level) generated by the webautoma adapter during execution.
        - Parameters: None used.
        - Example:
          {
            "id": "...",
            "command": "disableDebug",
            "target": "",
            "value": ""
          }
    `
	helpMap["enableDebug"] = `
        - Description: Re-enables detailed log output (debug level). Useful if previously disabled.
        - Parameters: None used.
        - Example:
          {
            "id": "...",
            "command": "enableDebug",
            "target": "",
            "value": ""
          }
    `
	helpMap["status"] = `
        - Description: Requests and prints WebDriver server status information (version, OS, availability) to the console (standard output). Useful for diagnosis.
        - Parameters: None used.
        - Example:
          {
            "id": "...",
            "command": "status",
            "target": "",
            "value": ""
          }
    `
	helpMap["execId"] = `
        - Description: Sets a global identifier for the entire current execution. This ID will be included in generated JSON logs, useful for correlating events across different runs.
        - Parameters:
          - target: The string to use as execution ID.
          - value: Not used.
        - Example:
          {
            "id": "...",
            "command": "execId",
            "target": "prod_run_evening",
            "value": ""
          }
    `
	helpMap["probe"] = `
        - Description: Sets an identifier for the "probe" or specific instance of the test in progress. This ID will be included in JSON logs, useful for distinguishing results when multiple instances of the same test run in parallel or to identify specific monitoring points.
        - Parameters:
          - target: The string to use as probe/instance ID.
          - value: Not used.
        - Example:
          {
            "id": "...",
            "command": "probe",
            "target": "monitor_login_rome",
            "value": ""
          }
    `
	helpMap["setTimeout"] = `
        - Description: Sets the maximum time (implicit timeout) in milliseconds that WebDriver will wait when searching for an element (findElement, findElements) before returning an "element not found" error.
        - Parameters:
          - target: The wait time in milliseconds.
          - value: Not used.
        - Example:
          {
            "id": "...",
            "command": "setTimeout",
            "target": "30000", // Sets timeout to 30 seconds
            "value": ""
          }
    `
	helpMap["activeElement"] = `
        - Description: Identifies the currently active element (the one with focus) on the page and prints its details (tag, text, state) to the console (standard output). Useful for debugging to understand where keyboard/interaction focus is.
        - Parameters: None used.
        - Example:
          {
            "id": "...",
            "command": "activeElement",
            "target": "",
            "value": ""
          }
    `
	helpMap["pageSource"] = `
        - Description: Retrieves the entire HTML source of the currently displayed page and prints it to the console (standard output). Useful for advanced debugging of DOM structure.
        - Parameters: None used.
        - Example:
          {
            "id": "...",
            "command": "pageSource",
            "target": "",
            "value": ""
          }
    `
	helpMap["noop"] = `
        - Description: "No Operation" command. Performs no action. Can be useful as a placeholder or to insert comments into the .side file (using Selenium IDE standard comment field; although webautoma does not actively read it, it can be useful for anyone reading the file).
        - Parameters: None used.
        - Example:
          {
            "id": "...",
            "command": "noop",
            "target": "",
            "value": "",
            "comment": "Checkout section begins here"
          }
    `

	// Add other commands here if necessary...

	return &Provider{
		commandHelp: helpMap,
	}
}

// GetCommand searches and returns the help string for the specified command.
// Returns an error if the command is not found in the help map.
func (h *Provider) GetCommand(commandName string) (string, error) {
	helpText, found := h.commandHelp[commandName]
	if !found {
		return "", fmt.Errorf("no help found for command: '%s'", commandName)
	}
	// Remove leading/trailing whitespaces from the text block for cleanliness
	return strings.TrimSpace(helpText), nil
}

// GetSupportedCommands returns a sorted list of command names
// for which help is available.
func (h *Provider) GetSupportedCommands() []string {
	keys := make([]string, 0, len(h.commandHelp))
	for k := range h.commandHelp {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (h *Provider) GetFormatted(commandName string) (string, error) {
	helpText, err := h.GetCommand(commandName)
	if err != nil {
		return "", err
	}
	lines := strings.Split(helpText, "\n")
	minIndent := 999
	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " ")
		if len(trimmed) == 0 {
			continue
		}
		indent := len(line) - len(trimmed)
		if indent < minIndent {
			minIndent = indent
		}
	}
	var builder strings.Builder
	for i, line := range lines {
		if len(line) >= minIndent {
			builder.WriteString(line[minIndent:])
		} else {
			builder.WriteString(line) // Keep empty lines or lines with smaller indentation
		}
		if i < len(lines)-1 {
			builder.WriteString("\n")
		}
	}
	return builder.String(), nil
}
