package datatree

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/kasiss-liu/kvtree/src/entity/dataitem"
	"github.com/kasiss-liu/kvtree/src/entity/store"
	"github.com/kasiss-liu/kvtree/src/module/datastore"
)

type DataTree struct {
	Name        string
	Root        *dataitem.DataNode
	Store       datastore.Store
	lock        sync.RWMutex
	LastModTime time.Time
	autoSyncd   bool
}

func NewDataTree(root string) *DataTree {
	dt := &DataTree{}
	dt.Name = root
	dt.Root = dataitem.NewDataNode(dt.Name, dataitem.DataMap{}, nil)
	return dt
}

func (dt *DataTree) Build() error {
	bs, err := dt.Store.Load()
	if err != nil {
		return err
	}
	if len(bs) == 0 {
		return nil
	}
	s := store.StoreItem{}
	err = json.Unmarshal(bs, &s)
	if err != nil {
		return err
	}
	dt.Unserialize(s)
	return nil
}

func (dt *DataTree) AutoSync() {
	if dt.autoSyncd {
		return
	}
	dt.autoSyncd = true
	lastSync := dt.LastModTime
	go func() {
		tk := time.NewTicker(time.Second)
		for range tk.C {
			if dt.LastModTime != lastSync {
				if err := dt.Sync(); err != nil {
					fmt.Println(err)
					continue
				}
				lastSync = dt.LastModTime
				fmt.Println(dt.Name, "auto sync modified")
				continue
			}
			// fmt.Println("no sync", dt.LastModTime, lastSync)
		}
	}()
}

func (dt *DataTree) Sync() error {
	data, err := dt.Serialize()
	if err != nil {
		return err
	}
	return dt.Store.Save(data.Bytes())
}

func (dt *DataTree) Del(key string) error {
	dt.lock.Lock()
	defer dt.lock.Unlock()
	return dt.del(key)
}

func (dt *DataTree) Set(key string, value *dataitem.DataCell) error {
	dt.lock.Lock()
	defer dt.lock.Unlock()
	return dt.set(key, value)
}
func (dt *DataTree) Get(key string) *dataitem.DataNode {
	dt.lock.RLock()
	defer dt.lock.RUnlock()
	return dt.get(key)
}

func (dt *DataTree) get(key string) *dataitem.DataNode {
	keys := strings.Split(dt.Name+"/"+strings.Trim(key, "/"), "/")
	items := dataitem.DataMap{dt.Root.Key: dt.Root}
	var v *dataitem.DataNode
	for _, key := range keys {
		if item, ok := items[key]; ok {
			v = item
			if item.Index != nil {
				if its := item.MapDataItem(); len(its) > 0 {
					items = its
				}
			}
		}
	}
	v = v.Copy()
	v.Index = nil
	return v
}

func (dt *DataTree) Prefix(key string) *dataitem.DataNode {
	dt.lock.RLock()
	defer dt.lock.RUnlock()
	keys := strings.Split(dt.Name+"/"+strings.Trim(key, "/"), "/")
	items := dataitem.DataMap{dt.Root.Key: dt.Root}
	var v *dataitem.DataNode
	for _, key := range keys {
		if item, ok := items[key]; ok && item.Index != nil {
			v = item
			if its := item.MapDataItem(); len(its) > 0 {
				items = its
			}
		}
	}
	v = v.Copy()
	return v
}

func (dt *DataTree) set(longKey string, value *dataitem.DataCell) error {
	// fmt.Println(longKey)
	keys := strings.Split(strings.Trim(longKey, "/"), "/")
	maxI := len(keys) - 1
	//不属于同一个
	var fn func(node *dataitem.DataNode, i int, value *dataitem.DataCell)
	fn = func(node *dataitem.DataNode, i int, value *dataitem.DataCell) {
		key := keys[i]
		//到达数组末尾 进行数据操作
		if maxI == i {
			//如果node是一个已存在的索引
			if m := node.MapDataItem(); m != nil {
				if _, ok := m[key]; ok {
					m[key].Val.Set(value)
				} else {
					m[key] = dataitem.NewDataNode(key, nil, value)
				}
			} else {
				//如果node下没有索引，需要为node创建一个索引，并赋值
				node.Index = dataitem.DataMap{key: dataitem.NewDataNode(key, nil, value)}
			}
			return
		}
		//查找节点
		//如果当前节点是一个索引
		if m := node.MapDataItem(); m != nil {
			//判断是否已存在当前key值,如果已经存在node 则使用此节点进行下一位操作
			if _, ok := m[key]; ok {
				node = m[key]
			} else {
				//否则为当前node索引增加索引键，创建空值
				m[key] = dataitem.NewDataNode(key, nil, nil)
				node = m[key]
			}
		} else {
			//如果不是索引，则创建索引
			//如果当前node内有值，需要将值迁移到索引中
			if node.Index != nil {
				node.Index = dataitem.DataMap{node.Key: dataitem.NewDataNode(node.Key, node.Index, nil)}
			} else {
				//否则创建空索引
				node.Index = dataitem.DataMap{key: dataitem.NewDataNode(key, nil, nil)}
			}
			m := node.MapDataItem()
			node = m[key]
		}
		fn(node, i+1, value)

	}
	fn(dt.Root, 0, value)
	dt.LastModTime = time.Now()
	return nil
}

