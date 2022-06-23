package store

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"sort"
	"testing"
)

func TestDataKey(t *testing.T) {
	var uid uint64
	// Example.
	// :resource:source:0-ventas:01

	// key with uid = 0 is invalid
	uid = 0
	key := DataKey(RKey("bad uid", Source), uid)
	_, err := Parse(key)
	require.Error(t, err)

	for uid = 1; uid < 1001; uid++ {
		// Use the uid to derive the attribute so it has variable length and the test
		// can verify that multiple sizes of attr work correctly.
		sattr := fmt.Sprintf("attr:%d", uid)
		key := DataKey(RKey(sattr, Source), uid)
		pk, err := Parse(key)
		require.NoError(t, err)

		require.True(t, pk.IsData())
		require.Equal(t, sattr, ParseAttr(pk.Attr))
		require.Equal(t, uid, pk.UId)
	}

	keys := make([]string, 0, 1024)
	for uid = 1024; uid >= 1; uid-- {
		key := DataKey(RKey("q4_2021.key", Space), uid)
		keys = append(keys, string(key))
	}
	// Test that sorting is as expected.
	sort.Strings(keys)
	require.True(t, sort.StringsAreSorted(keys))
	for i, key := range keys {
		exp := DataKey(RKey("q4_2021.key", Space), uint64(i+1))
		require.Equal(t, string(exp), key)
	}

	for uid = 1000; uid < 2001; uid++ {
		// Use the uid to derive the attribute so it has variable length and the test
		// can verify that multiple sizes of attr work correctly.
		sattr := fmt.Sprintf("%s", uuid.New().String())
		key := DataKey(RKey(sattr, Source), uid)
		pk, err := Parse(key)
		require.NoError(t, err)

		require.True(t, pk.IsData())
		require.Equal(t, sattr, ParseAttr(pk.Attr))
		require.Equal(t, uid, pk.UId)
	}

}
