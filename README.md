# ftp2s3

An FTP to S3 adapter written in Go. Access an S3 bucket and its contents via an FTP interface.

ftp2s3 is based on the [Graval][] FTP server framework. It is very much a work in progress, but major functions mostly work. You can list directory contents, upload files, download files, and delete files. All these functions could use more testing and thought since it has been built rather hastily…

[graval]: https://github.com/koofr/graval

## Usage

```sh
ftp2s3 -aws-region="eu-west-1" \
-aws-bucket-name="my-s3-bucket" \
-aws-access-key-id="AWS_KEY" \
-aws-secret-access-key="AWS_SECRET" \
-ftp-username="svimes" \
-ftp-username="swordfish" \
-ftp-server-name="PseudopolisYardFTPd"
```

Configuration via environment variables is also possible (e.g. `-ftp-server-name="…"` → `FTP_SERVER_NAME="…"`).

Alternatively you can supply a simple configuration file:

```
# Put this into myftp.conf
aws-region eu-west-1
aws-bucket-name my-s3-bucket
aws-access-key-id AWS_KEY
aws-secret-access-key AWS_SECRET
ftp-username svimes
ftp-username swordfish
ftp-server-name PseudopolisYardFTPd
```

```sh
ftp2s3 -config="myftp.conf"
```

## Why would I want this?

If you have legacy software/hardware that doesn't support S3 directly (for example some web cameras).

## What's implemented?

Some commands (such as those that deal with mode switching) are implicitly supported through Graval and other features/commands are explicitly supported (see below), though it is possible that some bugs persist.

### Explicitly supported FTP commands

* `USER` (login as the given user)
* `PASS` (authenticate with the given password)
* `LIST` (directory listing)
* `DELE` (delete file)
* `MDTM` (last-modified time of a specified file)
* `SIZE` (size of a file in bytes)
* `STOR` (upload file)
* `RETR` (download file)

### Kinda supported

* `CWD` (change directory)
    * Changing “directories” is kind of implemented, but all the implications haven't been fully thought out yet (since S3 doesn't really have directories). There might be some unexpected bugs.

### Explicitly unsupported FTP commands

S3 doesn't have directories in the conventional sense of the concept, so these commands have been left unimplemented for now.

* `MKD` (make directory)
* `RMD` (remove directory)

As S3 doesn't support renames as such, so `RNFR`/`RNTO` (rename from/rename to) have been left unimplemented.

If you come up with a good way of supporting these features, I'll be more than happy to accept PRs.

## I fixed a thing! / I implemented a new feature! / I wrote some unit tests! / This Go code isn't idiomatic!

Great! Please submit a pull request and I'll be more than happy to merge it (as long as it's sane)…

## License and Copyright

Copyright © 2015 Matias Korhonen.

Licensed under the MIT License, see the [LICENSE.txt](LICENSE.txt) file for details.
