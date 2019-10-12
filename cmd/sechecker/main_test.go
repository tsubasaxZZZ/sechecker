package main

import (
	"sechecker/pkg/sechecker"
	"testing"
)


// 前回のイベント情報と今回のイベント情報の比較テスト

func TestCmd_DoAction(t *testing.T) {
	cases := []struct {
		name            string
		currentMetaData sechecker.MetaData
	}{
		{
			"TestA",
			sechecker.MetaData{sechecker.StateClosed, 0, []sechecker.ScheduleEvent{}},
		},
		{
			"TestB",
			sechecker.MetaData{sechecker.StateClosed, 0,
				[]sechecker.ScheduleEvent{
					{"602d9444-d2cd-49c7-8624-8643e7171297", "Reboot", "VirtualMachine", []string{"FrontEnd_IN_0", "BackEnd_IN_0"}, "Scheduled", "Mon, 19 Sep 2016 18:29:47 GMT"},
					{"602d9444-d2cd-49c7-8624-8643e7171297", "Reboot", "VirtualMachine", []string{"FrontEnd_IN_0", "BackEnd_IN_0"}, "Scheduled", "Mon, 19 Sep 2016 18:29:47 GMT"},
				},
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			if err := doAction(c.currentMetaData); err != nil {
				t.Fatal("doAction is failed")
			}
		})
	}

}
