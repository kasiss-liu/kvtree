package httpweb

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"github.com/kasiss-liu/kvtree/apps/server/api/httpfunc"
	"github.com/kasiss-liu/kvtree/apps/server/static"
	"github.com/kasiss-liu/kvtree/apps/server/useradmin"
)

type LoginItem struct {
	Name   string `json:"username" toml:"username"`
	Passwd string `json:"password" toml:"password"`
}

func NewLoginItem(r *http.Request) (*LoginItem, error) {
	defer r.Body.Close()
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	ni := &LoginItem{}
	err = json.Unmarshal(bs, &ni)
	if err != nil {
		return nil, err
	}
	return ni, nil
}

func ApiLogin(w http.ResponseWriter, r *http.Request) {
	li, err := NewLoginItem(r)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	auth := static.UserList.Get(li.Name)
	if auth == nil || auth.Passwd != li.Passwd {
		w.WriteHeader(401)
		w.Write([]byte("user info error"))
		return
	}
	hash := md5.New()
	hash.Write([]byte(fmt.Sprintf("%s%s%d", li.Name, li.Passwd, time.Now().Unix())))
	token := fmt.Sprintf("%x", hash.Sum(nil))
	static.LoginUser[token] = &useradmin.UserLogin{Name: li.Name, Token: token, Expired: 86400, Time: time.Now()}
	w.WriteHeader(200)
	w.Write([]byte(token))
}

func ApiLogout(w http.ResponseWriter, r *http.Request) {
	token := getLoginToken(r)
	delete(static.LoginUser, token)
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}

func ApiNamespaceList(w http.ResponseWriter, r *http.Request) {
	if err := checkLoginToken(r); err != nil {
		w.WriteHeader(401)
		w.Write([]byte("login user invalid"))
		return
	}
	token := getLoginToken(r)
	auth := static.LoginUser[token]
	if auth == nil {
		w.WriteHeader(401)
		w.Write([]byte("user info error"))
		return
	}
	user := static.UserList.Get(auth.Name)
	namelist := make([]string, 0)
	if user.Super {
		for _, tree := range static.DataStoreSet.List {
			namelist = append(namelist, tree.Name)
		}
	}
	for k := range user.NameSpaceAuth {
		if !checkin(namelist, k) {
			namelist = append(namelist, k)
		}
	}

	sort.Strings(namelist)
	bs, _ := json.Marshal(namelist)
	w.Write(bs)

}
func checkin(arr []string, i string) bool {
	for _, ii := range arr {
		if i == ii {
			return true
		}
	}
	return false
}

type ApiKeyValueItem struct {
	Key   string      `json:"key"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

func ApiGetAllKey(w http.ResponseWriter, r *http.Request) {
	if err := checkLoginToken(r); err != nil {
		w.WriteHeader(401)
		w.Write([]byte("login user invalid"))
		return
	}
	ns := r.URL.Query().Get("namespace")
	if ns == "" {
		w.WriteHeader(404)
		w.Write([]byte("no namespace found"))
		return
	}
	tree := static.DataStoreSet.GetTree(ns)
	data := []byte{}
	keyfound := false
	if tree != nil {
		node := tree.Prefix("/")
		if node != nil {
			keyfound = true
			m := node.KeyValues()
			mm := make([]ApiKeyValueItem, 0)
			keys := make([]string, 0)
			for k := range m {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, key := range keys {
				var v interface{} = ""
				typ := ""
				if len(m[key].List) > 0 {
					v = m[key].List[0].Value
					typ = m[key].List[0].Type
				}
				mm = append(mm, ApiKeyValueItem{Key: key, Type: typ, Value: v})
			}
			bs, _ := json.Marshal(mm)
			data = bs
		}
	}
	if !keyfound {
		w.WriteHeader(404)
	}
	w.Write(data)
}

func ApiSetKey(w http.ResponseWriter, r *http.Request) {
	if err := checkLoginToken(r); err != nil {
		w.WriteHeader(401)
		w.Write([]byte("login user invalid"))
		return
	}
	item := httpfunc.NewHttpReqItem(r)
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
			w.Write([]byte(err.Error()))
			return
		}
	}
	w.Write([]byte("ok"))
}

func ApiDelKey(w http.ResponseWriter, r *http.Request) {
	if err := checkLoginToken(r); err != nil {
		w.WriteHeader(401)
		w.Write([]byte("login user invalid"))
		return
	}
	item := httpfunc.NewHttpReqItem(r)
	if err := item.CheckDelRequired(); err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	tree := static.DataStoreSet.GetTree(item.Namespace)
	if tree != nil {
		if err := tree.Del(item.Key); err != nil {
			w.Write([]byte(err.Error()))
			return
		}
	}
	w.Write([]byte("ok"))
}

func ApiKeyVersionAll(w http.ResponseWriter, r *http.Request) {
	if err := checkLoginToken(r); err != nil {
		w.WriteHeader(401)
		w.Write([]byte("login user invalid"))
		return
	}
	ns := r.URL.Query().Get("namespace")
	if ns == "" {
		w.WriteHeader(404)
		w.Write([]byte("no namespace found"))
		return
	}
	key := r.URL.Query().Get("key")
	if key == "" {
		w.WriteHeader(403)
		w.Write([]byte("need key name"))
		return
	}
	tree := static.DataStoreSet.GetTree(ns)
	data := ""
	keyfound := false
	if tree != nil {
		node := tree.Get(key)
		if node != nil {
			keyfound = true
			data = node.ValueAll().String()
		}
	}
	if !keyfound {
		w.WriteHeader(404)
	}
	w.Write([]byte(data))
}

func getLoginToken(r *http.Request) string {
	return r.Header.Get("Token")
}

func checkLoginToken(r *http.Request) error {
	t := getLoginToken(r)
	if v, ok := static.LoginUser[t]; ok && time.Since(v.Time) < time.Duration(v.Expired*int(time.Second)) {
		return nil
	}
	return fmt.Errorf("user token invalid")
}
