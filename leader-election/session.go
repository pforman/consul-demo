package main

import (
	//"fmt"

	"github.com/hashicorp/consul/api"
)

func consulSession(name string, ttl string) string {
	// Get a new client
	c := api.DefaultConfig()
	c.Address = "consul:8500"
	client, err := api.NewClient(c)
	if err != nil {
		panic(err)
	}

	s := client.Session()
	se := &api.SessionEntry{Name: name, TTL: ttl}

	id, _, err := s.Create(se, nil)
	if err != nil {
		panic(err)
	}
	return id
}

func acquireLock(sid string, key string, value string) bool {
	// Get a new client
	c := api.DefaultConfig()
	c.Address = "consul:8500"
	client, err := api.NewClient(c)
	if err != nil {
		panic(err)
	}

	// Get a handle to the KV API
	kv := client.KV()

	kvp := &api.KVPair{
		Key:     key,
		Value:   []byte(value),
		Session: sid,
	}
	// GET a new KV pair
	leader, _, err := kv.Acquire(kvp, nil)
	if err != nil {
		panic(err)
	}

	return leader
}

func releaseLock(sid string, key string) bool {
	// Get a new client
	c := api.DefaultConfig()
	c.Address = "consul:8500"
	client, err := api.NewClient(c)
	if err != nil {
		panic(err)
	}

	// Get a handle to the KV API
	kv := client.KV()

	kvp := &api.KVPair{
		Key:     key,
		Session: sid,
	}
	// GET a new KV pair
	ok, _, err := kv.Acquire(kvp, nil)
	if err != nil {
		panic(err)
	}

	return ok
}
