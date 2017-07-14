package oracles

import "testing"

func TestRemoteOracle(t *testing.T) {
	oracle := NewRemoteOracle("localhost", "/zk/tso")
	m := map[uint64]struct{}{}
	for i := 0; i < 100000; i++ {
		ts, err := oracle.GetTimestamp()
		if err != nil {
			t.Error(err)
		}
		m[ts] = struct{}{}
	}

	if len(m) != 100000 {
		t.Error("generated same ts, ", len(m))
	}
}
