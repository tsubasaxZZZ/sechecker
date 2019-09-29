package main

import (
	"github.com/google/go-cmp/cmp"
	"os"
	"sechecker/pkg/sechecker"
	"testing"
)

// イベントファイル読み込みテスト
func TestCmd_ReadEventFile(t *testing.T) {
	file, err := os.OpenFile(`event.json`, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		t.Fatal("File open error")
	}
	defer file.Close()

	file.Write(([]byte)(`
{
    "DocumentIncarnation": 0,
    "Events": [
        {
            "EventId": "602d9444-d2cd-49c7-8624-8643e7171297",
            "EventType": "Reboot",
            "ResourceType": "VirtualMachine",
            "Resources": ["FrontEnd_IN_0", "BackEnd_IN_0"],
            "EventStatus": "Scheduled",
            "NotBefore": "Mon, 19 Sep 2016 18:29:47 GMT"
        }
    ]
}
	`))
	var metadata sechecker.MetaData
	readEventFile("event.json", &metadata)

	expected := sechecker.MetaData{
		sechecker.StateClosed,
		0,
		[]sechecker.ScheduleEvent{
			{"602d9444-d2cd-49c7-8624-8643e7171297", "Reboot", "VirtualMachine", []string{"FrontEnd_IN_0", "BackEnd_IN_0"}, "Scheduled", "Mon, 19 Sep 2016 18:29:47 GMT"},
		},
	}

	if diff := cmp.Diff(metadata, expected); diff != "" {
		t.Errorf("differs: (-got +want)\n%s", diff)

	}

	os.Remove(`event.json`)

}

