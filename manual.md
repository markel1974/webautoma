Introduction

Webautoma is an advanced tool for web application automation and testing, written entirely in Go. Its primary purpose is to execute tests defined in the standard Selenium IDE JSON format (`.side` files), communicating with browsers via the WebDriver protocol (primarily supporting ChromeDriver, but potentially extensible).

Beyond executing standard Selenium IDE commands, `webautoma` significantly extends automation capabilities by introducing:

* A rich set of **custom commands** to handle complex scenarios such as:
    * Precise execution time measurement (Timers).
    * Management of dynamic variables extracted from the page (Stack).
    * Automatic retrieval of One-Time Passwords (OTP) from IMAP email accounts.
    * File downloads.
    * Advanced scrolling of the page or specific elements.
    * Complex mouse and keyboard interactions.
* An **Interactive Console (SAM - Step-Aside-Mode)** that allows executing `.side` tests step-by-step, inspecting state, modifying commands on-the-fly, and debugging much more effectively than standard execution.
* **Detailed logging** in JSON format, including execution times for each step, errors, network error counts, and captured HTTP headers (if configured).
* **Screenshot capture** automatically on error or on demand, with metadata saved in a separate JSON file.
* **Support for Variables and Templating**, allowing test parameterization via the command line or external files.
* A **Server Mode** to execute tests via HTTP requests (useful for integrations).
* **Optional lifecycle management** of the WebDriver service (e.g., starting ChromeDriver).

**Who is webautoma for?**

Webautoma is designed for developers, testers, and automation engineers who:

* Use or want to use the Selenium IDE `.side` format but require more advanced features.
* Need interactive debugging capabilities for their automation scripts.
* Must automate complex scenarios involving email OTPs or precise performance measurements.
* Prefer a stand-alone tool written in Go, easily distributable, and with potential performance advantages.

This manual guides the installation, configuration, and usage of all `webautoma` features.