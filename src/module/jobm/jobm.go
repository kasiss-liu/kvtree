package jobm

import (
	"log"
	"reflect"
	"time"

	"github.com/kasiss-liu/kvtree/src/module/dataset"
)

type JobModule interface {
	Ready() bool
	Run(*dataset.DataSet) error
	Alarm(error)
}

type JobModuleExcutor struct {
	data  *dataset.DataSet
	jobs  []JobModule
	Debug bool
	timer *time.Ticker
}

func NewJobModuleExcutor(data *dataset.DataSet, debug bool, ticker *time.Ticker) *JobModuleExcutor {
	je := &JobModuleExcutor{}
	je.Debug = debug
	je.data = data
	je.timer = ticker
	if je.timer == nil {
		je.timer = time.NewTicker(time.Second)
	}
	return je
}

func (je *JobModuleExcutor) RegisterJob(job JobModule) {
	je.jobs = append(je.jobs, job)
}
func (je *JobModuleExcutor) RegisterJobs(job []JobModule) {
	je.jobs = append(je.jobs, job...)
}

func (je *JobModuleExcutor) Start() {
	go func() {
		for range je.timer.C {
			for _, job := range je.jobs {
				go func(job JobModule) {
					if job.Ready() {
						if err := job.Run(je.data); err != nil {
							job.Alarm(err)
						}
					}
					if je.Debug {
						log.Println(reflect.TypeOf(job).Elem(), false)
					}
				}(job)
			}
		}
	}()
}

func (je *JobModuleExcutor) Stop() {
	je.timer.Stop()
	if je.Debug {
		log.Println("timer stopped")
	}
}
