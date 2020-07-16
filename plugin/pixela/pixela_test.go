package main

import (
	"testing"
)

func TestPostPixela(t *testing.T) {
	api := NewPixelaClient("tsubasaxzzz", "scheduleevent", "organza-faun-weak")

	err2 := api.PostEvent()
	if err2 != nil {
		t.Error(err2)
	}

}
