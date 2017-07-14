package themis

import (
	"runtime"
	"strconv"
	"sync"

	"github.com/insionng/yougam/libraries/ngaut/log"
	. "github.com/insionng/yougam/libraries/pingcap/check"
	"github.com/insionng/yougam/libraries/pingcap/go-hbase"
)

type ParallelTestSuit struct{}

var _ = Suite(&ParallelTestSuit{})

func (s *ParallelTestSuit) TestParallelHbaseCall(c *C) {
	runtime.GOMAXPROCS(runtime.NumCPU() / 2)
	cli, err := createHBaseClient()
	c.Assert(err, Equals, nil)

	err = createNewTableAndDropOldTable(cli, themisTestTableName, "cf", nil)
	c.Assert(err, Equals, nil)

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			tx := newTxn(cli, defaultTxnConf)
			p := hbase.NewPut(getTestRowKey(c))
			p.AddValue(cf, q, []byte(strconv.Itoa(i)))
			tx.Put(themisTestTableName, p)
			tx.Commit()
		}(i)
	}
	wg.Wait()

	g := hbase.NewGet(getTestRowKey(c)).AddColumn(cf, q)
	rs, err := cli.Get(themisTestTableName, g)
	if err != nil {
		log.Fatal(err)
	}
	log.Info(string(rs.SortedColumns[0].Value))
}
