package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

func registerService(service, httpAddr string) {
	// Get a new client
	c := api.DefaultConfig()
	c.Address = "consul:8500"
	client, err := api.NewClient(c)
	if err != nil {
		panic(err)
	}

	// Get a handle to the KV API
	kv := client.KV()

	key := fmt.Sprintf("%s/%s", service, httpAddr)

	// PUT a new KV pair
	//p := &api.KVPair{Key: key, Value: []byte('stuff')}
	p := &api.KVPair{Key: key}
	_, err = kv.Put(p, nil)
	if err != nil {
		panic(err)
	}
}

func deRegisterService(service, httpAddr string) {
	// Get a new client
	c := api.DefaultConfig()
	c.Address = "consul:8500"
	client, err := api.NewClient(c)
	if err != nil {
		panic(err)
	}

	// Get a handle to the KV API
	kv := client.KV()

	key := fmt.Sprintf("%s/%s", service, httpAddr)

	_, err = kv.Delete(key, nil)
	if err != nil {
		panic(err)
	}
}