// 前回のイベント情報と今回のイベント情報の比較テスト
func TestCmd_EventDiff(t *testing.T) {
	cases := []struct {
		name            string
		prevMetaData    sechecker.MetaData
		currentMetaData sechecker.MetaData
		expected        sechecker.EventState
	}{
		//-------------------------------------------------------------------
		{
			"イベント数: 0 -> 1",
			sechecker.MetaData{sechecker.StateClosed, 0, []sechecker.ScheduleEvent{}},
			sechecker.MetaData{sechecker.StateClosed, 0,
				[]sechecker.ScheduleEvent{
					{"602d9444-d2cd-49c7-8624-8643e7171297", "Reboot", "VirtualMachine", []string{"FrontEnd_IN_0", "BackEnd_IN_0"}, "Scheduled", "Mon, 19 Sep 2016 18:29:47 GMT"},
				},
			},
			sechecker.StateNew,
		},
		//-------------------------------------------------------------------
		{
			"イベント数: 1 -> 3",
			sechecker.MetaData{sechecker.StateClosed, 0,
				[]sechecker.ScheduleEvent{
					{"602d9444-d2cd-49c7-8624-8643e7171297", "Reboot", "VirtualMachine", []string{"FrontEnd_IN_0", "BackEnd_IN_0"}, "Scheduled", "Mon, 19 Sep 2016 18:29:47 GMT"},
				},
			},
			sechecker.MetaData{sechecker.StateClosed, 0,
				[]sechecker.ScheduleEvent{
					{"602d9444-d2cd-49c7-8624-8643e7171297", "Reboot", "VirtualMachine", []string{"FrontEnd_IN_0", "BackEnd_IN_0"}, "Scheduled", "Mon, 19 Sep 2016 18:29:47 GMT"},
					{"602d9444-d2cd-49c7-8624-8643e7171297", "Reboot", "VirtualMachine", []string{"FrontEnd_IN_0", "BackEnd_IN_0"}, "Scheduled", "Mon, 19 Sep 2016 18:29:47 GMT"},
					{"602d9444-d2cd-49c7-8624-8643e7171297", "Reboot", "VirtualMachine", []string{"FrontEnd_IN_0", "BackEnd_IN_0"}, "Scheduled", "Mon, 19 Sep 2016 18:29:47 GMT"},
				},
			},
			sechecker.StateNew,
		},
		//-------------------------------------------------------------------
		{
			"イベント数: 1 -> 1",
			sechecker.MetaData{sechecker.StateClosed, 0,
				[]sechecker.ScheduleEvent{
					{"602d9444-d2cd-49c7-8624-8643e7171297", "Reboot", "VirtualMachine", []string{"FrontEnd_IN_0", "BackEnd_IN_0"}, "Scheduled", "Mon, 19 Sep 2016 18:29:47 GMT"},
				},
			},
			sechecker.MetaData{sechecker.StateClosed, 0,
				[]sechecker.ScheduleEvent{
					{"602d9444-d2cd-49c7-8624-8643e7171297", "Reboot", "VirtualMachine", []string{"FrontEnd_IN_0", "BackEnd_IN_0"}, "Scheduled", "Mon, 19 Sep 2016 18:29:47 GMT"},
				},
			},
			sechecker.StateActive,
		},
		//-------------------------------------------------------------------
		{
			"イベント数: 1 -> 0",
			sechecker.MetaData{sechecker.StateClosed, 0,
				[]sechecker.ScheduleEvent{
					{"602d9444-d2cd-49c7-8624-8643e7171297", "Reboot", "VirtualMachine", []string{"FrontEnd_IN_0", "BackEnd_IN_0"}, "Scheduled", "Mon, 19 Sep 2016 18:29:47 GMT"},
				},
			},
			sechecker.MetaData{sechecker.StateClosed, 0, []sechecker.ScheduleEvent{}},
			sechecker.StateClosed,
		},
		//-------------------------------------------------------------------
		{
			"イベント数: 3 -> 2",
			sechecker.MetaData{sechecker.StateClosed, 0,
				[]sechecker.ScheduleEvent{
					{"602d9444-d2cd-49c7-8624-8643e7171297", "Reboot", "VirtualMachine", []string{"FrontEnd_IN_0", "BackEnd_IN_0"}, "Scheduled", "Mon, 19 Sep 2016 18:29:47 GMT"},
					{"602d9444-d2cd-49c7-8624-8643e7171297", "Reboot", "VirtualMachine", []string{"FrontEnd_IN_0", "BackEnd_IN_0"}, "Scheduled", "Mon, 19 Sep 2016 18:29:47 GMT"},
					{"602d9444-d2cd-49c7-8624-8643e7171297", "Reboot", "VirtualMachine", []string{"FrontEnd_IN_0", "BackEnd_IN_0"}, "Scheduled", "Mon, 19 Sep 2016 18:29:47 GMT"},
				},
			},
			sechecker.MetaData{sechecker.StateClosed, 0,
				[]sechecker.ScheduleEvent{
					{"602d9444-d2cd-49c7-8624-8643e7171297", "Reboot", "VirtualMachine", []string{"FrontEnd_IN_0", "BackEnd_IN_0"}, "Scheduled", "Mon, 19 Sep 2016 18:29:47 GMT"},
					{"602d9444-d2cd-49c7-8624-8643e7171297", "Reboot", "VirtualMachine", []string{"FrontEnd_IN_0", "BackEnd_IN_0"}, "Scheduled", "Mon, 19 Sep 2016 18:29:47 GMT"},
				},
			},
			sechecker.StateActive,
		},
		//-------------------------------------------------------------------
		{
			"イベント数: 0 -> 0",
			sechecker.MetaData{sechecker.StateClosed, 0, []sechecker.ScheduleEvent{}},
			sechecker.MetaData{sechecker.StateClosed, 0, []sechecker.ScheduleEvent{}},
			sechecker.StateClosed,
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			eventState := diffEvent(c.prevMetaData, c.currentMetaData)
			if eventState != c.expected {
				t.Fatalf("eventState doesn't match: expected=%d, actual=%d\n", c.expected, eventState)
			}
			t.Logf("Test: %s, EventState %d\n", c.name, c.expected)
		})
	}

}

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
