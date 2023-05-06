package datatree

import (
	"testing"

	"github.com/kasiss-liu/kvtree/src/entity/dataitem"
)

// func TestDataItemFind(t *testing.T) {

// 	data := &DataItem{
// 		Key: "a",
// 		Value: DataMap{"b": {
// 			Key: "b",
// 			Value: DataMap{"c": {
// 				Key: "c",
// 				Value: DataMap{"d": {
// 					Key: "d",
// 					Value: DataMap{"d": {
// 						Key:   "d",
// 						Value: "hello",
// 					}},
// 				}},
// 			}},
// 		}},
// 	}

// 	t.Logf("%#v", data.Get("/a/b/c/d/d/"))

// }

// func TestDataItemSet(t *testing.T) {
// 	data := NewDataItem("root", nil)
// 	data.Set("/a/b/c/d/e", "hello")
// 	data.Show()

// }
// func TestDataItemSet2(t *testing.T) {
// 	data := &DataItem{
// 		Key: "a",
// 		Value: DataMap{"b": {
// 			Key: "b",
// 			Value: DataMap{"c": {
// 				Key: "c",
// 				Value: DataMap{"dd": {
// 					Key: "dd",
// 					Value: DataMap{"d": {
// 						Key:   "d",
// 						Value: "hello",
// 					}},
// 				}},
// 			}},
// 		}},
// 	}
// 	data.Show()
// 	t.Log("---------------------")
// 	data.Set("/a/b/c/d", "hello1")
// 	data.Show()

// }

func TestDataTreeShow(t *testing.T) {
	data := DataTree{Name: "root", Root: &dataitem.DataNode{
		Key: "root",
		Index: dataitem.DataMap{
			"a": &dataitem.DataNode{
				Key: "a",
				Index: dataitem.DataMap{"bc": {
					Key: "bc",
					Index: dataitem.DataMap{"d": {
						Key:   "d",
						Index: nil,
						Val:   &dataitem.DataLeaf{List: []*dataitem.DataCell{dataitem.NewDataCell("1", "float64", 1)}},
					}},
				}},
			},
			"dd": &dataitem.DataNode{
				Key: "dd",
				Index: dataitem.DataMap{"a": {
					Key: "a",
					Index: dataitem.DataMap{"a": {
						Key:   "a",
						Index: nil,
						Val:   &dataitem.DataLeaf{List: []*dataitem.DataCell{dataitem.NewDataCell("1", "float64", 3)}},
					}},
				}},
				Val: &dataitem.DataLeaf{List: []*dataitem.DataCell{dataitem.NewDataCell("1", "string", "test_dd")}},
			},
		},
	}}
	data.Show()

}

func TestDataTree(t *testing.T) {
	m := map[string]interface{}{
		"a/b/c":   1,
		"a/b/d":   2,
		"a/bc/d":  3,
		"dd/a/a":  1,
		"a/b/c/d": 0.4,
		"a":       1,
	}

	dt := NewDataTree("root")
	// t.Log(dt.Root)
	// dt.set("a/bc/d", m["a/bc/d"])

	// dt.Show()
	// dt.set("dd/a/a", m["dd/a/a"])
	for k, v := range m {
		vv := dataitem.NewDataCell("1", "float64", v)
		dt.Set(k, vv)
	}
	dt.Show()

	se, err := dt.Serialize()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(se)

	t.Log("get before del", dt.Get("a").Value())
	t.Log("del error", dt.Del("a"))
	t.Log("get after del", dt.Get("a").Value())

	t.Log("get a/b", dt.Get("a/b").Value())

	dt.Del("a/b/c/d")
	t.Log("prefix a/b", dt.Prefix("a/b").RawJson())
}
