# faasflow-consul-statestore
A **[faas-flow](https://github.com/s8sg/faas-flow)** statestore implementation that uses Consul to store state

## Consul State-Store configuration
Below are the configuration 
```
consul_url
consul_dc
```

### Use Consul StateStore in `faasflow`
* Set the `stack.yml` with the necessary environments
```yaml
    consul_url: "consul.faasflow:8500"
    consul_dc: "dc1"
```
* Use the `ConsulStateStore` as a DataStore on `handler.go`
```go
consulStateStore "github.com/s8sg/faas-flow-consul-statestore"

func DefineStateStore() (faasflow.StateStore, error) {
        consulss, err := consulStateStore.GetConsulStateStore(os.Getenv("consul_url"), os.Getenv("consul_dc"))
        if err != nil {
                return nil, err
        }
        return consulss, nil
}
```
