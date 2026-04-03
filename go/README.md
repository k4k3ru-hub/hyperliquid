# Hyperliquid SDK for Go

[![License](https://img.shields.io/github/license/k4k3ru-hub/hyperliquid/go)](./LICENSE)

This is a Go library for interacting with Hyperliquid, providing convenient utilities for developers building applications with Hyperliquid's API.


## Features

- Easy-to-use Go CLI for Hyperliquid.
- Support for key Hyperliquid functionalities such as order execution, market data retrieval, and account management.
- Simple and efficient integration with Hyperliquid API.


## Installation

### CLI

```
git clone https://github.com/k4k3ru-hub/hyperliquid.git
cd hyperliquid/go
go build -o build/hyperliquid-cli cli/main.go
ls -l build/hyperliquid-cli
```

### Module

```
import "github.com/k4k3ru-hub/hyperliquid/go"
```


## Usage

Currently supported features:
- REST:
  - metaAndAssetCtxs: Retrieve perpetuals asset contexts (includes mark price, current funding, open interest, etc.)
- WebSocket:
  - Under development...


### Using as a CLI

Hyperliquid SDK for Go can be executed as a CLI command.

- REST:
  - metaAndAssetCtxs: `./build/hyperliquid-cli rest metaAndAssetCtxs`


### Using as a Go Module

The SDK can also be imported and used as a Go module.

1. Install the package

```
go get github.com/k4k3ru-hub/hyperliquid/go
```

2. Import the package for the function you want to use.

- REST:
  - metaAndAssetCtxs: `import "github.com/k4k3ru-hub/hyperliquid/go/rest/info/meta_and_asset_ctxs"`

3. Initialize the client for the function you want to use.

- REST:
  - metaAndAssetCtxs: `c := meta_and_asset_ctxs.NewMetaAndAssetCtxsClient()`

4. Send a API request.

```go
result, err := c.Send()
```


## References

- [Hyperliquid Docs](https://hyperliquid.gitbook.io/hyperliquid-docs)


## Support me
I am a Japanese developer, and your support is a great encouragement for my work!
In addition to support, feel free to reach out with comments, feature requests, or development inquiries!

Thank you for your support😊

[![Support on Ko-fi](https://img.shields.io/badge/Ko--fi-Support%20Me-blue?style=flat-square&logo=ko-fi)](https://ko-fi.com/k4k3ru)
[![Support on Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-Support%20Me-yellow?style=flat-square&logo=buy-me-a-coffee)](https://buymeacoffee.com/k4k3ru)


## License
This repository is open-source and distributed under the MIT License.
