package misc

func NewCache(size int) *Cache {
	cache := &Cache{}
	cache.Init(size)
	return cache
}

type Cache struct {
	heap  []Element
	index int
	size  int
}

func (cache *Cache) Init(size int) {
	if cache == nil {
		return
	}

	cache.size = size
}

func (cache *Cache) Alloc() *Element {
	if cache == nil {
		return &Element{}
	}

	if cache.index >= len(cache.heap) {
		if cache.size <= 0 {
			return &Element{}
		}

		cache.index = 0
		cache.heap = make([]Element, cache.size)
	}

	e := &cache.heap[cache.index]
	cache.index++

	return e
}
