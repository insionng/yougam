package themis

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/insionng/yougam/libraries/ngaut/log"
	. "github.com/insionng/yougam/libraries/pingcap/check"
	"github.com/insionng/yougam/libraries/pingcap/go-hbase"
)

type TransactionTestSuit struct {
	cli hbase.HBaseClient
}

var _ = Suite(&TransactionTestSuit{})

func Test(t *testing.T) { TestingT(t) }

func (s *TransactionTestSuit) SetUpSuite(c *C) {
	var err error
	s.cli, err = createHBaseClient()
	c.Assert(err, Equals, nil)

	log.Warn("new test, reset tables")
	err = createNewTableAndDropOldTable(s.cli, themisTestTableName, string(cf), nil)
	c.Assert(err, IsNil)
}

func (s *TransactionTestSuit) SetUpTest(c *C) {
}

func getTestRowKey(c *C) []byte {
	return []byte("test_row_" + c.TestName())
}

func (s *TransactionTestSuit) TestAsyncCommit(c *C) {
	conf := defaultTxnConf
	conf.brokenCommitSecondaryTest = true

	tx := newTxn(s.cli, conf)
	// simulating broken commit
	for i := 0; i < 10; i++ {
		p := hbase.NewPut([]byte(fmt.Sprintf("test_%d", i)))
		p.AddValue(cf, q, []byte(fmt.Sprintf("%d", tx.(*themisTxn).GetStartTS())))
		tx.Put(themisTestTableName, p)
	}
	err := tx.Commit()
	c.Assert(err, Equals, nil)

	//  wait until lock expired.
	log.Warn("Wait for lock expired. Sleep...")
	tick := 6
	for tick > 0 {
		time.Sleep(1 * time.Second)
		tick--
		log.Infof("remain %ds...", tick)
	}

	log.Warn("Try commit again")
	// new transction will not see lock
	for {
		tx = newTxn(s.cli, defaultTxnConf)
		for i := 0; i < 5; i++ {
			p := hbase.NewPut([]byte(fmt.Sprintf("test_%d", i)))
			p.AddValue(cf, q, []byte(fmt.Sprintf("%d", tx.(*themisTxn).GetStartTS())))
			tx.Put(themisTestTableName, p)
		}
		err = tx.Commit()
		if err == nil || !errorEqual(err, ErrRetryable) {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	c.Assert(err, Equals, nil)

	for {
		tx = newTxn(s.cli, defaultTxnConf)
		for i := 5; i < 10; i++ {
			p := hbase.NewPut([]byte(fmt.Sprintf("test_%d", i)))
			p.AddValue(cf, q, []byte(fmt.Sprintf("%d", tx.(*themisTxn).GetStartTS())))
			tx.Put(themisTestTableName, p)
		}
		err = tx.Commit()
		if err == nil || !errorEqual(err, ErrRetryable) {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	c.Assert(err, Equals, nil)
}

func (s *TransactionTestSuit) TestBrokenPrewriteSecondary(c *C) {
	tx := newTxn(s.cli, defaultTxnConf)
	ts := tx.(*themisTxn).GetStartTS()
	// simulating broken commit
	for i := 0; i < 10; i++ {
		p := hbase.NewPut([]byte(fmt.Sprintf("test_%d", i)))
		p.AddValue(cf, q, []byte(fmt.Sprintf("%d", ts)))
		tx.Put(themisTestTableName, p)
	}
	err := tx.Commit()
	c.Assert(err, IsNil)

	// TODO: check rallback & cleanup locks
	conf := defaultTxnConf
	conf.brokenPrewriteSecondaryTest = true

	tx = newTxn(s.cli, conf)
	ts = tx.GetStartTS()
	// simulating broken commit
	for i := 0; i < 10; i++ {
		p := hbase.NewPut([]byte(fmt.Sprintf("test_%d", i)))
		p.AddValue(cf, q, []byte(fmt.Sprintf("%d", ts)))
		tx.Put(themisTestTableName, p)
	}
	err = tx.Commit()
	c.Assert(err, NotNil)

	// check if locks are cleaned successfully
	tx = newTxn(s.cli, defaultTxnConf)
	for i := 0; i < 10; i++ {
		g := hbase.NewGet([]byte(fmt.Sprintf("test_%d", i)))
		r, err := tx.Get(themisTestTableName, g)
		c.Assert(err, Equals, nil)
		c.Assert(r == nil || string(r.SortedColumns[0].Value) != fmt.Sprintf("%d", ts), Equals, true)
	}
}

func (s *TransactionTestSuit) TestPrimaryLockTimeout(c *C) {
	// TODO: check if lock can be cleaned up when secondary prewrite failed and
	// rollback is also failed
	conf := defaultTxnConf
	conf.brokenPrewriteSecondaryTest = true
	conf.brokenPrewriteSecondaryAndRollbackTest = true
	tx := newTxn(s.cli, conf)
	ts := tx.GetStartTS()
	// simulating broken commit
	for i := 0; i < 2; i++ {
		p := hbase.NewPut([]byte(fmt.Sprintf("test_%d", i)))
		p.AddValue(cf, q, []byte(fmt.Sprintf("%d", ts)))
		tx.Put(themisTestTableName, p)
	}
	err := tx.Commit()
	c.Assert(err, NotNil)
	log.Error(err)

	//  wait until lock expired.
	log.Warn("Wait for lock expired. Sleep...")
	tick := 6
	for tick > 0 {
		time.Sleep(1 * time.Second)
		tick--
		log.Infof("remain %ds...", tick)
	}

	// check if locks are cleaned successfully
	tx = newTxn(s.cli, defaultTxnConf)
	for i := 0; i < 2; i++ {
		g := hbase.NewGet([]byte(fmt.Sprintf("test_%d", i)))
		r, err := tx.Get(themisTestTableName, g)
		c.Assert(err, Equals, nil)
		// this commit must rollback
		c.Assert(r == nil || string(r.SortedColumns[0].Value) != fmt.Sprintf("%d", ts), Equals, true)
	}
}

func checkCommitSuccess(s *TransactionTestSuit, c *C, row []byte) {
	tx := newTxn(s.cli, defaultTxnConf)
	colMap := make(map[string]string)
	colMap["#p:"+string(cf)+"#q"] = ""
	colMap[string(cf)+":q"] = ""
	r, err := tx.(*themisTxn).client.Get(themisTestTableName, hbase.NewGet(row))
	c.Assert(err, Equals, nil)
	c.Assert(2, Equals, len(r.Columns))
	for _, v := range r.Columns {
		_, exist := colMap[string(v.Family)+":"+string(v.Qual)]
		c.Assert(exist, Equals, true)
	}
}

func (s *TransactionTestSuit) TestLockRow(c *C) {
	tx := newTxn(s.cli, defaultTxnConf)
	row := []byte("lockRow")
	put := hbase.NewPut(row)
	put.AddValue(cf, q, []byte("v"))
	tx.Put(themisTestTableName, put)
	tx.Commit()

	checkCommitSuccess(s, c, row)

	tx = newTxn(s.cli, defaultTxnConf)
	err := tx.LockRow(themisTestTableName, row)
	c.Assert(err, Equals, nil)

	tx.(*themisTxn).selectPrimaryAndSecondaries()
	err = tx.(*themisTxn).prewritePrimary()
	c.Assert(err, Equals, nil)
	colMap := make(map[string]string)
	colMap["#p:"+string(cf)+"#q"] = ""
	colMap[string(cf)+":q"] = ""
	colMap["L:"+string(cf)+"#q"] = ""
	var r *hbase.ResultRow
	r, err = tx.(*themisTxn).client.Get(themisTestTableName, hbase.NewGet(row))
	c.Assert(err, Equals, nil)
	c.Assert(3, Equals, len(r.Columns))
	for _, v := range r.Columns {
		_, exist := colMap[string(v.Family)+":"+string(v.Qual)]
		c.Assert(exist, Equals, true)
	}
	tx.(*themisTxn).commitTs = tx.GetStartTS() + 1
	tx.(*themisTxn).commitPrimary()
	checkCommitSuccess(s, c, row)
}

func (s *TransactionTestSuit) TestBatchGet(c *C) {
	batchGetTestTbl := "batch_get_test"
	err := createNewTableAndDropOldTable(s.cli, batchGetTestTbl, string(cf), [][]byte{
		// split in middle
		[]byte("batch_test_5"),
	})
	defer dropTable(s.cli, batchGetTestTbl)
	// prepare data
	tx := newTxn(s.cli, defaultTxnConf)
	for i := 0; i < 10; i++ {
		p := hbase.NewPut([]byte(fmt.Sprintf("batch_test_%d", i))).AddValue(cf, q, []byte("v"))
		tx.Put(batchGetTestTbl, p)
	}
	err = tx.Commit()
	c.Assert(err, IsNil)

	// batch get
	var gets []*hbase.Get
	for i := 0; i < 10; i++ {
		g := hbase.NewGet([]byte(fmt.Sprintf("batch_test_%d", i))).AddColumn(cf, q)
		gets = append(gets, g)
	}
	for i := 5; i < 10; i++ {
		g := hbase.NewGet([]byte(fmt.Sprintf("batch_test_no_such_row_%d", i))).AddColumn(cf, q)
		gets = append(gets, g)
	}
	tx = newTxn(s.cli, defaultTxnConf)
	_, err = tx.Gets(batchGetTestTbl, gets)
	c.Assert(isWrongRegionErr(err), Equals, true)

	gets = nil
	for i := 0; i < 5; i++ {
		g := hbase.NewGet([]byte(fmt.Sprintf("batch_test_%d", i))).AddColumn(cf, q)
		gets = append(gets, g)
	}
	tx = newTxn(s.cli, defaultTxnConf)
	rs, err := tx.Gets(batchGetTestTbl, gets)
	c.Assert(err, IsNil)
	c.Assert(len(rs), Equals, 5)
}

func (s *TransactionTestSuit) TestBatchGetWithLocks(c *C) {
	// simulating locks
	conf := defaultTxnConf
	conf.brokenCommitSecondaryTest = true

	tx := newTxn(s.cli, conf)
	ts := tx.GetStartTS()
	// simulating broken commit
	for i := 0; i < 10; i++ {
		p := hbase.NewPut([]byte(fmt.Sprintf("batch_test_with_lock_%d", i)))
		p.AddValue(cf, q, []byte(fmt.Sprintf("%d", ts)))
		tx.Put(themisTestTableName, p)
	}
	tx.Commit()

	tx = newTxn(s.cli, defaultTxnConf)

	var gets []*hbase.Get
	for i := 0; i < 10; i++ {
		g := hbase.NewGet([]byte(fmt.Sprintf("batch_test_with_lock_%d", i))).AddColumn(cf, q)
		gets = append(gets, g)
	}
	rs, err := tx.Gets(themisTestTableName, gets)
	c.Assert(err, IsNil)
	// we had already cleaned secondary locks
	c.Assert(len(rs), Equals, 10)
}

func (s *TransactionTestSuit) TestAsyncSecondaryCommit(c *C) {
	conf := defaultTxnConf
	conf.brokenCommitSecondaryTest = true
	tx := newTxn(s.cli, conf)
	for i := 0; i < 10; i++ {
		p := hbase.NewPut([]byte(fmt.Sprintf("async_commit_test_%d", i))).AddValue(cf, q, []byte(fmt.Sprintf("%d", tx.GetStartTS())))
		tx.Put(themisTestTableName, p)
	}
	err := tx.Commit()
	c.Assert(err, IsNil)

	tx = newTxn(s.cli, conf)
	for i := 0; i < 10; i++ {
		g := hbase.NewGet([]byte(fmt.Sprintf("async_commit_test_%d", i)))
		rs, err := tx.Get(themisTestTableName, g)
		c.Assert(err, IsNil)
		c.Assert(len(rs.SortedColumns), Greater, 0)
	}
}

func (s *TransactionTestSuit) TestTTL(c *C) {
	conf := defaultTxnConf
	conf.brokenCommitPrimaryTest = true
	tx := newTxn(s.cli, conf)
	p := hbase.NewPut(getTestRowKey(c)).AddValue(cf, q, []byte("val"))
	tx.Put(themisTestTableName, p)
	tx.Commit()

	startTs := time.Now()
	conf = defaultTxnConf
	conf.TTLInMs = 1000
	tx = newTxn(s.cli, conf)
	rs, err := tx.Get(themisTestTableName, hbase.NewGet(getTestRowKey(c)).AddColumn(cf, q))
	c.Assert(time.Since(startTs).Seconds(), Greater, float64(1))
	c.Assert(time.Since(startTs).Seconds(), Less, float64(1.5))
	// transction timeout, alreay rolled back.
	c.Assert(rs, IsNil)
	c.Assert(err, IsNil)
	tx.Commit()
}

type mockOracle struct {
	tick uint64
}

func (o *mockOracle) GetTimestamp() (uint64, error) {
	return o.tick, nil
}

func (o *mockOracle) IsExpired(beginMs uint64, TTL uint64) bool {
	return false
}

func (s *TransactionTestSuit) TestPhantomRead(c *C) {
	o := &mockOracle{}

	conf := defaultTxnConf
	conf.brokenCommitPrimaryTest = true
	o.tick = 1
	tx, _ := NewTxnWithConf(s.cli, conf, o)
	p := hbase.NewPut(getTestRowKey(c)).AddValue(cf, q, []byte("val"))
	o.tick = 3
	tx.Put(themisTestTableName, p)
	tx.Commit()

	o.tick = 2
	tx, _ = NewTxn(s.cli, o)
	rs, err := tx.Get(themisTestTableName, hbase.NewGet(getTestRowKey(c)).AddColumn(cf, q))
	c.Assert(err, NotNil)
	c.Assert(rs, IsNil)
}

func (s *TransactionTestSuit) TestExceedMaxRows(c *C) {
	conf := defaultTxnConf
	tx := newTxn(s.cli, conf)
	for i := 0; i < conf.MaxRowsInOneTxn+1; i++ {
		tx.Put(themisTestTableName, hbase.NewPut([]byte(strconv.Itoa(i))).AddValue(cf, q, []byte("test")))
	}
	err := tx.Commit()
	c.Assert(err, NotNil)
}

func (s *TransactionTestSuit) TestCheckCommitStatus(c *C) {
	conf := defaultTxnConf
	hook := newHook()
	hook.beforeCommitSecondary = func(txn *themisTxn, ctx interface{}) (bool, interface{}, error) {
		// add before commit secondary hook, just return, do not commit
		// secondaries
		log.Info("before commit secondary")
		return false, nil, nil
	}
	tx := newTxn(s.cli, conf)
	tx.(*themisTxn).setHook(hook)
	for i := 0; i < 10; i++ {
		tx.Put(themisTestTableName, hbase.NewPut([]byte(fmt.Sprintf("%s_%d", getTestRowKey(c), i))).AddValue(cf, q, []byte("test")))
	}
	tx.Commit()
	commitTs := tx.(*themisTxn).commitTs

	pRow := append([]byte(nil), tx.(*themisTxn).primaryRow.row...)
	sRow := append([]byte(nil), tx.(*themisTxn).secondaryRows[0].row...)

	for i := 0; i < 10; i++ {
		conf = defaultTxnConf
		tx = newTxn(s.cli, conf)
		tx.Put(themisTestTableName, hbase.NewPut([]byte(pRow)).AddValue(cf, q, []byte("newVal")))
		tx.Commit()
	}

	hook = newHook()
	hook.beforePrewriteLockClean = func(txn *themisTxn, ctx interface{}) (bool, interface{}, error) {
		lock := ctx.(Lock)
		cc := lock.Primary().Coordinate()
		exists, err := txn.lockCleaner.IsLockExists(lock.Coordinate(), 0, lock.Timestamp())
		c.Assert(err, IsNil)
		c.Assert(exists, Equals, true)
		ts, err := txn.lockCleaner.GetCommitTimestamp(cc, lock.Timestamp())
		c.Assert(err, IsNil)
		c.Assert(ts, Equals, commitTs)
		return true, nil, nil
	}
	tx = newTxn(s.cli, conf)
	tx.(*themisTxn).setHook(hook)
	tx.Put(themisTestTableName, hbase.NewPut([]byte(sRow)).AddValue(cf, q, []byte("newVal")))
	tx.Commit()
}

func createChoosePrimaryRowHook(target string) *txnHook {
	hook := newHook()
	hook.afterChoosePrimaryAndSecondary = func(txn *themisTxn, ctx interface{}) (bool, interface{}, error) {
		txn.secondary = nil
		txn.secondaryRows = nil
		for tblName, rowMutations := range txn.mutationCache.mutations {
			for _, rowMutation := range rowMutations {
				row := rowMutation.row
				for i, mutation := range rowMutation.mutationList(true) {
					colcord := hbase.NewColumnCoordinate([]byte(tblName), row, mutation.Family, mutation.Qual)
					// set the first column as primary if primary is not set by user
					if string(row) == target {
						txn.primary = colcord
						txn.primaryRowOffset = i
						txn.primaryRow = rowMutation
					} else {
						txn.secondary = append(txn.secondary, colcord)
					}
				}
				if string(row) != target {
					txn.secondaryRows = append(txn.secondaryRows, rowMutation)
				}
			}
		}
		return true, nil, nil
	}
	return hook
}

// Fix https://yougam/libraries/pingcap/go-themis/issues/19
func (s *TransactionTestSuit) TestPrewriteSecondaryMissingRows(c *C) {
	conf := defaultTxnConf
	hook := createChoosePrimaryRowHook("A")
	hook.beforePrewriteSecondary = func(txn *themisTxn, ctx interface{}) (bool, interface{}, error) {
		go func() {
			hook2 := createChoosePrimaryRowHook("B")
			hook2.onSecondaryOccursLock = func(txn *themisTxn, ctx interface{}) (bool, interface{}, error) {
				log.Info("tx2 occurs secondary lock", ctx)
				return true, nil, nil
			}
			tx2 := newTxn(s.cli, conf)
			tx2.(*themisTxn).setHook(hook2)
			tx2.Put(themisTestTableName, hbase.NewPut([]byte("A")).AddValue(cf, q, []byte("A")))
			tx2.Put(themisTestTableName, hbase.NewPut([]byte("B")).AddValue(cf, q, []byte("B")))
			tx2.Put(themisTestTableName, hbase.NewPut([]byte("C")).AddValue(cf, q, []byte("C")))
			tx2.Commit()
		}()
		time.Sleep(500 * time.Millisecond)
		return true, nil, nil
	}

	hook.onSecondaryOccursLock = func(txn *themisTxn, ctx interface{}) (bool, interface{}, error) {
		log.Info("tx1", ctx)
		return true, nil, nil
	}

	hook.onPrewriteRow = func(txn *themisTxn, ctx interface{}) (bool, interface{}, error) {
		containPrimary := ctx.([]interface{})[1].(bool)
		if !containPrimary {
			rm := ctx.([]interface{})[0].(*rowMutation)
			log.Info(string(rm.row))
		}

		return true, nil, nil
	}

	tx1 := newTxn(s.cli, conf)
	tx1.(*themisTxn).setHook(hook)
	tx1.Put(themisTestTableName, hbase.NewPut([]byte("A")).AddValue(cf, q, []byte("A")))
	tx1.Put(themisTestTableName, hbase.NewPut([]byte("B")).AddValue(cf, q, []byte("B")))
	tx1.Put(themisTestTableName, hbase.NewPut([]byte("C")).AddValue(cf, q, []byte("C")))
	err := tx1.Commit()
	c.Assert(err, IsNil)

	tx3 := newTxn(s.cli, conf)
	rs, err := tx3.Get(themisTestTableName, hbase.NewGet([]byte("C")).AddColumn(cf, q))
	c.Assert(rs, NotNil)
	c.Assert(err, IsNil)
	tx3.Commit()
}
