package git

import "testing"

func TestGetTime(t *testing.T) {
	gt := &GitTag{
		path: "../glstats-sample",
	}

	_, err := gt.Analyze("deploy")
	if err != nil {
		t.Logf("Error")
	}
}
