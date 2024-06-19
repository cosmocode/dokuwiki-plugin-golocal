# golocal - cross platform protocl handler

## use

Running it without any arguments starts it in GUI mode, giving options to install or remove the protocol handler. It registers a handler for `locallink://` - everything after is passed to `xdg-open` on Linux and `cmd \C open` on Windows. 

## build

Additional build dependencies on ArchLinux (next to go, gcc, make, etc)

  * mingw-w64-gcc
    * mingw-w64-binutils-bin
    * mingw-w64-crt-bin
    * mingw-w64-headers-bin
    * mingw-w64-winpthreads-bin

For building a `make all` should be enough
