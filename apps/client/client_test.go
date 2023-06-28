package client

import "testing"

func TestBmccClient(t *testing.T) {
	a := "http://127.0.0.1"
	n := "biddsp"
	au := "testuser:123qwe"
	c := NewBmccHttpClient(a, au, n)

	t.Log(c.Get("dspadmin/config.json"))

}
