package main

import (
	//"fmt"

	"github.com/hashicorp/consul/api"
)

func watchKey(key string, index uint64) (string, uint64) {
	// Get a new client
	c := api.DefaultConfig()
	c.Address = "consul:8500"
	client, err := api.NewClient(c)
	if err != nil {
		panic(err)
	}

	// Get a handle to the KV API
	kv := client.KV()

	q := &api.QueryOptions{WaitIndex: index}
	// GET a new KV pair
	kvp, meta, err := kv.Get(key, q)
	if err != nil {
		panic(err)
	}
	if kvp != nil {
		return string(kvp.Value), meta.LastIndex
	} else {
		// handle the case where the key doesn't exist
		return "false", 0
	}
}
