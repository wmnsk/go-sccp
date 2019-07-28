# go-sccp: SCCP in Golang

Package sccp provides simple and painless handling of SCCP(Signaling Connection Control Part) in SS7/SIGTRAN stack, implemented in the Go Programming Language.

[![CircleCI](https://circleci.com/gh/wmnsk/go-sccp.svg?style=shield)](https://circleci.com/gh/wmnsk/go-sccp)
[![GolangCI](https://golangci.com/badges/github.com/wmnsk/go-sccp.svg)](https://golangci.com/r/github.com/wmnsk/go-sccp)
[![GoDoc](https://godoc.org/github.com/wmnsk/go-sccp?status.svg)](https://godoc.org/github.com/wmnsk/go-sccp)
[![GitHub](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/wmnsk/go-sccp/blob/master/LICENSE)

## Disclaimer

This is still an experimental project, and currently in its very early stage of development. Any part of implementations(including exported APIs) may be changed before released as v1.0.0.

## Getting started

The following package should be installed before getting started.

```shell-session
go get -u github.com/google/go-cmp
go get -u github.com/pascaldekloe/goe
go get -u github.com/pkg/errors
```

If you use Go 1.11+, you can also use Go Modules.


```shell-session
GO111MODULE=on go [test | build | run | etc...]
```

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
| Extended unitdata              | XUDT         | 4.18      | -          |
| Extended unitdata service      | XUDTS        | 4.19      | -          |
| Long unitdata                  | LUDT         | 4.20      | -          |
| Long unitdata service          | LUDTS        | 4.21      | -          |

### Parameters

| Parameter name              | Reference | Supported? |
| --------------------------- | --------- | ---------- |
| End of optional parameters  | 3.1       |            |
| Destination local reference | 3.2       |            |
| Source local reference      | 3.3       |            |
| Called party address        | 3.4       | Yes        |
| Calling party address       | 3.5       | Yes        |
| Protocol class              | 3.6       | Yes        |
| Segmenting/reassembling     | 3.7       |            |
| Receive sequence number     | 3.8       |            |
| Sequencing/segmenting       | 3.9       |            |
| Credit                      | 3.10      |            |
| Release cause               | 3.11      |            |
| Return cause                | 3.12      |            |
| Reset cause                 | 3.13      |            |
| Error cause                 | 3.14      |            |
| Refusal cause               | 3.15      |            |
| Data                        | 3.16      | Yes        |
| Segmentation                | 3.17      |            |
| Hop counter                 | 3.18      |            |
| Importance                  | 3.19      |            |
| Long data                   | 3.20      |            |

## Author(s)

Yoshiyuki Kurauchi ([My Website](https://wmnsk.com/) / [Twitter](https://twitter.com/wmnskdmms))

I'm always open to welcome co-authors! Please feel free to talk to me.

## LICENSE

[MIT](https://github.com/wmnsk/go-sccp/blob/master/LICENSE)
