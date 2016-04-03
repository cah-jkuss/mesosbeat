package beater

// Defaults for config variables which are not set
const (
	DefaultRootDir string = "/etc/mesos-slave"
	DefaultUrl string = "http://localhost:5051/metrics/snapshot"
)

type MesosbeatConfig struct {
	Period 			*int64
	URLs			[]string
	RootDir 		string 
}

type ConfigSettings struct {
	Input MesosbeatConfig
}
