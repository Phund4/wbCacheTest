package internal

// Задание: "Написать сервис (класс/структура) кэширования".

// Основные условия:
// - Кэш ограниченной емкости, метод вытеснения ключей LRU;
// - Сервис должен быть потокобезопасный;
// - Сервис должен принимать любые значения;
// - Реализовать unit-тесты.

// Сервис должен реализовывать следующий интерфейс:
//     type ICache interface {
//         Cap() int
//         Len() int
//         Clear() // удаляет все ключи
//         Add(key, value interface{})
//         AddWithTTL(key, value interface{}, ttl time.Duration) // добавляет ключ со сроком жизни ttl
//         Get(key interface{}) (value interface{}, ok bool)
//         Remove(key interface{})
//     }

import (
	"container/list"
	"sync"
	"time"
	"fmt"
)

type ICache interface {
	Cap() int
	Len() int
	Clear() // удаляет все ключи
	Add(key, value interface{})
	AddWithTTL(key, value interface{}, ttl time.Duration) // добавляет ключ со сроком жизни ttl
	Get(key interface{}) (value interface{}, ok bool)
	Remove(key interface{})
}

type node struct {
	key   interface{}
	value interface{}
	ttl   time.Time
}

type lruCache struct {
	capacity int
	items    map[interface{}]*list.Element
	order    *list.List
	mutex    sync.Mutex
}

// Функция инициализации LRU кэша
func NewLRUCache(capacity int) ICache {
	return &lruCache{
		capacity: capacity,
		items:    make(map[interface{}]*list.Element),
		order:    list.New(),
		mutex:    sync.Mutex{},
	}
}

func (c *lruCache) Cap() int {
	return c.capacity
}

func (c *lruCache) Len() int {
	// блокируем доступ
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.order.Len()
}

func (c *lruCache) Clear() {
	// блокируем доступ
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// обнуляем хранилище и очередь из элементов
	c.items = make(map[interface{}]*list.Element)
	c.order.Init()
}

func (c *lruCache) Add(key, value interface{}) {
	c.AddWithTTL(key, value, 0)
}

func (c *lruCache) AddWithTTL(key, value interface{}, ttl time.Duration) {
	// блокируем доступ
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// элемент уже найден
	if elem, ok := c.items[key]; ok {
		// обновляем его в очереди
		c.order.MoveToFront(elem)
		// обновляем значение и ttl
		elem.Value.(*node).value = value
		elem.Value.(*node).ttl = time.Now().Add(ttl)
		return
	}

	// если кэш переполнен
	if c.order.Len() >= c.capacity {
		back := c.order.Back()
		// проходимся по элементам, пока не найдем тот, у кого ttl expiry
		for back != nil {
			if time.Now().After(back.Value.(*node).ttl) {
				// Удаляем его
				c.order.Remove(back)
				delete(c.items, back.Value.(*node).key)
				break;
			}
			back = back.Prev();
		}
		// если в кэше все элементы нужны
		if back == nil {
			fmt.Println("error: cache is overflowing. Try later!");
			return 
		}
	}

	// создаем новый элемент
	item := &node{
		key:   key,
		value: value,
		ttl:   time.Now().Add(ttl),
	}
	elem := c.order.PushFront(item)
	c.items[key] = elem
}

func (c *lruCache) Get(key interface{}) (interface{}, bool) {
	// блокируем доступ
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// ищем нужный элемент по ключу
	if elem, ok := c.items[key]; ok {
		// если он есть, но ttl expiry, то удаляем его
		if time.Now().After(elem.Value.(*node).ttl) {
			c.order.Remove(elem)
			delete(c.items, key)
			return nil, false
		}

		// обновляем в очереди
		c.order.MoveToFront(elem)
		return elem.Value.(*node).value, true
	}
	return nil, false
}

func (c *lruCache) Remove(key interface{}) {
	// блокируем доступ
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// ищем нужный элемент по ключу
	if elem, ok := c.items[key]; ok {
		c.order.Remove(elem)
		delete(c.items, key)
	}
}

// Функция вывода всего кэша в консоль
func PrintCache(cache ICache) {
	for key, value := range cache.(*lruCache).items {
		fmt.Printf("key: %v, value: %v\n", key, value.Value.(*node).value);
	}
	fmt.Println()
}