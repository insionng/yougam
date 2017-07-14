package iterator_test

import (
	"testing"

	"github.com/insionng/yougam/libraries/syndtr/goleveldb/leveldb/testutil"
)

func TestIterator(t *testing.T) {
	testutil.RunSuite(t, "Iterator Suite")
}
