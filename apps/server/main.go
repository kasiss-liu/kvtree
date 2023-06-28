package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kasiss-liu/kvtree/apps/server/config"
	"github.com/kasiss-liu/kvtree/apps/server/jobs"
	"github.com/kasiss-liu/kvtree/apps/server/services"
	"github.com/kasiss-liu/kvtree/apps/server/static"
	"github.com/kasiss-liu/kvtree/apps/server/useradmin"
	"github.com/kasiss-liu/kvtree/src/module/dataset"
	"github.com/kasiss-liu/kvtree/src/module/jobm"
)

//go:embed resources/index.html resources/favicon.ico
var staticIndex embed.FS

//go:embed resources/assets/*
var staticAssets embed.FS

var (
	cnfPath = flag.String("config", "conf.toml", "a config file formatted in toml, default conf.toml")
)

func init() {
	flag.Parse()
}

func main() {
	bs, err := readfile(*cnfPath)
	if err != nil {
		fmt.Println("load config file error:", err)
		os.Exit(1)
	}
	static.ServerConf, err = config.NewServerConfFromBytes(bs)
	if err != nil {
		fmt.Println("parse config file error:", err)
		os.Exit(1)
	}

	static.UserList = *useradmin.NewUserPermit(static.ServerConf.Users)

	fmt.Println("data store dir:", static.ServerConf.ServNode.DataStore)
	static.DataStoreSet, err = dataset.NewDataSetWithFileDir(static.ServerConf.ServNode.DataStore)
	if err != nil {
		fmt.Println("init data store set error:", err)
		os.Exit(1)
	}
	static.DataStoreSet.AutoSyncd = static.ServerConf.ServNode.AutoSync
	if static.DataStoreSet.AutoSyncd {
		static.DataStoreSet.AutoSync()
		fmt.Println("data set auto sync.")
	}

	var timer *time.Ticker
	if static.ServerConf.JobCnf.Ticker > 0 {
		timer = time.NewTicker(time.Duration(static.ServerConf.JobCnf.Ticker) * time.Second)
	}
	fmt.Println("starting ticker job module...")
	static.JobExcutor = jobm.NewJobModuleExcutor(static.DataStoreSet, static.ServerConf.JobCnf.Debug, timer)
	static.JobExcutor.RegisterJobs(jobs.RunJobs)
	fmt.Println("ticker job module started ...")
	fmt.Println(len(jobs.RunJobs), "jobs registered")

	services.NewServerWithStatic(staticIndex, staticAssets).Run()
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	static.JobExcutor.Stop()
}

func readfile(p string) ([]byte, error) {
	bs, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	return bs, nil
}
