# Proton Drive WebDAV Bridge

An alternative approach to interacting with Proton Drive on Linux.

Instead of integrating the Proton Drive API into a tool like rclone, this program will expose a local interface to it
using the WebDAV protocol. This approach has two benefits over the existing backend in rclone:

* You can use any client that supports WebDAV (even though rclone is still the recommended one, for reasons explained
  below)
* Because the bridge runs as a daemon, it can cache the entire directory tree and only update it when required, through
  the event system of Proton Drive. This enables near instant directory listings through WebDAV.

## Installation

To build this software, you need a recent Go toolchain (1.18 or newer) and git.

```bash
$ go version
go version go1.21.5 linux/amd64
$ git --version
git version 2.43.0
```

Building and installing the bridge is mostly automatic. You can use the GOBIN environment variable to select the target
directory. I recommend using `$HOME/.local/bin` because it is fairly standard and most likely already part of the PATH
variable in your shell. If you use a different target, you must make sure it is part of PATH.

```bash
$ git clone https://github.com/StollD/proton-webdav-bridge
$ cd proton-webdav-bridge
$ env GOBIN="$HOME/.local/bin" go install .
```

## Login

Before using the bridge, you need to provide it with your credentials and let it generate a login token. Only the token
will be stored, not your username or password.

Run the following command and follow the instructions it prints to your terminal.

```bash
$ proton-webdav-bridge --login
```

## Running the bridge

Running the WebDAV bridge is as simple as running the program without any arguments.

```bash
$ proton-webdav-bridge
```

By default, the WebDAV server will listen on http://127.0.0.1:7984, but you can change this with the `--addr` option.

Depending on the amount (not the size!) of files and directories in your drive, the startup might take quite a while,
because the bridge is caching the metadata of all objects, to speed up WebDAV lookups.

For starting the bridge automatically when you log in, I recommend using a systemd user service. A basic service file
that you can use is in the `systemd` directory of this repository.

Download it into `$HOME/.config/systemd/user`, and create the directory if it doesn't exist. Then enable it like this.

```bash
$ systemctl --user daemon-reload
$ systemctl --user enable --now proton-webdav-bridge.service
```

Keep in mind that this service will only work if you installed the bridge to `$HOME/.local/bin`. If you changed the
path, you need to adjust the service file as well.

## WebDAV, Clients and Rclone

The WebDAV standard does not include support for fetching file hashes, which makes it less suitable for a two-way sync,
because the client has to guess if two files are identical based on the modification time.

There is also no standard way for the client to set the modification time, which means that uploading existing files
resets their modification time.

To bridge this gap, the bridge implements a few extensions:
* Modification time can be set through a `X-OC-Mtime` header (OwnCloud extension)
* The SHA1 of a file can be read through the `checksums` property (OwnCloud extension)
* The SHA1 of a file can be read through the `sha1hex` property (FastMail extension)

To make full use of these extensions, you should use the rclone WebDAV backend and configure its vendor type to
`fastmail`. This maps pretty much 1:1 to what Proton Drive and this bridge support and will provide the best
experience.

If there are any other clients that support these extensions, or if there are useful extensions I missed, please open
an issue or send a pull request!

## Thanks

* henrybear327 for publishing https://github.com/henrybear327/Proton-API-Bridge
* Proton for publishing https://github.com/ProtonMail/go-proton-api
