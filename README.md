# Liver

A recursive file watcher that works by attaching itself to processes

## Installation

```bash
$ go install github.com/mnafees/liver@latest
```

## Getting started

To use Liver, you need a `liver.json` file in the root of your directory with the following contents
```json
{
    "paths": [
        "/some/path/or/file/to/watch",
        "/other/path/index.js"
    ],
    "procs": {
        "/some/path": "go run main.go",
        "/other/path": "node index.js"
    }
}
```

Then call Liver such as
```bash
$ liver
```
