package dataset

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/kasiss-liu/kvtree/src/module/datastore"
	"github.com/kasiss-liu/kvtree/src/module/datatree"
)

type DataSet struct {
	List      map[string]*datatree.DataTree
	AutoSyncd bool
	DataDir   string
	lock      sync.Mutex
}

func NewDataSetWithFileDir(dir string) (*DataSet, error) {
	fs, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dir, 0755)
			fs, err = os.ReadDir(dir)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	ds := &DataSet{}
	ds.DataDir = dir
	ds.List = make(map[string]*datatree.DataTree)
	for _, f := range fs {
		fname := dir + "/" + f.Name()
		if f.IsDir() {
			continue
		}
		ar := strings.Split(f.Name(), ".")
		dt := datatree.NewDataTree(ar[0])
		dt.Store, err = datastore.NewFileStore(fname)
		if err != nil {
			return nil, fmt.Errorf("filename %s init error:%w", fname, err)
		}
		err = dt.Build()
		if err != nil {
			return nil, fmt.Errorf("filename %s init error:%w", fname, err)
		}
		fmt.Println(fname, "data key-values:")
		dt.Show()

		ds.List[dt.Name] = dt
	}
	return ds, nil
}

func NewDataSetMem() *DataSet {
	ds := &DataSet{}
	ds.List = make(map[string]*datatree.DataTree)
	return ds
}

func (ds *DataSet) SetTree(name string) error {
	ds.lock.Lock()
	defer ds.lock.Unlock()
	var err error
	if _, ok := ds.List[name]; !ok {
		dt := datatree.NewDataTree(name)
		if ds.DataDir != "" {
			fname := ds.DataDir + "/" + name + ".dat"
			dt.Store, err = datastore.NewFileStore(fname)
			if err != nil {
				return fmt.Errorf("filename %s init error:%w", fname, err)
			}
			err = dt.Build()
			if err != nil {
				return fmt.Errorf("filename %s init error:%w", fname, err)
			}
			if ds.AutoSyncd {
				dt.AutoSync()
			}
		}
		ds.List[dt.Name] = dt
	}
	return nil
}

func (ds *DataSet) GetTree(name string) *datatree.DataTree {
	return ds.List[name]
}

func (ds *DataSet) AutoSync() {
	if ds.AutoSyncd {
		for _, d := range ds.List {
			d.AutoSync()
		}
	}
}
