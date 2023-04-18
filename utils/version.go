package utils

import "runtime"

type version struct{}

var (
	Version   = "version not set, please build your app with appropriate ldflags"
	Revision  = ""
	GitSHA1   = ""
	BuildDate = ""
)

// V the version
var V = &version{}

func (v *version) String() string {
	var versionString = Version
	if GitSHA1 != "" {
		versionString += " (" + GitSHA1 + ")"
	}
	if Revision != "" {
		versionString += ", " + Revision
	}
	if BuildDate != "" {
		versionString += ", " + BuildDate
	}
	versionString += ", " + runtime.Version()
	return versionString
}
