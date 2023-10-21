package cache

import (
	"errors"
	"log"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"go_app/pkg/order"
)

func ParseCleanupExpiration() (time.Duration, time.Duration) {

	envExpiration := os.Getenv("DEFAULT_EXPIRATION")
	envCleanupInterval := os.Getenv("CLEANUP_INTERVAL")
	defaultExp, err := strconv.Atoi(envExpiration)
	if err != nil {
		log.Fatal("env variable DEFAULT_EXPIRATION must be integer")
	}
	cleanup, err := strconv.Atoi(envCleanupInterval)
	if err != nil {
		log.Fatal("env variable CLEANUP_INTERVAL must be integer")
	}
	defaultExpiration := time.Duration(defaultExp * int(time.Minute))
	cleanupInterval := time.Duration(cleanup * int(time.Minute))
	return defaultExpiration, cleanupInterval
}

func SortOrders(orders []*order.Order) {
	sort.SliceStable(orders, func(i, j int) bool { return orders[i].DateCreated > orders[j].DateCreated })
}

type Cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	items             map[string]Item
}

type Item struct {
	Value      interface{}
	Created    time.Time
	Expiration int64
}

func New(defaultExpiration, cleanupInterval time.Duration) *Cache {

	// инициализируем карту(map) в паре ключ(string)/значение(Item)
	items := make(map[string]Item)

	cache := Cache{
		items:             items,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}

	// Если интервал очистки больше 0, запускаем GC (удаление устаревших элементов)
	if cleanupInterval > 0 {
		cache.StartGC() // данный метод рассматривается ниже
	}

	return &cache
}

func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	c.Lock()
	defer c.Unlock()
	var expiration int64

	// Если продолжительность жизни равна 0 - используется значение по-умолчанию
	if duration == 0 {
		duration = c.defaultExpiration
	}

	// Устанавливаем время истечения кеша
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}
	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}

}

func (c *Cache) SetAllOrders(value []*order.Order, duration time.Duration) {
	c.Lock()
	defer c.Unlock()
	var expiration int64

	// Если продолжительность жизни равна 0 - используется значение по-умолчанию
	if duration == 0 {
		duration = c.defaultExpiration
	}

	// Устанавливаем время истечения кеша
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}
	created := time.Now()
	for _, order := range value {
		c.items[order.OrderID] = Item{
			Value:      order,
			Expiration: expiration,
			Created:    created,
		}
	}
}

func (c *Cache) Get(key string) (*order.Order, bool) {
	c.RLock()
	defer c.RUnlock()
	item, found := c.items[key]
	// ключ не найден
	if !found {
		return nil, false
	}
	// Проверка на установку времени истечения, в противном случае он бессрочный
	if item.Expiration > 0 {
		// Если в момент запроса кеш устарел возвращаем nil
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}
	}
	order := item.Value.(*order.Order)
	return order, true
}

func (c *Cache) GetAllOrders() ([]*order.Order, bool) {
	var orders []*order.Order
	c.RLock()
	defer c.RUnlock()
	for _, item := range c.items {
		order, ok := item.Value.(*order.Order)
		if !ok {
			log.Println("Not order type")
			continue
		}
		if item.Expiration > 0 {
			// Если в момент запроса кеш устарел, не добавляем его в слайс
			if time.Now().UnixNano() > item.Expiration {
				continue
			}
		}
		orders = append(orders, order)
	}
	if len(orders) == 0 {
		return nil, false
	}
	SortOrders(orders)
	return orders, true
}

func (c *Cache) Delete(key string) error {

	c.Lock()

	defer c.Unlock()

	if _, found := c.items[key]; !found {
		return errors.New("key not found")
	}

	delete(c.items, key)

	return nil
}

func (c *Cache) StartGC() {
	go c.GC()
}

func (c *Cache) GC() {

	for {
		// ожидаем время установленное в cleanupInterval
		<-time.After(c.cleanupInterval)

		if c.items == nil {
			return
		}

		// Ищем элементы с истекшим временем жизни и удаляем из хранилища
		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearItems(keys)

		}

	}

}

// expiredKeys возвращает список "просроченных" ключей
func (c *Cache) expiredKeys() (keys []string) {

	c.RLock()

	defer c.RUnlock()

	for k, i := range c.items {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}

	return
}

// clearItems удаляет ключи из переданного списка, в нашем случае "просроченные"
func (c *Cache) clearItems(keys []string) {

	c.Lock()

	defer c.Unlock()

	for _, k := range keys {
		delete(c.items, k)
	}
}
