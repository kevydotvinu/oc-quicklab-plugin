# oc-quicklab-plugin
The oc-quicklab-plugin is a CLI plugin that works with oc/kubectl and gives the quicklab shared cluster information in the CLI.

## Prerequisites
- `oc` binary

## Installation
```bash
git clone https://github.com/kevydotvinu/oc-quicklab-plugin.git
cd oc-quicklab-plugin
go build -o oc-quicklab main.go
mv oc-quicklab /usr/local/bin/oc-quicklab
oc quicklab
```

## Usage
- List quicklab shared clusters
```bash
oc quicklab list
```