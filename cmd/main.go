package main

import (
	"time"

	"github.com/Phund4/wbCacheTest/internal"

)

func main() {
	lru := internal.NewLRUCache(4);
	lru.AddWithTTL("c", "ccc", time.Second * 5)
	lru.AddWithTTL("d", "ddd", time.Second * 1);
	lru.Add("a", "aaa");
	lru.Add("b", "bbb");

	internal.PrintCache(lru)

	lru.Add("g", "ggg");
	internal.PrintCache(lru)

	time.Sleep(time.Second * 3);

	lru.Add("f", "fff");
	internal.PrintCache(lru)

	lru2 := internal.NewLRUCache(2);
	lru2.AddWithTTL("a", "aaa", time.Second * 2);
	lru2.AddWithTTL("b", "bbb", time.Second * 2);

	internal.PrintCache(lru2);

	lru2.Add("c", "ccc")

	internal.PrintCache(lru2);
}