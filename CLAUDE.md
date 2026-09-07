# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this project is

A Go reimplementation of Hikvision's SADPTool — a command-line tool that discovers Hikvision devices (IPCs/NVRs) on the local network by sending a UDP multicast probe and writing the replies to `device.json`. The protocol was reverse-engineered from a Wireshark packet capture of the official SADPTool.

Note: the README and almost all code comments are written in Chinese.

## Build & run

```bash
go run main.go       # run discovery
go build ./...       # compile
go vet ./...         # static analysis
```

There are no test files. The single dependency is `github.com/google/uuid` (Go 1.20+).

Running the program sends a probe, listens ~2 seconds (see `SetReadDeadline` in `main.go`), and writes results to `device.json` (a JSON array of devices) then exits. `device.json` is gitignored (`*.json`), so it only exists after a run. To re-run, either delete it first or just run again — it's overwritten each time.

## Architecture

Two packages:

- **`main`** (`main.go`) — all program logic: opens a UDP socket, sends an XML probe to the multicast group `239.255.255.250:37020`, then loops reading `ProbeMatch` XML replies until the 2-second read deadline expires, appending each reply to a device list and rewriting `device.json` after each read. The XML probe is a `model.Probe` marshaled with an uppercase UUID and `Types: "inquiry"`.
- **`model`** (`model/model.go`) — XML wire types only (no logic):
  - `Probe` — the outbound discovery request.
  - `Device` — a discovered device (`<ProbeMatch>`, root element bound via `XMLName`), with fields mirroring Hikvision's reply XML.
  - `DeviceList` — a mutex-guarded slice. Note: `main.go` marshals `deviceList.Devices` (the bare slice), **not** the `DeviceList` struct, so the output is a JSON array.

## Gotchas

- **Flaky field types in `Device`**: Hikvision's firmware sends `"flase"` (a typo for "false") in some boolean fields, which fails `strconv.ParseBool`. Fields affected by this — `HCPlatformEnable`, `SupportEzvizUnbind`, `SupportIPv6`, `SupportModifyIPv6` — are deliberately typed as `string` rather than `bool` (see the comment in `model.go`). Keep unknown reply fields either `string` or guarded accordingly; don't change these types back to `bool`.
- The receive loop ends by returning on read error/timeout rather than a controlled break, which is why the program has no explicit exit after the deadline.
- `ParseBool` failures inside the reply are only logged (`xml转换结构体异常`), not fatal; the rest of the device is still appended.
