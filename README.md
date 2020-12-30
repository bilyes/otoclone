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

1. Make sure `inotify-tools` is installed

2. Install `rclone`. To verify if it's installed, run `rclone --version`

3. Configure your remotes (backup destinations) in rclone: `rclone config`. For
   more details on how to install and configure rclone check out [their website](https://rclone.org/).

4. Create configuration file for otoclone called `config.yml` in
   `$XDG_CONFIG_HOME/otoclone`. If you don't have `$XDG_CONFIG_HOME`, put the
   configuration file in `$HOME/.config/otoclone`.

   ```
   config.yml
   ----------
   
   folders:
     documents:
       path: /home/jim/documents
       strategy: copy
       remotes:
         - Dropbox
         - OneDrive
         - S3
       ignoreList:
         - ignore-me.txt
   ```

## Keep Alive
WIP
