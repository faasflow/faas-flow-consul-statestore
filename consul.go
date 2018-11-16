package ConsulStateStore

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
	faasflow "github.com/s8sg/faasflow"
	"strconv"
)

type ConsulStateStore struct {
	consulKeyPath string
	consulClient  *consul.Client
	kv            *consul.KV

	RetryCount int
}

func GetConsulStateStore() (faasflow.StateStore, error) {

	consulST := &ConsulStateStore{}

	cc, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		return nil, err
	}
	consulST.consulClient = cc
	consulST.kv = cc.KV()

	consulST.RetryCount = 10

	return consulST, nil
}

func (consulStore *ConsulStateStore) Init(flowName string, requestId string) error {
	consulStore.consulKeyPath = fmt.Sprintf("faasflow/%s/%s", flowName, requestId)
	return nil
}

func (consulStore *ConsulStateStore) Create(vertexs []string) error {

	for _, vertex := range vertexs {
		key := fmt.Sprintf("%s/%s", consulStore.consulKeyPath, vertex)
		p := &consul.KVPair{Key: key, Value: []byte("0")}
		_, err := consulStore.kv.Put(p, nil)
		if err != nil {
			return fmt.Errorf("failed to create vertex %s, error %v", vertex, err)
		}
	}
	return nil
}

func (consulStore *ConsulStateStore) IncrementCounter(vertex string) (int, error) {
	count := 0
	key := fmt.Sprintf("%s/%s", consulStore.consulKeyPath, vertex)
	for i := 0; i < consulStore.RetryCount; i++ {
		pair, _, err := consulStore.kv.Get(key, nil)
		if err != nil {
			return 0, fmt.Errorf("failed to get vertex %s, error %v", vertex, err)
		}
		modifyIndex := pair.ModifyIndex
		counter, err := strconv.Atoi(string(pair.Value))
		if err != nil {
			return 0, fmt.Errorf("failed to convert counter for %s, error %v", vertex, err)
		}

		count := counter + 1
		counterStr := fmt.Sprintf("%d", count)

		p := &consul.KVPair{Key: key, Value: []byte(counterStr), ModifyIndex: modifyIndex}
		_, err = consulStore.kv.Put(p, nil)
		if err != nil {
			continue
		}
	}
	return count, nil
}

func (consulStore *ConsulStateStore) Cleanup() error {
	_, err := consulStore.kv.DeleteTree(consulStore.consulKeyPath, nil)
	if err != nil {
		return fmt.Errorf("error removing %s, error %v", consulStore.consulKeyPath, err)
	}
	return nil
}
