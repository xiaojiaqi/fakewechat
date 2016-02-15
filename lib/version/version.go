package version

import (
	"fmt"
)

type Version uint64

var ProgramVersion = "No Version Provided"
var Buildstamp = "No Provided"
var Githash = "No Provided"

func PrintVersion() {
	fmt.Printf("Version: %s Build: %s  Git: %s\n", ProgramVersion, Buildstamp, Githash)
}
