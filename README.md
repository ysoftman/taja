# taja 게임

터미널 환경에서 떨어지는 단어를 빠르게 타이핑하여 제거하는 게임

## Enviroment

- Windows : cmd.exe
- Linux : terminal
- Mac : terminal(iterm2)

## Build & Run

```bash
# build
go get "github.com/mattn/go-runewidth"
go get "github.com/nsf/termbox-go"
go build

# run
./taja
```

- 현 위치에 패키지 다운로드(필요시)

  ```bash
  echo $GOPATH
  export GOPATH_BACKUP=$GOPATH
  export GOPATH=$PWD
  go get github.com/nsf/termbox-go
  export GOPATH=$GOPATH_BACKUP
  unset GOPATH_BACKUP && echo $GOPATH_BACKUP
  ```
