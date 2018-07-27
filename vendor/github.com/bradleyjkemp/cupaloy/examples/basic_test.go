package examples

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
)

func TestString(t *testing.T) {
	result := "Hello world"
	err := cupaloy.Snapshot(result)
	if err != nil {
		t.Errorf("This will pass because \"Hello world\" is in the snapshot %s", err)
	}

	err = cupaloy.Snapshot("Hello world!")
	if err == nil {
		t.Errorf("Now it will fail because the snapshot doesn't have an exclamation mark %s", err)
	}
}

// Tests are independent of each other
func TestSecondString(t *testing.T) {
	result := "Hello Universe!"
	err := cupaloy.Snapshot(result)
	if err != nil {
		t.Errorf("This will pass because Snapshots are per test function %s", err)
	}
}

// Multiple snapshots can be taken in a single test
func TestMultipleSnapshots(t *testing.T) {
	result1 := "Hello"
	err := cupaloy.Snapshot(result1)
	if err != nil {
		t.Errorf("This will pass as normal %s", err)
	}

	result2 := "World"
	err = cupaloy.SnapshotMulti("result2", result2)
	if err != nil {
		t.Errorf("This will pass also as we've specified a unique (to this function) id %s", err)
	}
}

// Snapshot() takes an arbitrary number of values
func TestMultipleValues(t *testing.T) {
	result1 := "Hello"
	result2 := "World"

	err := cupaloy.Snapshot(result1, result2)
	if err != nil {
		t.Errorf("You can snapshot multiple values in the same call to Snapshot %s", err)
	}
}

// All types can be snapshotted. Maps are snapshotted in a deterministic way
func TestMap(t *testing.T) {
	result := map[int]string{
		1: "Hello",
		3: "!",
		2: "World",
	}

	err := cupaloy.Snapshot(result)
	if err != nil {
		t.Errorf("Snapshots can be taken of any type %s", err)
	}
}
