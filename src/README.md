# golocal - cross platform protocol handler

## use

Running it without any arguments starts it in GUI mode, giving options to install or remove the protocol handler.

It registers a handler for `golocal://` - everything after is passed to `xdg-open` on Linux and `cmd \C start` on Windows.

## build

Additional build dependencies on ArchLinux (next to go, gcc, make, etc)

  * mingw-w64-gcc

For building a `make all` should be enough.
