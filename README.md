# clipboard-monitor

Two way sync between the clipboard and a file.

If the system clipboard is updated (e.g. by pressing `ctrl+c`), the changes are written to `clipboard.txt`.
If the `clipboard.txt` file is updated (e.g. by `echo 123 > clipboard.txt`) this change is written to the system clipboard.

## Why?

I am using this script to have a shared clipboard between host and a virtual machine. Both machines need to have access to the `clipboard.txt` file. This can for example be accomplished by samba or SMB or any other file sync software.

Build on Linux for Linux:

``` bash
go build

```

Build on Linux for Windows:

```bash
GOOS=windows GOARCH=amd64 go build
```

## How to use?

### Windows Host, Linux Guest

1. Configure and enable samba server in Linux
2. Mount samba share in Windows
3. Copy both the Linux and Windows binary on the share (TODO this is not very secure, because now the VM can modify the binary run on the host)
   1. On the Linux machine run the Linux binary
   2. On the Windows machine run the Windows binary
4. Clipboard sharing should now work

## Limitations

* This is most probably not very secure and should be improved
* notify doesn't work from Linux to Windows (only the other way around). Thus there is a fallback to poll the file very 5 seconds
