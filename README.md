# k8s client with crd

For use k8s client easily

## usage:

Add client to go mod

```bash
go get github.com/changqings/k8scrd/client

```

New k8s client

```go
client, err := client.NewClient()
if err != nil {
    panic(err)
}
```
