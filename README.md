[godoclink]: https://godoc.org/github.com/sarvalabs/go-moi-identifiers
[latestrelease]: https://github.com/sarvalabs/go-moi-identifiers/releases/latest
[issueslink]: https://github.com/sarvalabs/go-moi-identifiers/issues
[pullslink]: https://github.com/sarvalabs/go-moi-identifiers/pulls

![go version](https://img.shields.io/github/go-mod/go-version/sarvalabs/go-moi-identifiers?style=for-the-badge)
![license](https://img.shields.io/badge/license-MIT%2FApache--2.0-informational?style=for-the-badge)
[![go docs](http://img.shields.io/badge/go-documentation-blue.svg?style=for-the-badge)][godoclink]
[![latest tag](https://img.shields.io/github/v/tag/sarvalabs/go-moi-identifiers?color=blue&label=latest%20tag&sort=semver&style=for-the-badge)][latestrelease]

![ci status](https://img.shields.io/github/actions/workflow/status/sarvalabs/go-moi-identifiers/ci.yaml?label=CI&style=for-the-badge)
[![issue count](https://img.shields.io/github/issues/sarvalabs/go-moi-identifiers?style=for-the-badge&color=yellow)][issueslink]
[![pulls count](https://img.shields.io/github/issues-pr/sarvalabs/go-moi-identifiers?style=for-the-badge&color=brightgreen)][pullslink]

# MOI Identifiers
**go-moi-identifiers** is a package that contains implementations for all primitive identifiers
used in the MOI Protocol such as the `Address`, `LogicID` and `AssetID` standards. The `KramaID` 
primitive will also be implemented soon.

### Specifications & Implementations
| Identifier |                                              Specification                                              | Implemented |
|:----------:|:-------------------------------------------------------------------------------------------------------:|:-----------:|
| `Address`  |                                                _pending_                                                |     Yes     |
| `AssetID`  | [Asset ID Spec](https://sarvalabs.notion.site/Asset-ID-Standard-e4fcd9151e7d4e7eb2447f1d8edf4672?pvs=4) |     Yes     |
| `LogicID`  | [Logic ID Spec](https://sarvalabs.notion.site/Logic-ID-Standard-174a2cc6e3dc42e4bbf4dd708af0cd03?pvs=4) |     Yes     |
| `KramaID`  |                                                _pending_                                                |     No      |

## Install
Install the latest [release](https://github.com/sarvalabs/go-moi-engineio/releases) using the following command
```sh
go get -u github.com/sarvalabs/go-moi-identifiers
```

## Contributing
Unless you explicitly state otherwise, any contribution intentionally submitted
for inclusion in the work by you, as defined in the Apache-2.0 license, shall be
dual licensed as below, without any additional terms or conditions.

## License
&copy; 2023 Sarva Labs Inc. & MOI Protocol Developers.

This project is licensed under either of
- [Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0) ([`LICENSE-APACHE`](LICENSE-APACHE))
- [MIT license](https://opensource.org/licenses/MIT) ([`LICENSE-MIT`](LICENSE-MIT))

at your option.

The [SPDX](https://spdx.dev) license identifier for this project is `MIT OR Apache-2.0`.