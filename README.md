# Liver

A dead simple recursive file watcher that works by attaching itself to processes

## Installation

1. If you have Go installed, then
```bash
$ go install github.com/mnafees/liver@latest
```
2. Or else, directly download from the [latest release](https://github.com/mnafees/liver/releases/latest) for your respective OS and architecture and place the binary in a location such as `/usr/local/bin` which is part of your `PATH`

## Getting started

To use Liver, you need a `liver.json` file in the root of your directory with the following contents
```json
{
    "paths": [
        "/some/path/or/file/to/watch",
        "/other/path/index.js"
    ],
    "procs": {
        "/some/path": [ "go run main.go" ],
        "/other/path": [
            "node index.js",
            "node index2.js"
        ]
    }
}
```

Then call Liver such as
```bash
$ liver
```
