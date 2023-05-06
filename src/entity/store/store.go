package store

import (
	"encoding/json"

	"github.com/kasiss-liu/kvtree/src/entity/dataitem"
)

type StoreItem struct {
	Name string
	Data map[string]*dataitem.DataLeaf
}

func (si *StoreItem) Bytes() []byte {
	bs, _ := json.Marshal(si)
	return bs
}
