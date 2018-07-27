package cupaloy

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

// New constructs a new, configured instance of cupaloy using the given Configurators.
func New(configurators ...Configurator) *Config {
	return defaultConfig().WithOptions(configurators...)
}

// Snapshot calls Snapshotter.Snapshot with the default config.
func Snapshot(i ...interface{}) error {
	return defaultConfig().snapshot(getNameOfCaller(), i...)
}

// SnapshotMulti calls Snapshotter.SnapshotMulti with the default config.
func SnapshotMulti(snapshotID string, i ...interface{}) error {
	snapshotName := fmt.Sprintf("%s-%s", getNameOfCaller(), snapshotID)
	return defaultConfig().snapshot(snapshotName, i...)
}

// SnapshotT calls Snapshotter.SnapshotT with the default config.
func SnapshotT(t *testing.T, i ...interface{}) {
	t.Helper()
	defaultConfig().SnapshotT(t, i...)
}

// Snapshot compares the given value to the it's previous value stored on the filesystem.
// An error containing a diff is returned if the snapshots do not match.
// Snapshot determines the snapshot file automatically from the name of the calling function.
func (c *Config) Snapshot(i ...interface{}) error {
	return c.snapshot(getNameOfCaller(), i...)
}

// SnapshotMulti is identical to Snapshot but can be called multiple times from the same function.
// This is done by providing a unique snapshotId for each invocation.
func (c *Config) SnapshotMulti(snapshotID string, i ...interface{}) error {
	snapshotName := fmt.Sprintf("%s-%s", getNameOfCaller(), snapshotID)
	return c.snapshot(snapshotName, i...)
}

// SnapshotT is identical to Snapshot but gets the snapshot name using
// t.Name() and calls t.Fail() directly if the snapshots do not match.
func (c *Config) SnapshotT(t *testing.T, i ...interface{}) {
	t.Helper()
	if t.Failed() {
		return
	}

	snapshotName := strings.Replace(t.Name(), "/", "-", -1)
	err := c.snapshot(snapshotName, i...)
	if err != nil {
		t.Error(err)
	}
}

// WithOptions allows the modification of an existing Config. This can usefully be
// used to use a different option for a single call e.g.
//  snapshotter.WithOptions(cupaloy.SnapshotSubdirectory("testdata")).SnapshotT(t, result)
func (c *Config) WithOptions(configurators ...Configurator) *Config {
	clonedConfig := c.clone()

	for _, configurator := range configurators {
		configurator(clonedConfig)
	}

	return clonedConfig
}

func (c *Config) snapshot(snapshotName string, i ...interface{}) error {
	snapshot := takeSnapshot(i...)

	prevSnapshot, err := c.readSnapshot(snapshotName)
	if os.IsNotExist(err) {
		return c.updateSnapshot(snapshotName, snapshot)
	}
	if err != nil {
		return err
	}

	if snapshot == prevSnapshot || takeV1Snapshot(i...) == prevSnapshot {
		// previous snapshot matches current value
		return nil
	}

	if c.shouldUpdate() {
		// updates snapshot to current value and upgrades snapshot format
		return c.updateSnapshot(snapshotName, snapshot)
	}

	diff := diffSnapshots(prevSnapshot, snapshot)
	return fmt.Errorf("snapshot not equal:\n%s", diff)
}
