# OtoClone

An automatic backup and sync utility that watches directories. It listens to
filesystem events on these directories and copies their content to remote
destinations. A remote destination can be a folder on the local filesystem, a
network drive mounted on the filesystem, or a bucket on a cloud storage
provider. See the list of supported providers [here](https://rclone.org/#providers). 

## Compatibility

Linux

## Dependencies

- [rclone](https://github.com/rclone/rclone)
- inotify-tools

## Installation

The only repository supported for now is the [AUR](https://aur.archlinux.org/otoclone.git). To install `otoclone` on Arch
linux, use any AUR helper of your choice. Example with `yay`:
```
yay -S otoclone
```

## Configuration

1. Make sure `inotify-tools` is installed

2. Install `rclone`. To verify if it's installed, run `rclone --version`

3. Configure your remotes (backup destinations) in rclone by running `rclone config`. For
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
       strategy: sync
       remotes:
         - name: Dropbox
           bucket: documents
         - name: OneDrive
           bucket: docs
         - name: S3
           bucket: Documents
       ignoreList:
         - ignore-me.txt
       excludePattern: "*.jpg"

     photos:
       path: /home/jim/photos
       strategy: copy
       remotes:
         - name: GoogleDrive
           bucket: Pix
   ```

| Field            | Required? | Description                                                                                                                                                                                                                                                        | Example                                      |
|------------------|:---------:|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------|
| `path`           |    Yes    | The path of the folder to watch                                                                                                                                                                                                                                    | `home/jim/documents`                         |
| `strategy`       |    Yes    | The backup strategy to use. Currently there are 2 supported strategies: `copy` and `sync`                                                                                                                                                                          | `copy`                                       |
| `remotes`        |    Yes    | The remote destination to tranfer data to. Two fields are required for the remote:<br>- `name`: the name of the remote as defined in rclone. To list configured remotes, run `rclone listremotes`<br>- `bucket`: the path of the destination folder on the remote. | `name: Dropbox`<br><br>`bucket: backup/docs` |
| `ignoreList`     |    No     | A list of files whose filesystem events should<br> be ignored. These files **are not** ignored during backup however.                                                                                                                                              | `file1.lock`                                 |
| `excludePattern` |    No     | The pattern to exclude based on file globs as used by the unix shell. For more details see `rclone`'s [documentation](https://rclone.org/filtering/).                                                                                                              | `*.jpg`                                      |


## Build from source

#### Download the source code and build the binary

```
git clone https://github.com/ilyessbachiri/otoclone
cd otoclone/
go build -o otoclone
```

#### Install the binary

Add the binary to `/usr/local/bin`
```
sudo mv otoclone /usr/local/bin
```

If you don't have root privileges or you choose not to touch `/usr/local/bin`, you can
copy the binary into `~/.local/bin/`. If this is your case, make sure
`~/.local/bin` is part of your `$PATH` env variable.

## Keep Alive

`systemd` can be used on Linux to run `otoclone` as a service on boot and
restart it if it ever fails. Here's an example of how to accomplish this.

1. Add a systemd service file to manage `otoclone` under `/etc/systemd/system/`.
   You can name it `otoclone.service` for example.
   ```
   otoclone.service
   ----------------
   
   [Unit]
   Description=Automatic backup
   After=network.target

   [Service]
   Type=simple
   # Another Type: forking
   User=jim
   WorkingDirectory=/home/jim
   ExecStart=otoclone
   Restart=on-failure
   # Other restart options: always, on-abort, etc

   # The install section is needed to use
   # `systemctl enable` to start on boot
   # For a user service that you want to enable
   # and start automatically, use `default.target`
   # For system level services, use `multi-user.target`
   [Install]
   WantedBy=default.target
   ```
2. Enable the service using: `systemctl enable otoclone`
3. Start the service using: `systemctl start otoclone`

If everything is configured properly you should see the `otoclone` service
active when you run: `systemctl status otoclone`

Like any other systemd service the logs can be found in `journalclt` 
```
journalctl -u otoclone
```

