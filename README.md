# FSSCANNER

This package allows file system directory recursive walks and files listing.
It uses https://github.com/spf13/afero allowing the use of any type of storage as
long as an adapter for the file system abstraction interface exists.

Here is a pseudo code on how to use this package:

```
// Select a file system interface
fs, err := fsscanner.NewFS(fsscanner.MEM)

// Request recursive-scan into a directory in search for files
scanner, err := fs.Scan("this-directory");
if err != nil {
    panic(err)
}

//List file results from the scan
files, err := scanner.List()
if err != nil {
    panic(err)
}

//Do something with files found
for _, file := range files {
    ...
}
```
