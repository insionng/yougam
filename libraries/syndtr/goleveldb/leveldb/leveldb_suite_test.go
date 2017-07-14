package leveldb

import (
	"testing"

	"github.com/insionng/yougam/libraries/syndtr/goleveldb/leveldb/testutil"
)

func TestLevelDB(t *testing.T) {
	testutil.RunSuite(t, "LevelDB Suite")
}
