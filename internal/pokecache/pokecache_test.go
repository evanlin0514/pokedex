package pokecache

import (
	"testing"
	"time"
	"fmt"
)


func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://pokeapi.co/api/v2/location/",
			val: []byte(`{"count":1039,"next":"https://pokeapi.co/api/v2/location/?offset=20&limit=20","previous":null,"results":[{"name":"canalave-city","url":"https://pokeapi.co/api/v2/location/1/"},{"name":"eterna-city","url":"https://pokeapi.co/api/v2/location/2/"},{"name":"pastoria-city","url":"https://pokeapi.co/api/v2/location/3/"},{"name":"sunyshore-city","url":"https://pokeapi.co/api/v2/location/4/"},{"name":"sinnoh-pokemon-league","url":"https://pokeapi.co/api/v2/location/5/"},{"name":"oreburgh-mine","url":"https://pokeapi.co/api/v2/location/6/"},{"name":"valley-windworks","url":"https://pokeapi.co/api/v2/location/7/"},{"name":"eterna-forest","url":"https://pokeapi.co/api/v2/location/8/"},{"name":"fuego-ironworks","url":"https://pokeapi.co/api/v2/location/9/"},{"name":"mt-coronet","url":"https://pokeapi.co/api/v2/location/10/"},{"name":"great-marsh","url":"https://pokeapi.co/api/v2/location/11/"},{"name":"solaceon-ruins","url":"https://pokeapi.co/api/v2/location/12/"},{"name":"sinnoh-victory-road","url":"https://pokeapi.co/api/v2/location/13/"},{"name":"ravaged-path","url":"https://pokeapi.co/api/v2/location/14/"},{"name":"oreburgh-gate","url":"https://pokeapi.co/api/v2/location/15/"},{"name":"stark-mountain","url":"https://pokeapi.co/api/v2/location/16/"},{"name":"spring-path","url":"https://pokeapi.co/api/v2/location/17/"},{"name":"turnback-cave","url":"https://pokeapi.co/api/v2/location/18/"},{"name":"snowpoint-temple","url":"https://pokeapi.co/api/v2/location/19/"},{"name":"wayward-cave","url":"https://pokeapi.co/api/v2/location/20/"}]}`),
		},
		{
			key: "https://pokeapi.co/api/v2/berry-firmness/2/",
			val: []byte(`{"firmness":{"name":"soft","url":"https://pokeapi.co/api/v2/berry-firmness/2/"},"flavors":[{"flavor":{"name":"spicy","url":"https://pokeapi.co/api/v2/berry-flavor/1/"},"potency":10},{"flavor":{"name":"dry","url":"https://pokeapi.co/api/v2/berry-flavor/2/"},"potency":0},{"flavor":{"name":"sweet","url":"https://pokeapi.co/api/v2/berry-flavor/3/"},"potency":0},{"flavor":{"name":"bitter","url":"https://pokeapi.co/api/v2/berry-flavor/4/"},"potency":0},{"flavor":{"name":"sour","url":"https://pokeapi.co/api/v2/berry-flavor/5/"},"potency":0}],"growth_time":3,"id":1,"item":{"name":"cheri-berry","url":"https://pokeapi.co/api/v2/item/126/"},"max_harvest":5,"name":"cheri","natural_gift_power":60,"natural_gift_type":{"name":"fire","url":"https://pokeapi.co/api/v2/type/10/"},"size":20,"smoothness":25,"soil_dryness":15}`),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
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
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("https://pokeapi.co/api/v2/berry-firmness/2/", []byte(`{"firmness":{"name":"soft","url":"https://pokeapi.co/api/v2/berry-firmness/2/"},"flavors":[{"flavor":{"name":"spicy","url":"https://pokeapi.co/api/v2/berry-flavor/1/"},"potency":10},{"flavor":{"name":"dry","url":"https://pokeapi.co/api/v2/berry-flavor/2/"},"potency":0},{"flavor":{"name":"sweet","url":"https://pokeapi.co/api/v2/berry-flavor/3/"},"potency":0},{"flavor":{"name":"bitter","url":"https://pokeapi.co/api/v2/berry-flavor/4/"},"potency":0},{"flavor":{"name":"sour","url":"https://pokeapi.co/api/v2/berry-flavor/5/"},"potency":0}],"growth_time":3,"id":1,"item":{"name":"cheri-berry","url":"https://pokeapi.co/api/v2/item/126/"},"max_harvest":5,"name":"cheri","natural_gift_power":60,"natural_gift_type":{"name":"fire","url":"https://pokeapi.co/api/v2/type/10/"},"size":20,"smoothness":25,"soil_dryness":15}`))

	_, ok := cache.Get("https://pokeapi.co/api/v2/berry-firmness/2/")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://pokeapi.co/api/v2/berry-firmness/2/")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}