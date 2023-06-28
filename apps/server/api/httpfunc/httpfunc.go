package httpfunc

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/kasiss-liu/kvtree/apps/server/static"
	"github.com/kasiss-liu/kvtree/src/entity/dataitem"
)

type HttpReqItem struct {
	Namespace string `json:"namespace"`
	Key       string `json:"key"`
	Type      string `json:"type"`
	Val       []byte `json:"value"`
	Version   string `json:"ver"`
}

func (i HttpReqItem) RealValue() *dataitem.DataCell {
	ver := strconv.FormatInt(time.Now().Unix(), 10)
	switch i.Type {
	case "string":
		return dataitem.NewDataCell(ver, "string", string(i.Val))
	case "int":
		ii, _ := strconv.ParseInt(string(i.Val), 10, 64)
		return dataitem.NewDataCell(ver, "int", ii)
	case "float":
		ii, _ := strconv.ParseFloat(string(i.Val), 64)
		return dataitem.NewDataCell(ver, "float", ii)
	case "json":
		return dataitem.NewDataCell(ver, "json", string(i.Val))
	}

	return dataitem.NewDataCell(ver, "unknow", i.Val)
}

func (i HttpReqItem) CheckRequired() error {
	if i.Namespace == "" || i.Key == "" || i.Type == "" {
		return fmt.Errorf("data error")
	}
	return nil
}

func (i HttpReqItem) CheckDelRequired() error {
	if i.Namespace == "" || i.Key == "" {
		return fmt.Errorf("data error")
	}
	return nil
}

func NewHttpReqItem(r *http.Request) *HttpReqItem {
	i := &HttpReqItem{}
	switch r.Method {
	case http.MethodGet:
		i.Namespace = r.URL.Query().Get("namespace")
		i.Key = r.URL.Query().Get("key")
		i.Version = r.URL.Query().Get("ver")
	case http.MethodPost:
		i.Namespace = r.URL.Query().Get("namespace")
		i.Key = r.URL.Query().Get("key")
		i.Type = r.URL.Query().Get("type")
		bs, _ := io.ReadAll(r.Body)
		defer r.Body.Close()
		i.Val = bs
	}
	return i
}

func GetKey(w http.ResponseWriter, r *http.Request) {
	item := NewHttpReqItem(r)
	tree := static.DataStoreSet.GetTree(item.Namespace)
	data := ""
	keyfound := false
	if tree != nil {
		node := tree.Get(item.Key)
		if node != nil {
			keyfound = true
			data = node.Value(item.Version)
		}
	}
	if !keyfound {
		w.WriteHeader(404)
	}
	w.Write([]byte(data))
}

func GetPrefix(w http.ResponseWriter, r *http.Request) {
	item := NewHttpReqItem(r)
	tree := static.DataStoreSet.GetTree(item.Namespace)
	data := ""
	keyfound := false
	if tree != nil {
		node := tree.Prefix(item.Key)
		if node != nil {
			keyfound = true
			data = node.RawJson()
		}
	}
	if !keyfound {
		w.WriteHeader(404)
	}
	w.Write([]byte(data))
}

func SetKey(w http.ResponseWriter, r *http.Request) {
	item := NewHttpReqItem(r)
	if err := item.CheckRequired(); err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	if err := static.DataStoreSet.SetTree(item.Namespace); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	tree := static.DataStoreSet.GetTree(item.Namespace)
	if tree != nil {
		if err := tree.Set(item.Key, item.RealValue()); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
	}
	w.Write([]byte("ok"))
}

func DataSync(w http.ResponseWriter, r *http.Request) {
	item := NewHttpReqItem(r)
	tree := static.DataStoreSet.GetTree(item.Namespace)
	data := "document not found"
	if tree != nil {
		err := tree.Sync()
		if err != nil {
			data = err.Error()
		} else {
			data = "ok"
		}
	}
	w.Write([]byte(data))
}
