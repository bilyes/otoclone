# OtoClone

An automatic backup and sync tool

## Compatibility
Linux

## Requirements
- [rclone](https://github.com/rclone/rclone)
- inotify-tools

## Installation

### Build from source

#### Download the source code and build the binary

```
git clone https://github.com/ilyessbachiri/otoclone
cd otoclone/
go build -o otoclone
```

#### Install the binary

Add the binary to `/usr/bin`
```
sudo mv otoclone /usr/bin
```

If you don't have root privileges or you choose not to touch `/usr/bin`, you can
copy the binary into `~/.local/bin/`. If this is your case, make sure
`~/.local/bin` is part of your `$PATH` env variable.

## Configuration

- Make sure `inotify-tools` is installed
- Install `rclone`. To verify if it's installed, run `rclone --version`
- Configure your remotes (backup destinations) in rclone: `rclone config`
For more details on how to install and configure rclone check out [their website](https://rclone.org/).

## Keep Alive
WIP
