package config

import "github.com/BurntSushi/toml"

type Node struct {
	Name string `json:"name" toml:"name"`
	Port int    `json:"port" toml:"port"`
	Addr string `json:"addr" toml:"addr"`

	DataStore string `json:"datastore" toml:"datastore"`
	AutoSync  bool   `json:"autosync" toml:"autosync"`
}

type User struct {
	Name   string   `json:"name" toml:"name"`
	Passwd string   `json:"passwd" toml:"passwd"`
	Source []string `json:"source" toml:"source"`
	Super  bool     `json:"super" toml:"super"`
}

type ServerConf struct {
	ServNode *Node   `json:"server" toml:"server"`
	Users    []*User `json:"user" toml:"user"`
	JobCnf   *JobCnf `json:"job" toml:"job"`
}

type JobCnf struct {
	Ticker int64 `json:"ticker" bson:"ticker"`
	Debug  bool  `json:"debug" bson:"debug"`
}

func NewServerConfFromBytes(bs []byte) (*ServerConf, error) {
	sc := &ServerConf{}
	err := toml.Unmarshal(bs, &sc)
	if err != nil {
		return nil, err
	}
	return sc, nil
}
