package cache

import (
	"container/list"
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
	"runtime"
	"sync"
	"time"

	sizing "github.com/DmitriyVTitov/size"
)

const (
	DefaultTTL     = 30              //time.Second
	DefaultminSize = 1 << 20         // 1Mb
	DefaultmaxSize = (1 << 20) * 128 // 128Mb
)

type KeyT int64
type ValueT State

type elem struct {
	key      KeyT
	val      ValueT
	deadTime int64
}

func (e *elem) size() uint32 {
	size := uint32((reflect.ValueOf(e.key).Type().Size())+16) * 2 // map index key size + map key size
	size += uint32(sizing.Of(e.val))                              // map value size
	size += 32                                                    // list.Element size
	size += uint32(reflect.ValueOf(e.deadTime).Type().Size())     // time.Duration size
	return size
}

type Cache struct {
	table         map[KeyT]*list.Element
	queue         *list.List
	ttl           uint32
	checkTTL      uint32
	curSize       uint32
	minSize       uint32
	maxSize       uint32
	isCollapse    chan bool
	isClose       chan bool
	closeCollapse sync.WaitGroup
	sync.RWMutex
}

// Creates a new Cache.
//  Input free or less argument: TTL, minBytesize, maxBytesize.
//  Example:
//  	c, _ := NewCache(10) set TTL to 10 sec
//		c, _ := NewCache(10, 512) set TTL to 10 sec, minByteSize to 512 byte
// 		c, _ := NewCache(10, 512, 1024) set TTL to 10 sec, minByteSize to 512 byte, maxByteSize to 1024 byte
/*const (
	DefaultTTL     = 30              //time.Second
	DefaultminSize = 1 << 20         // 1Mb
	DefaultmaxSize = (1 << 20) * 128 // 128Mb
)*/
func NewCache(param ...uint32) (*Cache, error) {
	count := len(param)
	flagForIndex := 0
	switch {
	case count == 1:
		if param[0] < 5 {
			return nil, errors.New("TTL is to low. It must be >5")
		}
		flagForIndex = 1
	case count == 2:
		if param[0] < 5 {
			return nil, errors.New("TTL is to low. It must be >5")
		}
		if param[1] > DefaultmaxSize {
			return nil, errors.New("minBytesize must be low then DefaultmaxSize")
		}
	case count == 3:
		if param[0] < 5 {
			return nil, errors.New("TTL is to low. It must be >5")
		}
		if param[1] > param[2] {
			return nil, errors.New("minBytesize must be low then maxBytesize")
		}
	case count > 3:
		return nil, errors.New("number of arguments must be less than 3")

	}
	param = append(param, DefaultTTL, DefaultminSize, DefaultmaxSize)
	c := &Cache{
		table:         make(map[KeyT]*list.Element),
		queue:         list.New(),
		ttl:           param[0],
		checkTTL:      uint32(math.Sqrt(float64(param[0]))),
		curSize:       0,
		minSize:       param[1+flagForIndex],
		maxSize:       param[2+count%3],
		isCollapse:    make(chan bool, 1),
		isClose:       make(chan bool),
		closeCollapse: sync.WaitGroup{},
		RWMutex:       sync.RWMutex{},
	}
	c.closeCollapse.Add(1)
	go c.collapse()
	return c, nil
}

// Deleted the cache struct.
//  Use like: defer Cache.Destroy()
func (c *Cache) Destroy() {
	if c != nil {
		c.Lock()
		if c.isClose != nil {
			close(c.isClose)
			c.closeCollapse.Wait()
			c.isClose = nil
		}
		c.Unlock()
	}
	c = nil
}

// Puts element: {key, value} to Cache or update value if key exist.
func (c *Cache) Put(key KeyT, val State) {
	c.Lock()
	defer c.Unlock()

	if e, ok := c.table[key]; ok {
		c.queue.MoveToFront(e)
		c.curSize -= e.Value.(*elem).size()
		e.Value.(*elem).val = ValueT(val)
		e.Value.(*elem).deadTime = time.Now().Add(time.Duration(c.ttl) * time.Second).Unix()
		c.curSize += e.Value.(*elem).size()
	} else {
		c.table[key] = c.queue.PushFront(&elem{key, ValueT(val),
			time.Now().Add(time.Duration(c.ttl) * time.Second).Unix()})
		c.curSize += c.queue.Front().Value.(*elem).size()
		if c.curSize >= c.maxSize {
			c.isCollapse <- true
		}
	}
}

// Returns the value by key, ok = true - if the item exists.
// If the key is not in the Cache its return nil, false.
func (c *Cache) Get(key KeyT) (val State, ok bool) {
	c.Lock()
	defer c.Unlock()

	if e, ok := c.table[key]; ok && time.Now().Unix() <= e.Value.(*elem).deadTime {
		c.queue.MoveToFront(e)
		e.Value.(*elem).deadTime = time.Now().Add(time.Duration(c.ttl) * time.Second).Unix()
		//typeVal := reflect.ValueOf(e.Value.(*elem).val).Elem() //reflect.TypeOf(e.Value.(*elem).val).PkgPath()
		return State(e.Value.(*elem).val), true
	}
	return State{}, false
}

// Displays the contents of the cache.
func (c *Cache) Display() {
	c.RLock()
	defer c.RUnlock()

	str := "{"
	for e := c.queue.Front(); e != nil; e = e.Next() {
		str += fmt.Sprintf("{%v: %v}, ", e.Value.(*elem).key, e.Value.(*elem).val)
	}
	if str != "{" {
		str = str[:len(str)-2]
	}
	str += "}"
	fmt.Println(str)
}

// Delete element from cache.
func (c *Cache) Del(key KeyT) {
	c.Lock()
	defer c.Unlock()
	e, ok := c.table[key]
	if ok { // bug: don't use Lock() after RLock()
		c.curSize -= e.Value.(*elem).size()
		delete(c.table, e.Value.(*elem).key)
		c.queue.Remove(e)
	}
}

func (c *Cache) removeElem(e *list.Element) {
	c.Lock()
	defer c.Unlock()
	c.curSize -= e.Value.(*elem).size()
	delete(c.table, e.Value.(*elem).key)
	c.queue.Remove(e)
}

func (c *Cache) collapse() {
	defer c.closeCollapse.Done()
	timer := time.NewTicker(time.Second * time.Duration(c.checkTTL))
	defer timer.Stop()

	for {
		select {
		case <-c.isCollapse:
			for e := c.queue.Back(); c.curSize >= c.maxSize && e != nil; e = e.Prev() {
				c.removeElem(e)
			}
		case <-timer.C:
			for e := c.queue.Back(); e != nil &&
				time.Now().Unix() >= e.Value.(*elem).deadTime; e = e.Prev() {
				c.removeElem(e)
			}
			runtime.GC()
		case <-c.isClose:
			return
		}
		log.Println("Collapse Cache!")
	}
}