func (dt *DataTree) del(longKey string) error {
	// fmt.Println(longKey)
	keys := strings.Split(strings.Trim(longKey, "/"), "/")
	maxI := len(keys) - 1

	// var lastNode *dataitem.DataNode
	//不属于同一个
	var fn func(node *dataitem.DataNode, i int)
	fn = func(node *dataitem.DataNode, i int) {
		key := keys[i]
		//到达数组末尾 进行数据操作
		if maxI == i {
			//尝试从索引中找到最后一个值
			if m := node.MapDataItem(); m != nil {
				if _, ok := m[key]; ok {
					//要删除的这个节点 不是叶子节点 还有后续存值
					if mm := m[key].MapDataItem(); mm != nil {
						// 只删除节点内数据
						m[key].Val.Del()
					} else {
						// 找到的这个节点是叶子终端 直接干掉叶子
						delete(m, key)
					}
				}
				//else 索引中不存在最后一个key 算了 没东西可以删
			}
			// else 没有索引就算了 key不存在
			return
		}
		//查找节点
		//如果当前节点是一个索引
		var newNode *dataitem.DataNode
		if m := node.MapDataItem(); m != nil {
			//判断是否已存在当前key值,如果已经存在node 则使用此节点进行下一位操作
			if _, ok := m[key]; ok {
				newNode = m[key]
			}
		}
		//没有找到下一个节点 退出
		if newNode == nil {
			return
		}
		fn(newNode, i+1)
	}
	fn(dt.Root, 0)
	dt.LastModTime = time.Now()
	return nil
}

func (dt *DataTree) Show() {
	dt.lock.RLock()
	defer dt.lock.RUnlock()

	var fn func(string, *dataitem.DataNode)
	fn = func(index string, root *dataitem.DataNode) {
		if root != nil {
			if root.Val != nil && len(root.Val.List) > 0 && root.Index != nil {
				fmt.Println(index, "=>", root.VersionAll())
			}
			if root.MapDataItem() != nil {
				// fmt.Println(root.MapDataItem())
				for i, v := range root.MapDataItem() {
					fn(index+"/"+i, v)
				}
			} else {
				fmt.Println(index, "=>", root.VersionAll())
			}
		}
	}
	fn(dt.Root.Key, dt.Root)
}

func (dt *DataTree) Serialize() (store.StoreItem, error) {
	dt.lock.RLock()
	defer dt.lock.RUnlock()

	m := make(map[string]*dataitem.DataLeaf)
	var fn func(string, *dataitem.DataNode)
	fn = func(index string, root *dataitem.DataNode) {
		if root != nil {
			if root.Val != nil && root.Index != nil {
				if len(root.Val.List) > 0 {
					index := strings.TrimPrefix(index, dt.Root.Key+"/")
					m[index] = root.Val
				}
			}
			if root.MapDataItem() != nil {
				for i, v := range root.MapDataItem() {
					fn(index+"/"+i, v)
				}
			} else {
				if len(root.Val.List) > 0 {
					index := strings.TrimPrefix(index, dt.Root.Key+"/")
					m[index] = root.Val
				}
			}
		}
	}
	fn(dt.Root.Key, dt.Root)
	item := store.StoreItem{}
	item.Name = dt.Name
	item.Data = m
	return item, nil
}

func (dt *DataTree) Unserialize(data store.StoreItem) {
	dt.Name = data.Name
	for key, v := range data.Data {
		sort.Sort(dataitem.DataCellList(v.List))
		for _, item := range v.List {
			dt.Set(key, item)
		}
	}
}

func (dt *DataTree) Close() {
	dt.Store.Close()
}
