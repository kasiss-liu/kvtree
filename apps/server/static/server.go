package static

import (
	"github.com/kasiss-liu/kvtree/apps/server/config"
	"github.com/kasiss-liu/kvtree/apps/server/useradmin"
	"github.com/kasiss-liu/kvtree/src/module/dataset"
	"github.com/kasiss-liu/kvtree/src/module/jobm"
)

var ServerConf *config.ServerConf

var DataStoreSet *dataset.DataSet

var UserList useradmin.UsersPermit

var LoginUser = make(map[string]*useradmin.UserLogin)

var JobExcutor *jobm.JobModuleExcutor
