package jobm

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/kasiss-liu/kvtree/src/module/dataset"
)

type TestJobM struct {
	counter int
}

func (jm *TestJobM) Ready() bool {
	a := rand.Float64()
	fmt.Println("ready", a)
	return a > 0.5
}

func (jm *TestJobM) Run(ds *dataset.DataSet) error {
	a := rand.Float64()
	fmt.Println("run", a)
	if a > 0.5 {
		return fmt.Errorf("test error")
	}
	jm.counter++
	return nil
}

func (jm *TestJobM) Alarm(err error) {
	fmt.Println(err)
}

func TestJobExcutor(t *testing.T) {

	je := NewJobModuleExcutor(nil, true, nil)
	je.RegisterJob(&TestJobM{})

	je.Start()

	time.Sleep(5 * time.Second)
	je.Stop()
	t.Log(je.jobs)

}
