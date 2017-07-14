package oracles

import (
	"testing"
	"time"
)

func TestLocalOracle(t *testing.T) {
	l := NewLocalOracle()
	m := map[uint64]struct{}{}
	for i := 0; i < 100000; i++ {
		ts, err := l.GetTimestamp()
		if err != nil {
			t.Error(err)
		}
		m[ts] = struct{}{}
	}

	if len(m) != 100000 {
		t.Error("generated same ts")
	}
}

func TestExpired(t *testing.T) {
	o := NewLocalOracle()
	ts, _ := o.GetTimestamp()
	time.Sleep(1 * time.Second)
	if !o.IsExpired(uint64(ts), 500) {
		t.Error("should expired")
	}
	if o.IsExpired(uint64(ts), 2000) {
		t.Error("should not expired")
	}
}
