# Building


### Prerequisites

[Go](https://go.dev/doc/install) | [MSYS2](https://www.msys2.org/)

Install packages
```console
C:\msys64\msys2_shell.cmd -defterm -use-full-path -no-start -mingw64 -here -c "pacman --noconfirm --needed -S $MINGW_PACKAGE_PREFIX-gcc"
```

Build
```console
C:\msys64\msys2_shell.cmd -defterm -use-full-path -no-start -mingw64 -here -c "go build ."
```

Release
```console
C:\msys64\msys2_shell.cmd -defterm -use-full-path -no-start -mingw64 -here
rev=$(git rev-parse --short HEAD)
go build -tags release -ldflags "-H=windowsgui -X main.rev=$rev"
```