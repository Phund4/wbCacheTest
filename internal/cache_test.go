package internal_test

import (
	"testing"
	"time"

	"github.com/Phund4/wbCacheTest/internal"
)

func TestLRUCache(t *testing.T) {
	lru := internal.NewLRUCache(4);
	lru.Add("a", "aaa")
	if _, ok := lru.Get("a"); ok {
		t.Errorf("expected nil, but element was found")
	}
	
	lru.AddWithTTL("a", "aaa", time.Second * 3);
	lru.Clear()
	if lru.Len() != 0 {
		t.Errorf("after clean() function len(cache) != 0");
	}

	lru.AddWithTTL("c", "ccc", time.Second * 5)
	lru.AddWithTTL("d", "ddd", time.Second * 1);
	lru.Add("a", "aaa");
	lru.Add("b", "bbb");

	if _, ok := lru.Get("c"); !ok {
		t.Errorf("expected ccc, but element wasn't found")
	}

	lru.Add("g", "ggg");

	if _, ok := lru.Get("a"); ok {
		t.Errorf("expected nil, but element was found");
	}

	time.Sleep(time.Second * 2);
	lru.Add("f", "fff");

	if _, ok := lru.Get("d"); ok {
		t.Errorf("expected nil, but element was found");
	}

	lru2 := internal.NewLRUCache(2);
	lru2.AddWithTTL("a", "aaa", time.Second * 2);
	lru2.AddWithTTL("b", "bbb", time.Second * 2);
	lru2.Add("c", "ccc")

	if _, ok := lru2.Get("c"); ok {
		t.Errorf("expected nil, but element was found")
	}
}
