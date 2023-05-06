package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/kasiss-liu/kvtree/src/entity/dataitem"
	"github.com/kasiss-liu/kvtree/src/module/dataset"
)

var dt *dataset.DataSet

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
	tree := dt.GetTree(item.Namespace)
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

func SetKey(w http.ResponseWriter, r *http.Request) {
	item := NewHttpReqItem(r)
	if err := item.CheckRequired(); err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	if err := dt.SetTree(item.Namespace); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	tree := dt.GetTree(item.Namespace)
	if tree != nil {
		if err := tree.Set(item.Key, item.RealValue()); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
	}
	w.Write([]byte("ok"))
}

func main() {

	http.HandleFunc("/get", GetKey)
	http.HandleFunc("/set", SetKey)

	dt = dataset.NewDataSetMem()

	http.ListenAndServe(":8080", nil)
}
