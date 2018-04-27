# taja 연습

- 현 위치에 패키지 추가(필요시)

```bash
echo $GOPATH
export GOPATH_BACKUP=$GOPATH
export GOPATH=$PWD
go get github.com/nsf/termbox-go
export GOPATH=$GOPATH_BACKUP
unset GOPATH_BACKUP && echo $GOPATH_BACKUP
```
