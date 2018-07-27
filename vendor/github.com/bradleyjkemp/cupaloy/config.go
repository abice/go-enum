package cupaloy

// Configurator is a functional option that can be passed to cupaloy.New() to change snapshotting behaviour.
type Configurator func(*Config)

// EnvVariableName can be used to customize the environment variable that determines whether snapshots should be updated
// e.g.
//  cupaloy.New(EnvVariableName("UPDATE"))
// Will create an instance where snapshots will be updated if the UPDATE environment variable is set,
// instead of the default of UPDATE_SNAPSHOTS.
func EnvVariableName(name string) Configurator {
	return func(c *Config) {
		c.shouldUpdate = func() bool {
			return envVariableSet(name)
		}
	}
}

// ShouldUpdate can be used to provide custom logic to decide whether or not to update a snapshot
// e.g.
//   var update = flag.Bool("update", false, "update snapshots")
//   cupaloy.New(ShouldUpdate(func () bool { return *update })
// Will create an instance where snapshots are updated if the --update flag is passed to go test.
func ShouldUpdate(f func() bool) Configurator {
	return func(c *Config) {
		c.shouldUpdate = f
	}
}

// SnapshotSubdirectory can be used to customize the location that snapshots are stored in.
// e.g.
//  cupaloy.New(SnapshotSubdirectory("testdata"))
// Will create an instance where snapshots are stored in testdata/ rather than the default .snapshots/
func SnapshotSubdirectory(name string) Configurator {
	return func(c *Config) {
		c.subDirName = name
	}
}

// Config provides the same snapshotting functions with additional configuration capabilities.
type Config struct {
	shouldUpdate func() bool
	subDirName   string
}

func defaultConfig() *Config {
	return (&Config{}).WithOptions(
		SnapshotSubdirectory(".snapshots"),
		EnvVariableName("UPDATE_SNAPSHOTS"),
	)
}

func (c *Config) clone() *Config {
	return &Config{
		shouldUpdate: c.shouldUpdate,
		subDirName:   c.subDirName,
	}
}
