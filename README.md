# go-peerdisco

轻量 UDP gossip 成员发现库 + `peerd` 管理服务：Join/Leave、Ping/Ack、间接探测、Suspicion→Dead、元数据同步。

## Build / Test

```bash
go build ./...
go test ./... -count=1
```

## Daemon

```bash
go run ./cmd/peerd -addr :8105 -web web
```

打开 http://127.0.0.1:8105/ 查看成员表与怀疑队列。
