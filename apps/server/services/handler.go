package services

import (
	"embed"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/kasiss-liu/kvtree/apps/server/api/httpfunc"
	"github.com/kasiss-liu/kvtree/apps/server/api/httpfunc/httpweb"
	"github.com/kasiss-liu/kvtree/apps/server/static"
)

type ServerHttpHandler struct {
	router map[string]http.HandlerFunc
}

func (serv *ServerHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if handler, ok := serv.router[r.URL.Path]; ok {
		serv.safeRun(
			handler,
		)(w, r)
		return
	}
	w.WriteHeader(404)
	w.Write([]byte("not found."))
}

func (serv *ServerHttpHandler) safeRun(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
				fmt.Println(string(debug.Stack()))
				w.WriteHeader(500)
				w.Write([]byte("server internal error"))
			}
		}()
		h(w, r)
	}
}

func (serv *ServerHttpHandler) checkAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		authi := strings.Split(auth, ":")
		if len(authi) != 2 {
			w.WriteHeader(401)
			w.Write([]byte("Auth info not allowed"))
			return
		}
		name, passwd := authi[0], authi[1]
		user := static.UserList.Get(name)
		if user == nil || user.Passwd != passwd {
			w.WriteHeader(401)
			w.Write([]byte("Auth info not allowed"))
			return
		}
		if _, ok := user.NameSpaceAuth[r.URL.Query().Get("namespace")]; !ok && !user.Super {
			w.WriteHeader(401)
			w.Write([]byte("Auth info not allowed"))
			return
		}

		h(w, r)
	}
}

func NewServerHttpHandler() *ServerHttpHandler {
	s := &ServerHttpHandler{}
	router := make(map[string]http.HandlerFunc)
	router["/auth/get"] = s.checkAuth(httpfunc.GetKey)
	router["/auth/set"] = s.checkAuth(httpfunc.SetKey)
	router["/auth/sync"] = s.checkAuth(httpfunc.DataSync)
	router["/auth/prefix"] = s.checkAuth(httpfunc.GetPrefix)

	router["/api/login"] = httpweb.ApiLogin
	router["/api/logout"] = httpweb.ApiLogout
	router["/api/namespace_list"] = httpweb.ApiNamespaceList
	router["/api/key_all"] = httpweb.ApiGetAllKey
	router["/api/key_set"] = httpweb.ApiSetKey
	router["/api/key_del"] = httpweb.ApiDelKey
	router["/api/key_detail"] = httpweb.ApiKeyVersionAll
	s.router = router
	return s
}

func NewServerHttpWithStaticHandler(fs ...embed.FS) *ServerHttpHandler {
	s := &ServerHttpHandler{}
	router := make(map[string]http.HandlerFunc)
	router["/auth/get"] = s.checkAuth(httpfunc.GetKey)
	router["/auth/set"] = s.checkAuth(httpfunc.SetKey)
	router["/auth/sync"] = s.checkAuth(httpfunc.DataSync)
	router["/auth/prefix"] = s.checkAuth(httpfunc.GetPrefix)

	router["/api/login"] = httpweb.ApiLogin
	router["/api/logout"] = httpweb.ApiLogout
	router["/api/namespace_list"] = httpweb.ApiNamespaceList
	router["/api/key_all"] = httpweb.ApiGetAllKey
	router["/api/key_set"] = httpweb.ApiSetKey
	router["/api/key_del"] = httpweb.ApiDelKey
	router["/api/key_detail"] = httpweb.ApiKeyVersionAll

	RegisterStaticRouter(router, fs)
	s.router = router
	return s
}
