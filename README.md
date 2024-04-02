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
