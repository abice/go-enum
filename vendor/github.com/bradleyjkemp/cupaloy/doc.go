// Package cupaloy provides a simple api for snapshot testing in golang: test that your changes don't unexpectedly alter the results of your code.
//
// cupaloy takes a snapshot of a given value and compares it to a snapshot committed alongside your tests. If the values don't match then you'll be forced to update the snapshot file before the test passes.
//
// Snapshot files are handled automagically: just use the cupaloy.Snapshot(value) function in your tests and cupaloy will automatically find the relevant snapshot file and compare it with the given value.
//
// Installation
//   go get -u github.com/bradleyjkemp/cupaloy
//
// Usage
//  func TestExample(t *testing.T) {
//    result := someFunction()
//
//    // check that the result is the same as the last time the snapshot was updated
//    // if the result has changed then the test will be failed with an error containing
//    // a diff of the changes
//    cupaloy.SnapshotT(t, result)
//  }
//
// To update the snapshots simply set the UPDATE_SNAPSHOTS environment variable and run your tests e.g.
//   UPDATE_SNAPSHOTS=true go test ./...
// Your snapshot files will now have been updated to reflect the current output of your code.
package cupaloy
