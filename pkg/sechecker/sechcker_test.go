package sechecker_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"sechecker/pkg/sechecker"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func checkMetaData(t *testing.T, expect sechecker.MetaData, target sechecker.MetaData) {
	t.Helper()

}

// メタデータファイルの読み書きのテスト
func TestWriteMetadataToMetadatafile(t *testing.T) {
	m := sechecker.MetaData{
		sechecker.StateClosed,
		0,
		[]sechecker.ScheduleEvent{
			{"602d9444-d2cd-49c7-8624-8643e7171297", "Reboot", "VirtualMachine", []string{"FrontEnd_IN_0", "BackEnd_IN_0"}, "Scheduled", "Mon, 19 Sep 2016 18:29:47 GMT"},
		},
	}
	eventJsonFile := "event_test.json"
	// ファイルの書き込みテスト
	if err := m.WriteEventFile(eventJsonFile); err != nil {
		t.Fatalf("Write metadata faild: %s", err)
	}

	// ファイルの読み取りテスト
	read_m := sechecker.ReadEventFile(eventJsonFile)

	if diff := cmp.Diff(m, read_m); diff != "" {
		t.Errorf("differs: (-got +want)\n%s", diff)

	}
	os.Remove(eventJsonFile)

}
func TestGetScheduleEvent(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected sechecker.MetaData
	}{
		//-------------------------------------------------------------------
		{
			"testA",
			`
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
		`,
			sechecker.MetaData{
				sechecker.StateClosed,
				0,
				[]sechecker.ScheduleEvent{
					{"602d9444-d2cd-49c7-8624-8643e7171297", "Reboot", "VirtualMachine", []string{"FrontEnd_IN_0", "BackEnd_IN_0"}, "Scheduled", "Mon, 19 Sep 2016 18:29:47 GMT"},
				},
			},
		},
		//-------------------------------------------------------------------
		{
			"testB",
			`
{
    "DocumentIncarnation": 0,
    "Events": []
}
		`,
			sechecker.MetaData{
				sechecker.StateClosed,
				0,
				[]sechecker.ScheduleEvent{},
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			//			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.Write([]byte(c.input))
			}))
			defer server.Close()
			api, err0 := sechecker.New(server.Client(), server.URL)
			if err0 != nil {
				t.Fatalf("Error initialize API")
			}
			metadata := sechecker.MetaData{}
			body, err := api.Run(&metadata)

			if err != nil {
				t.Errorf("Error")
			}
			t.Logf("Metadata------------->\n %#v", metadata)

			if len(metadata.Events) != len(c.expected.Events) {
				t.Fatalf("Events doesn't match: input=%d, expected=%d", len(metadata.Events), len(c.expected.Events))
			}
			if metadata.DocumentIncarnation != c.expected.DocumentIncarnation {
				t.Errorf("Does not match metadata.DocumentIncarnation")
			}
			for i, v := range metadata.Events {
				if v.EventID != c.expected.Events[i].EventID {
					t.Fatalf("EventId doesn't match")
				}
				if v.EventType != c.expected.Events[i].EventType {
					t.Fatalf("EventType doesn't match")
				}
				if v.ResourceType != c.expected.Events[i].ResourceType {
					t.Fatalf("ResourceType doesn't match")
				}
				if v.EventStatus != c.expected.Events[i].EventStatus {
					t.Fatalf("EventStatus doesn't match")
				}
				if v.NotBefore != c.expected.Events[i].NotBefore {
					t.Fatalf("NotBefore doesn't match")
				}
				if len(v.Resources) != len(c.expected.Events[i].Resources) {
					t.Fatalf("Resources doesn't match")
				}
				for j, v := range v.Resources {
					if v != c.expected.Events[i].Resources[j] {
					}
				}
			}
			t.Logf("%s", body)
		})
	}
}

func TestEventDiff(t *testing.T) {
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
			eventState := sechecker.DiffEvent(c.prevMetaData, c.currentMetaData)
			if eventState != c.expected {
				t.Fatalf("eventState doesn't match: expected=%d, actual=%d\n", c.expected, eventState)
			}
			t.Logf("Test: %s, EventState %d\n", c.name, c.expected)
		})
	}

}
func TestPostPixela(t *testing.T) {
	metadata := sechecker.MetaData{
		sechecker.StateClosed,
		0,
		[]sechecker.ScheduleEvent{
			{"", "", "", []string{"", ""}, "", ""},
		},
	}
	api, err := sechecker.NewPixelaClient("tsubasaxzzz", "scheduleevent", "organza-faun-weak")
	if err != nil {
		t.Errorf("Error")
	}

	// ToDo: グラフがない時のエラー処理
	err2 := api.PostEvent(metadata)
	if err2 != nil {
		t.Errorf("Error")
	}

}
