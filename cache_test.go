package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/mecebeci/pokedexcli/internal/pokecache"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := pokecache.NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const interval = 50 * time.Millisecond  
	cache := pokecache.NewCache(interval) 

	key := "https://example.com"
	val := []byte("testdata")

	cache.Add(key, val)

	got, ok := cache.Get(key)
	if !ok || string(got) != string(val) {
		t.Errorf("expected to find key immediately after Add")
	}

	
	time.Sleep(2 * interval)

	_, ok = cache.Get(key)
	if ok {
		t.Errorf("expected key to be removed by reapLoop")
	}
}
