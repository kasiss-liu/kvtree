package dataitem

import (
	"encoding/json"
	"fmt"
	"strings"
)

type DataMap map[string]*DataNode

type DataLeaf struct {
	List []*DataCell `json:"vers,omitempty"`
}

func (dl *DataLeaf) Set(dc *DataCell) {
	for i, leaf := range dl.List {
		if leaf.Ver == dc.Ver {
			dl.List[i] = dc
			return
		}
	}
	dl.List = append([]*DataCell{dc}, dl.List[:]...)
	if len(dl.List) > 20 {
		dl.List = dl.List[:20]
	}
}

// Del 危险操作
func (dl *DataLeaf) Del() {
	dl.List = make([]*DataCell, 0)
}

type DataCell struct {
	Ver   string      `json:"ver"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

func NewDataCell(ver, t string, value interface{}) *DataCell {
	return &DataCell{Ver: ver, Type: t, Value: value}
}

func (dc *DataCell) Copy() *DataCell {
	d := &DataCell{}
	d.Type = dc.Type
	d.Ver = dc.Ver
	d.Value = dc.Value
	return d
}

type DataCellList []*DataCell

func (dcl DataCellList) Len() int {
	return len(dcl)
}
func (dcl DataCellList) Swap(i, j int) {
	dcl[i], dcl[j] = dcl[j], dcl[i]
}
func (dcl DataCellList) Less(i, j int) bool {
	return dcl[i].Ver < dcl[j].Ver
}

type DataNode struct {
	Key   string      `json:"key,omitempty"`
	Val   *DataLeaf   `json:"value,omitempty"`
	Index interface{} `json:"index,omitempty"`
}

func (di *DataNode) KeyValues() map[string]*DataLeaf {
	m := make(map[string]*DataLeaf)
	var fn func(string, *DataNode)
	fn = func(index string, root *DataNode) {
		if root != nil {
			if root.Val != nil && root.Index != nil {
				if len(root.Val.List) > 0 {
					index := strings.TrimPrefix(index, di.Key+"/")
					m[index] = root.Val
				}
			}
			if root.MapDataItem() != nil {
				for i, v := range root.MapDataItem() {
					fn(index+"/"+i, v)
				}
			} else {
				if len(root.Val.List) > 0 {
					index := strings.TrimPrefix(index, di.Key+"/")
					m[index] = root.Val
				}
			}
		}
	}
	fn(di.Key, di)
	return m
}

func NewDataNode(key string, index interface{}, value *DataCell) *DataNode {
	val := &DataLeaf{List: []*DataCell{}}
	if value != nil {
		val.List = append(val.List, value)
	}
	return &DataNode{Key: key, Index: index, Val: val}
}

// switch di.Val.Type {
// case "string":
// 	return fmt.Sprintf("string:%v", di.Val.Value)
// case "int":
// 	return fmt.Sprintf("int:%v", di.Val.Value)
// case "int64":
// 	return fmt.Sprintf("int64:%v", di.Val.Value)
// case "float64":
// 	return fmt.Sprintf("float64:%v", di.Val.Value)
// case "float":
// 	return fmt.Sprintf("float:%v", di.Val.Value)
// }
// if di.Index != nil {
// 	bs, _ := json.Marshal(di.Index)
// 	return string(bs)
// }
// return fmt.Sprintf("unexpected type value:%v", di.Val)

func (di *DataNode) Value(ver ...string) string {
	version := ""
	if len(ver) > 0 {
		version = ver[0]
	}
	var dc *DataCell
	if len(di.Val.List) > 0 {
		if version == "" {
			dc = di.Val.List[0]
		} else {
			for _, item := range di.Val.List {
				if item.Ver == version {
					dc = item.Copy()
					break
				}
			}
		}
	}
	bs := []byte{}
	if dc != nil {
		if dc.Type == "json" {
			m := make(map[string]interface{})
			json.Unmarshal([]byte(dc.Value.(string)), &m)
			dc.Value = m
			bs, _ = json.Marshal(dc)
		} else {
			bs, _ = json.Marshal(dc)
		}
	}

	return string(bs)
}
func (di *DataNode) RawJson() string {
	bs, _ := json.Marshal(di)
	return string(bs)
}

func (di *DataNode) VersionAll() string {
	bs, _ := json.Marshal(di.Val.List)
	return string(bs)
}

func (di *DataNode) MapDataItem() DataMap {
	if di.Index == nil {
		return nil
	}
	if v, ok := di.Index.(DataMap); ok {
		return v
	}
	return nil
}

func (di *DataNode) String() string {
	return fmt.Sprintf("[%s:%v:%v]", di.Key, di.Index, di.Val)
}

func (di *DataNode) Copy() *DataNode {
	d := &DataNode{}
	d.Index = di.Index
	d.Val = di.Val
	d.Key = di.Key
	return d
}
