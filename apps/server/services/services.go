package services

import (
	"embed"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/kasiss-liu/kvtree/apps/server/config"
	"github.com/kasiss-liu/kvtree/apps/server/static"
	"github.com/kasiss-liu/kvtree/src/module/datastore"
)

type ServiceNode struct {
	Name    string
	Port    int
	Addr    string
	Status  int //0异常 1正常
	handler http.Handler
	Store   datastore.Store
}

func (sn *ServiceNode) Handlers() map[string]http.HandlerFunc {
	return sn.handler.(*ServerHttpHandler).router
}

func (sn *ServiceNode) PrintHandlers() {
	handlers := sn.Handlers()
	paths := make([]string, 0, len(handlers))
	for p := range handlers {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	fmt.Println("hanlder list:")
	for _, p := range paths {
		if p == "" {
			continue
		}
		fmt.Println("=== " + p + " ===")
	}
}

func (sn *ServiceNode) Ping() error {
	return nil
}

func (sn *ServiceNode) Serve() error {
	go func() {
		addr := fmt.Sprintf("%s:%d", sn.Addr, sn.Port)
		fmt.Println("Server start at ", addr)
		sn.PrintHandlers()
		err := http.ListenAndServe(addr, sn.handler)
		if err != nil {
			fmt.Println("ServiceNode Serve:", err)
		}
	}()
	return nil
}

func NewServerNodeFromConf(cnf *config.Node) *ServiceNode {
	s := &ServiceNode{
		Addr:    cnf.Addr,
		Name:    cnf.Name,
		Port:    cnf.Port,
		handler: NewServerHttpHandler(),
	}
	return s
}

func NewServerNodeFromConfWithStatic(cnf *config.Node, fs ...embed.FS) *ServiceNode {
	s := &ServiceNode{
		Addr:    cnf.Addr,
		Name:    cnf.Name,
		Port:    cnf.Port,
		handler: NewServerHttpWithStaticHandler(fs...),
	}
	return s
}

type Server struct {
	SelfNode *ServiceNode
	NodeList []*ServiceNode
}

func (snl *Server) Run() {
	snl.HealthCheck()
	snl.SelfNode.Serve()
}

// 健康检查
func (snl *Server) HealthCheck() {
	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			for _, node := range snl.NodeList {
				if err := node.Ping(); err != nil {
					node.Status = 0
				} else {
					node.Status = 1
				}
			}
		}
	}()
}

func NewServer() *Server {
	s := &Server{}
	s.SelfNode = NewServerNodeFromConf(static.ServerConf.ServNode)
	return s
}

func NewServerWithStatic(fs ...embed.FS) *Server {
	s := &Server{}
	s.SelfNode = NewServerNodeFromConfWithStatic(static.ServerConf.ServNode, fs...)
	return s
}
