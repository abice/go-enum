<h1 align="center">
    <img src="https://github.com/bradleyjkemp/cupaloy/blob/master/mascot.png" alt="Mascot" width="200">
    <br>
    <a href="https://travis-ci.org/bradleyjkemp/cupaloy"><img src="https://travis-ci.org/bradleyjkemp/cupaloy.svg?branch=master" alt="Build Status" /></a>
    <a href="https://coveralls.io/github/bradleyjkemp/cupaloy?branch=master"><img src="https://coveralls.io/repos/github/bradleyjkemp/cupaloy/badge.svg" alt="Coverage Status" /></a>
    <a href="https://goreportcard.com/report/github.com/bradleyjkemp/cupaloy"><img src="https://goreportcard.com/badge/github.com/bradleyjkemp/cupaloy" alt="Go Report Card" /></a>
    <a href="https://godoc.org/github.com/bradleyjkemp/cupaloy"><img src="https://godoc.org/github.com/bradleyjkemp/cupaloy?status.svg" alt="GoDoc" /></a>
    <a href="https://sourcegraph.com/github.com/bradleyjkemp/cupaloy?badge"><img src="https://sourcegraph.com/github.com/bradleyjkemp/cupaloy/-/badge.svg" alt="Number of users" /></a>
</h1>

Incredibly simple Go snapshot testing: `cupaloy` takes a snapshot of your test output and compares it to a snapshot committed alongside your tests. If the values don't match then you'll be forced to update the snapshot file before the test passes.

Snapshot files are handled automagically: just use the `cupaloy.SnapshotT(t, value)` function in your tests and `cupaloy` will automatically find the relevant snapshot file and compare it with the given value.

### Usage
```golang
func TestExample(t *testing.T) {
    result := someFunction()

    // check that the result is the same as the last time the snapshot was updated
    // if the result has changed then the test will be failed with an error containing
    // a diff of the changes
    cupaloy.SnapshotT(t, result)
}
```

To update the snapshots simply set the ```UPDATE_SNAPSHOTS``` environment variable and run your tests e.g.
```bash
UPDATE_SNAPSHOTS=true go test ./...
```
This will fail all tests where the snapshot was updated (to stop you accidentally updating snapshots in CI) but your snapshot files will now have been updated to reflect the current output of your code.

### Installation
```bash
go get -u github.com/bradleyjkemp/cupaloy
```

### Further Examples
#### Table driven tests
```golang
var testCases = map[string][]string{
    "TestCaseOne": []string{......},
    "AnotherTestCase": []string{......},
    ....
}

func TestCases(t *testing.T) {
    for testName, args := range testCases {
        t.Run(testName, func(t *testing.T) {
            result := functionUnderTest(args...)
            cupaloy.SnapshotT(t, result)
        })
    }
}
```
#### Changing output directory
```golang
func TestSubdirectory(t *testing.T) {
    result := someFunction()
    snapshotter := cupaloy.New(cupaloy.SnapshotSubdirectory("testdata"))
    err := snapshotter.Snapshot(result)
    if err != nil {
        t.Fatalf("error: %s", err)
    }
}
```
For further usage examples see basic_test.go and advanced_test.go in the examples/ directory which are both kept up to date and run on CI.
