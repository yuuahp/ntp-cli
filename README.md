# ntp-cli

A tiny CLI tool to call NTP servers.

## Features

- ⏰ Get the current time from NTP servers.
- 📝 Supports multiple time formats.
- 🤫 Quiet mode to suppress output.

## Install

- Build it yourself: run `go build -o ntp-cli`
- Download from [GitHub Releases](https://github.com/yuuahp/ntp-cli/releases/latest)

## Usage

```bash
./ntp-cli [flags]
```

## Examples

Please note that `--address` conflicts with `--hostname` and `--port`.

```bash
$ ./ntp-cli
Calling NTP server at pool.ntp.org:123...
Current time: Tue Apr  1 23:50:00 JST 2025

$ ./ntp-cli -h time.apple.com
Calling NTP server at time.apple.com:123...
Current time: Tue Apr  1 23:50:59 JST 2025

$ ./ntp-cli -a time.apple.com:123 -q
$ ./ntp-cli -h time.apple.com -p 123 -q
Tue Apr  1 23:53:09 JST 2025

$ ./ntp-cli -f RFC3339 -q
2025-04-01T23:57:00+09:00

$ ./ntp-cli -h time.apple.com -a pool.ntp.org:123
invalid arguments: you can either specify an address or a hostname, but not both
```

## Flags

- `-h`, `--help`: Show help message.
- `-a <string>`, `--address <string>`:  
  The address of the NTP server. (default: `pool.ntp.org:123`)
- `-h <string>`, `--hostname <string>`:  
  The hostname of the NTP server. (default: `pool.ntp.org`)
- `-p <int>`, `--port <int>`:  
  The port of the NTP server. (default: `123`)
- `-f <string>`, `--format <string>`:  
  The format of the time. (default: `RFC3339`)
  <details>
  <summary>Available formats:</summary>

  | Format      | Example                             |
  |-------------|-------------------------------------|
  | Layout      | 01/02 03:04:05PM '06 -0700          |
  | ANSIC       | Mon Jan _2 15:04:05 2006            |
  | UnixDate    | Mon Jan _2 15:04:05 MST 2006        |
  | RubyDate    | Mon Jan 02 15:04:05 -0700 2006      |
  | RFC822      | 02 Jan 06 15:04 MST                 |
  | RFC822Z     | 02 Jan 06 15:04 -0700               |
  | RFC850      | Monday, 02-Jan-06 15:04:05 MST      |
  | RFC1123     | Mon, 02 Jan 2006 15:04:05 MST       |
  | RFC1123Z    | Mon, 02 Jan 2006 15:04:05 -0700     |
  | RFC3339     | 2006-01-02T15:04:05Z07:00           |
  | RFC3339Nano | 2006-01-02T15:04:05.999999999Z07:00 |
  | Kitchen     | 3:04PM                              |
  | Stamp       | Jan _2 15:04:05                     |
  | StampMilli  | Jan _2 15:04:05.000                 |
  | StampMicro  | Jan _2 15:04:05.000000              |
  | StampNano   | Jan _2 15:04:05.000000000           |
  | DateTime    | 2006-01-02 15:04:05                 |
  | DateOnly    | 2006-01-02                          |
  | TimeOnly    | 15:04:05                            |
  | Seconds1900 | 3952507337                          |
  | Seconds1970 | 1743518576                          |

  </details>
