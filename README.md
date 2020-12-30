# OtoClone

An automatic backup and sync tool

## Compatibility
Linux

## Requirements
- [rclone](https://github.com/rclone/rclone)
- inotify-tools

## Installation

```
go build -o otoclone
sudo mv otoclone /usr/bin
```

If you don't have root privileges, you can copy the binary into `~/.local/bin/`.
If this is your case, make sure `~/.local/bin` is part of your `$PATH` env
variable.

## Configuration
WIP

## Keep Alive
WIP
