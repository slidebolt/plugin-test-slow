### `plugin-test-slow` repository

#### Project Overview

This repository contains the `plugin-test-slow`, a simple plugin designed for testing the core Slidebolt system.

#### Architecture

This is a minimal Go plugin that implements the `runner.Plugin` interface. Its sole purpose is to simulate a plugin that is slow to initialize. It does this by intentionally sleeping for 1 second during the `OnInitialize` lifecycle hook.

This plugin is used to test how the Slidebolt launcher and other parts of the system handle slow-starting plugins, ensuring that they do not cause timeouts or other issues.

The plugin does not create any devices or entities, nor does it handle any commands or events.

#### Key Files

| File | Description |
| :--- | :--- |
| `go.mod` | Defines the Go module and its dependencies on the `sdk-runner` and `sdk-types`. |
| `main.go` | Contains the complete, minimal implementation of the slow-loading test plugin. |

#### Available Commands

This plugin does not handle any commands. It is intended for internal testing of the Slidebolt system.
