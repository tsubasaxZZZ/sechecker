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
	config := sechecker.Configs{
		[]sechecker.ActionConfig{
			{"Pixela1", "Pixela", &sechecker.PixelaConfig{UserID: "tsubasaxzzz", GraphID: "scheduleevent", Secret: "organza-faun-weak"}},
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			// t.Parallel()
			if acctionCount, err := doAction(c.currentMetaData, config); err != nil || acctionCount == 0 { // エラーがある時もしくはアクションが実行されなかったとき
				t.Fatalf("Error=%s, AcctionCount=%d", err, acctionCount)
			}
		})
	}

	// コンフィグファイルが空の時はアクションは実行されない
	emptyConfig := sechecker.Configs{
		[]sechecker.ActionConfig{},
	}
	t.Run("コンフィグファイルが空の時はアクション実行されない", func(t *testing.T) {
		if acctionCount, err := doAction(cases[0].currentMetaData, emptyConfig); err != nil || acctionCount != 0 { // エラーがある時もしくはアクションが実行されたとき
			t.Fatalf("Error=%s, AcctionCount=%d", err, acctionCount)
		}
	})

}
