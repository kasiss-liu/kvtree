package client

import (
	"fmt"
	"time"
)

type Task struct {
	TickGap    time.Duration
	Client     *HttpClient
	BmccKey    string
	CallFuncs  []func([]byte) error
	KeyLastVer string
}

func NewTask(key string, client *HttpClient, gap time.Duration, fn ...func([]byte) error) *Task {
	return &Task{
		Client:    client,
		TickGap:   gap,
		BmccKey:   key,
		CallFuncs: fn,
	}
}

func (task *Task) Run() error {
	data, err := task.Client.Get(task.BmccKey)
	if err != nil {
		return err
	}
	task.KeyLastVer = data.Ver
	for _, call := range task.CallFuncs {
		if call == nil {
			continue
		}
		err = call([]byte(data.String()))
		if err != nil {
			return fmt.Errorf("key %s err:%e", task.BmccKey, err)
		}
	}
	if task.TickGap > 0 {
		ticker := time.NewTicker(task.TickGap)
		go func() {
			for range ticker.C {
				data, err := task.Client.Get(task.BmccKey)
				if err != nil {
					fmt.Printf("[bmcc task]key: %s,ticker: %d,get data error: %s \n", task.BmccKey, task.TickGap, err.Error())
					continue
				}
				if task.KeyLastVer == data.Ver {
					continue
				}
				for i, call := range task.CallFuncs {
					if call == nil {
						continue
					}
					err = call([]byte(data.String()))
					if err != nil {
						fmt.Printf("[bmcc task]key: %s,ticker: %d,set data error: %s \n", task.BmccKey, task.TickGap, err.Error())
					} else {
						task.KeyLastVer = data.Ver
						fmt.Printf("[bmcc task]key: %s,ticker: %d,set data success ver: %s fn[%d]\n", task.BmccKey, task.TickGap, data.Ver, i)
					}
				}
			}
		}()
	}
	return nil
}
