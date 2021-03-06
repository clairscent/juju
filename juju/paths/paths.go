package paths

import (
	"os"
	"os/exec"

	"github.com/juju/errors"
	jujuos "github.com/juju/utils/os"
	"github.com/juju/utils/series"
)

type osVarType int

const (
	tmpDir osVarType = iota
	logDir
	dataDir
	storageDir
	confDir
	jujuRun
	certDir
	metricsSpoolDir
	uniterStateDir
	jujuDumpLogs
)

var nixVals = map[osVarType]string{
	tmpDir:          "/tmp",
	logDir:          "/var/log",
	dataDir:         "/var/lib/juju",
	storageDir:      "/var/lib/juju/storage",
	confDir:         "/etc/juju",
	jujuRun:         "/usr/bin/juju-run",
	jujuDumpLogs:    "/usr/bin/juju-dumplogs",
	certDir:         "/etc/juju/certs.d",
	metricsSpoolDir: "/var/lib/juju/metricspool",
	uniterStateDir:  "/var/lib/juju/uniter/state",
}

var winVals = map[osVarType]string{
	tmpDir:          "C:/Juju/tmp",
	logDir:          "C:/Juju/log",
	dataDir:         "C:/Juju/lib/juju",
	storageDir:      "C:/Juju/lib/juju/storage",
	confDir:         "C:/Juju/etc",
	jujuRun:         "C:/Juju/bin/juju-run.exe",
	jujuDumpLogs:    "C:/Juju/bin/juju-dumplogs.exe",
	certDir:         "C:/Juju/certs",
	metricsSpoolDir: "C:/Juju/lib/juju/metricspool",
	uniterStateDir:  "C:/Juju/lib/juju/uniter/state",
}

// osVal will lookup the value of the key valname
// in the apropriate map, based on the series. This will
// help reduce boilerplate code
func osVal(ser string, valname osVarType) (string, error) {
	os, err := series.GetOSFromSeries(ser)
	if err != nil {
		return "", err
	}
	switch os {
	case jujuos.Windows:
		return winVals[valname], nil
	default:
		return nixVals[valname], nil
	}
}

// TempDir returns the path on disk to the corect tmp directory
// for the series. This value will be the same on virtually
// all linux systems, but will differ on windows
func TempDir(series string) (string, error) {
	return osVal(series, tmpDir)
}

// LogDir returns filesystem path the directory where juju may
// save log files.
func LogDir(series string) (string, error) {
	return osVal(series, logDir)
}

// DataDir returns a filesystem path to the folder used by juju to
// store tools, charms, locks, etc
func DataDir(series string) (string, error) {
	return osVal(series, dataDir)
}

// MetricsSpoolDir returns a filesystem path to the folder used by juju
// to store metrics.
func MetricsSpoolDir(series string) (string, error) {
	return osVal(series, metricsSpoolDir)
}

// CertDir returns a filesystem path to the folder used by juju to
// store certificates that are added by default to the Juju client
// api certificate pool.
func CertDir(series string) (string, error) {
	return osVal(series, certDir)
}

// StorageDir returns a filesystem path to the folder used by juju to
// mount machine-level storage.
func StorageDir(series string) (string, error) {
	return osVal(series, storageDir)
}

// ConfDir returns the path to the directory where Juju may store
// configuration files.
func ConfDir(series string) (string, error) {
	return osVal(series, confDir)
}

// JujuRun returns the absolute path to the juju-run binary for
// a particular series.
func JujuRun(series string) (string, error) {
	return osVal(series, jujuRun)
}

// JujuDumpLogs returns the absolute path to the juju-dumplogs binary
// for a particular series.
func JujuDumpLogs(series string) (string, error) {
	return osVal(series, jujuDumpLogs)
}

func MustSucceed(s string, e error) string {
	if e != nil {
		panic(e)
	}
	return s
}

var osStat = os.Stat
var execLookPath = exec.LookPath

// MongorestorePath will look for a mongorestore binary and return it
// if it exists. It will look first for the juju provided one and if
// not found make a try at a system installed one.
func MongorestorePath() (string, error) {
	// TODO (perrito666) this seems to be a packaging decision we should not
	// rely on it and we should be aware of /usr/lib/juju if its something
	// of ours.
	const mongoRestoreFullPath string = "/usr/lib/juju/bin/mongorestore"

	if _, err := osStat(mongoRestoreFullPath); err == nil {
		return mongoRestoreFullPath, nil
	}

	path, err := execLookPath("mongorestore")
	if err != nil {
		return "", errors.Trace(err)
	}
	return path, nil
}
