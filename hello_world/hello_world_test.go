package hello_world_test

import (
	"testing"

	"github.com/real-time-footfall-analysis/rtfa-backend/hello_world"
)

func TestVerboseAdder(t *testing.T) {
	if hello_world.VerboseAdder(5, 3) != 8 {
		t.Fail()
	}
}
