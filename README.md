## oc-quicklab
[![licence](https://img.shields.io/github/license/kevydotvinu/oc-quicklab-plugin)](https://github.com/kevydotvinu/oc-quicklab-plugin/blob/master/LICENSE)
[![goversion](https://img.shields.io/github/go-mod/go-version/kevydotvinu/oc-quicklab-plugin)](https://github.com/kevydotvinu/oc-quicklab-plugin/blob/master/go.mod)
[![downloads](https://img.shields.io/github/downloads/kevydotvinu/oc-quicklab-plugin/total)](https://github.com/kevydotvinu/oc-quicklab-plugin/releases)
[![release](https://github.com/kevydotvinu/oc-quicklab-plugin/actions/workflows/release.yml/badge.svg)](https://github.com/kevydotvinu/oc-quicklab-plugin/actions/workflows/release.yml)
[![openshiftplugin](https://img.shields.io/badge/openshift%20cli-plug--in-orange)](https://docs.openshift.com/container-platform/latest/cli_reference/openshift_cli/extending-cli-plugins.html)

### OpenShift CLI plug-in for Quicklab
The oc-quicklab-plugin is a CLI plug-in that works with oc/kubectl and gives the quicklab shared cluster information in the CLI.

### Demo
[![asciicast](https://asciinema.org/a/hZ5EQRwJG2oLB9WOlRDIr195V.svg)](https://asciinema.org/a/hZ5EQRwJG2oLB9WOlRDIr195V)

### Prerequisites
- [OpenShift CLI](https://access.redhat.com/downloads/content/290)
- [Chromium](https://www.chromium.org/) or [Google Chrome Web Browser](https://www.google.com/chrome/)
- Configure `Integrated Authentication` and `Kerberos Credentials Delegation` in Chrome browser - Add the below in `/etc/opt/chrome/policies/managed/redhat-corp.json` file.
```bash
{
    "AuthServerWhitelist": "*.redhat.com",
    "AuthNegotiateDelegateWhitelist": "*.redhat.com"
}
```

### Installation
Either use [prebuilt binary](https://github.com/kevydotvinu/oc-quicklab-plugin/releases):
```bash
$ URL=(curl -s https://api.github.com/repos/kevydotvinu/oc-quicklab-plugin/releases/latest)
$ JQ=(jq -r '.assets | .[] | select(.name | contains("linux")) | .browser_download_url')
$ curl -#Lo oc-quicklab-plugin $(curl -s ${URL} | ${JQ})
$ sudo mv oc-quicklab /usr/local/bin/
$ sudo chmod +x /usr/local/bin/oc-quicklab
$ oc quicklab help
```
Or build the binary for an OS and ARCH from [source](https://github.com/kevydotvinu/oc-quicklab-plugin):
```bash
$ git clone https://github.com/kevydotvinu/oc-quicklab-plugin.git
$ cd oc-quicklab-plugin
$ go build -o oc-quicklab main.go
$ sudo mv oc-quicklab /usr/local/bin
$ oc quicklab
```

### Usage
##### Help
```bash
$ oc quicklab help
```
##### List Quicklab shared clusters
```bash
$ oc quicklab list
```
##### Login into Quicklab shared cluster
```bash
$ oc quicklab login --cluster <cluster-name>
```

### Troublshooting
If browser winodw asks for the credentials, please check the Kerberos tickets held in a credentials cache and its validity.
```bash
klist
```
Login with credential if expires.
```bash
kinit <username>
```
