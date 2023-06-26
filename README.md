# Liver

[![License](https://img.shields.io/github/license/mnafees/liver)](https://github.com/mnafees/liver/blob/main/LICENSE)
![Code size](https://img.shields.io/github/languages/code-size/mnafees/liver)
[![Go report card](https://goreportcard.com/badge/github.com/mnafees/liver)](https://goreportcard.com/report/github.com/mnafees/liver)
[![Go version](https://img.shields.io/github/go-mod/go-version/mnafees/liver.svg)](https://github.com/mnafees/liver)



A dead simple recursive file watcher and live reloading utility that works by attaching itself to processes and matching them to files.

## Installation

### 1. `go install`

```bash
go install github.com/mnafees/liver@latest
```

### 2. Manually

**macOS**

```bash
# arm64 (Apple Silicon)
curl -L -o liver.tar.gz "https://github.com/mnafees/liver/releases/latest/download/Liver_Darwin_arm64.tar.gz" && tar xvf liver.tar.gz Liver && mv Liver liver && sudo install -c -m 0755 liver /usr/local/bin && rm -f liver.tar.gz

# x86_64 (Intel)
curl -L -o liver.tar.gz "https://github.com/mnafees/liver/releases/latest/download/Liver_Darwin_x86_64.tar.gz" && tar xvf liver.tar.gz Liver && mv Liver liver && sudo install -c -m 0755 liver /usr/local/bin && rm -f liver.tar.gz
```

**Linux**

```bash
# arm64
curl -L -o liver.tar.gz "https://github.com/mnafees/liver/releases/latest/download/Liver_Linux_arm64.tar.gz" && tar xvf liver.tar.gz Liver && mv Liver liver && sudo install -c -m 0755 liver /usr/local/bin && rm -f liver.tar.gz

# x86_64
curl -L -o liver.tar.gz "https://github.com/mnafees/liver/releases/latest/download/Liver_Linux_x86_64.tar.gz" && tar xvf liver.tar.gz Liver && mv Liver liver && sudo install -c -m 0755 liver /usr/local/bin && rm -f liver.tar.gz

# i386
curl -L -o liver.tar.gz "https://github.com/mnafees/liver/releases/latest/download/Liver_Linux_i386.tar.gz" && tar xvf liver.tar.gz Liver && mv Liver liver && sudo install -c -m 0755 liver /usr/local/bin && rm -f liver.tar.gz
```

**Windows**

Download the latest Windows binaries from https://github.com/mnafees/liver/releases/latest

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
liver
```
