# OtoClone

An automatic backup and sync utility that watches directories. It listens to
filesystem events on these directories and copies their content to remote
destination. These destinations can be a folder on the local filesystem, a
network drive mounted on the filesystem, or a bucket on a cloud storage
provider. See the list of supported providers [here](https://rclone.org/#providers). 

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

Add the binary to `/usr/local/bin`
```
sudo mv otoclone /usr/local/bin
```

If you don't have root privileges or you choose not to touch `/usr/local/bin`, you can
copy the binary into `~/.local/bin/`. If this is your case, make sure
`~/.local/bin` is part of your `$PATH` env variable.

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
       strategy: copy
       remotes:
         - name: Dropbox
           bucket: documents
         - name: OneDrive
           bucket: docs
         - name: S3
           bucket: Documents
       ignoreList:
         - ignore-me.txt

     photos:
       path: /home/jim/photos
       strategy: copy
       remotes:
         - GoogleDrive
   ```

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

