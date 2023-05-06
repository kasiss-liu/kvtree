package datastore

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestFileStore(t *testing.T) {
	tf := fmt.Sprintf("/tmp/bmcc_test_file_%d.dat", time.Now().Unix())
	fs, err := NewFileStore(tf)
	if err != nil {
		t.Error(err)
		return
	}
	err = fs.Save([]byte("hello world"))
	if err != nil {
		t.Error(err)
		return
	}
	fs.Close()
	fs, _ = NewFileStore(tf)
	bs, err := fs.Load()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("load", string(bs))

	t.Log("rm ", tf, os.Remove(tf))
}
