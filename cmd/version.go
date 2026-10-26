package cmd

import (
	"runtime/debug"
)

// version returns the version of the moqueries module as recorded by the Go
// toolchain. go install module@version records the module version. Builds
// from a VCS checkout record a version derived from tags (Go 1.24+) along
// with the revision.
func version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	ver := info.Main.Version
	if ver == "" {
		ver = "unknown"
	}

	var revision, modified string
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			revision = s.Value
		case "vcs.modified":
			modified = s.Value
		}
	}

	if revision == "" {
		return ver
	}

	if modified == "true" {
		revision += ", modified"
	}

	return ver + " (" + revision + ")"
}
