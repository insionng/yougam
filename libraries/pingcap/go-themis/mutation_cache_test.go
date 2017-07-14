package themis

import (
	"bytes"

	"github.com/insionng/yougam/libraries/ngaut/log"
	. "github.com/insionng/yougam/libraries/pingcap/check"
	"github.com/insionng/yougam/libraries/pingcap/go-hbase"
)

type MutationCacheTestSuit struct{}

var _ = Suite(&MutationCacheTestSuit{})

func (s *MutationCacheTestSuit) TestMutationCache(c *C) {
	cache := newColumnMutationCache()
	row := []byte("r1")
	col := &hbase.Column{[]byte("f1"), []byte("q1")}
	cache.addMutation([]byte("tbl"), row, col, hbase.TypePut, []byte("test"), false)
	cache.addMutation([]byte("tbl"), row, col, hbase.TypeDeleteColumn, []byte("test"), false)
	cache.addMutation([]byte("tbl"), row, col, hbase.TypePut, []byte("test"), false)

	cc := &hbase.ColumnCoordinate{
		Table: []byte("tbl"),
		Row:   []byte("r1"),
		Column: hbase.Column{
			Family: []byte("f1"),
			Qual:   []byte("q1"),
		},
	}
	mutation := cache.getMutation(cc)
	if mutation == nil || mutation.typ != hbase.TypePut || bytes.Compare(mutation.value, []byte("test")) != 0 {
		c.Error("cache error")
	} else {
		log.Info(mutation)
	}

	p := hbase.NewPut([]byte("row"))
	p.AddStringValue("cf", "q", "v")
	p.AddStringValue("cf", "q1", "v")
	p.AddStringValue("cf", "q2", "v")
	p.AddStringValue("cf", "q3", "v")

	entries := getEntriesFromPut(p)
	c.Assert(len(entries), Equals, 4)
}
