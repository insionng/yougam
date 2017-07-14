package table

import (
	"testing"

	"github.com/insionng/yougam/libraries/syndtr/goleveldb/leveldb/testutil"
)

func TestTable(t *testing.T) {
	testutil.RunSuite(t, "Table Suite")
}
