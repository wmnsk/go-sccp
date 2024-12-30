# go-sccp

Package sccp provides simple and painless handling of SCCP (Signaling Connection Control Part) in SS7/SIGTRAN stack, implemented in the Go Programming Language.

[![CI status](https://github.com/wmnsk/go-sccp/actions/workflows/go.yml/badge.svg)](https://github.com/wmnsk/go-sccp/actions/workflows/go.yml)
[![golangci-lint](https://github.com/wmnsk/go-sccp/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/wmnsk/go-sccp/actions/workflows/golangci-lint.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/wmnsk/go-sccp.svg)](https://pkg.go.dev/github.com/wmnsk/go-sccp)
[![GitHub](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/wmnsk/go-sccp/blob/master/LICENSE)

## Disclaimer

This is still an experimental project, and currently in its very early stage of development. Any part of implementations(including exported APIs) may be changed before released as v1.0.0.

## Getting started

Run `go mod tidy` to download the dependency, and you're ready to start developing.

## Supported Features

### Message Types

| Message type                   | Abbreviation | Reference | Supported? |
| ------------------------------ | ------------ | --------- | ---------- |
| Connection request             | CR           | 4.2       | -          |
| Connection confirm             | CC           | 4.3       | -          |
| Connection refused             | CREF         | 4.4       | -          |
| Released                       | RLSD         | 4.5       | -          |
| Release complete               | RLC          | 4.6       | -          |
| Data form 1                    | DT1          | 4.7       | -          |
| Data form 2                    | DT2          | 4.8       | -          |
| Data acknowledgement           | AK           | 4.9       | -          |
| Unitdata                       | UDT          | 4.10      | Yes        |
| Unitdata service               | UDTS         | 4.11      | -          |
| Expedited data                 | ED           | 4.12      | -          |
| Expedited data acknowledgement | EA           | 4.13      | -          |
| Reset request                  | RSR          | 4.14      | -          |
| Reset confirm                  | RSC          | 4.15      | -          |
| Protocol data unit error       | ERR          | 4.16      | -          |
| Inactivity test                | IT           | 4.17      | -          |
| Extended unitdata              | XUDT         | 4.18      | Yes        |
| Extended unitdata service      | XUDTS        | 4.19      | -          |
| Long unitdata                  | LUDT         | 4.20      | -          |
| Long unitdata service          | LUDTS        | 4.21      | -          |

### Parameters

| Parameter name              | Reference | Supported? |
| --------------------------- | --------- | ---------- |
| End of optional parameters  | 3.1       | Yes        |
| Destination local reference | 3.2       | Yes        |
| Source local reference      | 3.3       | Yes        |
| Called party address        | 3.4       | Yes        |
| Calling party address       | 3.5       | Yes        |
| Protocol class              | 3.6       | Yes        |
| Segmenting/reassembling     | 3.7       | Yes        |
| Receive sequence number     | 3.8       | Yes        |
| Sequencing/segmenting       | 3.9       | Yes        |
| Credit                      | 3.10      | Yes        |
| Release cause               | 3.11      | Yes        |
| Return cause                | 3.12      | Yes        |
| Reset cause                 | 3.13      | Yes        |
| Error cause                 | 3.14      | Yes        |
| Refusal cause               | 3.15      | Yes        |
| Data                        | 3.16      | Yes        |
| Segmentation                | 3.17      | Yes        |
| Hop counter                 | 3.18      | Yes        |
| Importance                  | 3.19      | Yes        |
| Long data                   | 3.20      | Yes        |

## Author(s)

Yoshiyuki Kurauchi ([Website](https://wmnsk.com/)) and [contributors](https://github.com/wmnsk/go-sccp/graphs/contributors).

## LICENSE

[MIT](https://github.com/wmnsk/go-sccp/blob/master/LICENSE)
