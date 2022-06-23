package log

import (
	"github.com/google/uuid"
	badger "github.com/sharop/nopaldb/store"
	"io/ioutil"
	"os"
	"testing"
	"time"

	api "github.com/sharop/nopaldb/pb/v1/log"
	"github.com/stretchr/testify/require"
)

func TestLog(t *testing.T) {

	for scenario, fn := range map[string]func(
		t *testing.T, log *Log,
	){
		"append and read a record succeeds": testAppendRead,
		"init with existing segments":       testInitExisting,
		"inserting items to database":       testLog_Insert,
	} {
		t.Run(scenario, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "store-test")
			require.NoError(t, err)
			defer os.RemoveAll(dir)

			c := Config{}
			log, err := NewLog(dir, c)
			require.NoError(t, err)

			fn(t, log)
		})
	}
}

func testLog_Insert(t *testing.T, log *Log) {
	value := `{SOMEDATA:"SomeJSON"}`
	item := badger.Item{
		ID:       uuid.New(),
		Path:     "/resources/sources",
		Type:     badger.Source,
		Version:  1,
		Key:      "TestingSource",
		Value:    &value,
		Created:  time.Now(),
		Modified: time.Now(),
		//Meta:     []string{"SOME", "META"},
	}
	err := log.Insert(item)
	require.NoError(t, err)

	item2 := badger.Item{
		ID:       uuid.New(),
		Path:     "/resources/spaces",
		Type:     badger.Space,
		Version:  1,
		Key:      "TestingSpace",
		Value:    &value,
		Created:  time.Now(),
		Modified: time.Now(),
		//Meta:     []string{"NAME", "TEST Space"},
	}
	err = log.Insert(item2)
	require.NoError(t, err)

	require.NoError(t, log.Close())

	n, err := NewLog(log.Dir, log.Config)
	require.NoError(t, err)

	var field = "Path"
	var itemRead []badger.Item
	itemRead, err = n.FindBy(field, "/resources/sources")
	require.NoError(t, err)
	require.True(t, len(itemRead) == 1)
	require.Equal(t, "TestingSource", itemRead[0].Key)

	require.NoError(t, n.Delete("TestingSource"))
	require.NoError(t, n.Delete("TestingSpace"))
}

func testAppendRead(t *testing.T, log *Log) {

	appendRecord := &api.Record{
		Value: "hello world",
	}
	off, err := log.Set(appendRecord)
	require.NoError(t, err)

	read, err := log.Get(off)
	require.NoError(t, err)
	require.Equal(t, appendRecord.Value, read.Value)
}

func testInitExisting(t *testing.T, o *Log) {

	testRecord := &api.Record{
		Key:   "First",
		Value: "Hello WW III",
	}
	for i := 0; i < 3; i++ {
		_, err := o.Set(testRecord)
		require.NoError(t, err)
	}
	require.NoError(t, o.Close())

	n, err := NewLog(o.Dir, o.Config)
	require.NoError(t, err)

	record, err := n.Get("First")
	require.NoError(t, err)
	require.Equal(t, "Hello WW III", record.Value)
}
