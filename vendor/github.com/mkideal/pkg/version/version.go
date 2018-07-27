package version

import (
	"fmt"
)

// Version interface
type Version interface {
	App() string
	Major() int16
	Minor() int16
	Build() int64
}

// implementes Version
type versionImpl struct {
	app   string
	major int16
	minor int16
	build int64
}

// New creates a version
func New(app string, major, minor int16, build int64) Version {
	return versionImpl{
		app:   app,
		major: major,
		minor: minor,
		build: build,
	}
}

// implementes Version interface
func (v versionImpl) App() string  { return v.app }
func (v versionImpl) Major() int16 { return v.major }
func (v versionImpl) Minor() int16 { return v.minor }
func (v versionImpl) Build() int64 { return v.build }

func (v versionImpl) String() string {
	return fmt.Sprintf("%s v%d.%d.%d", v.app, v.major, v.minor, v.build)
}
