/*
Copyright © 2014–5 Brad Ackerman.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package test

import (
	"log"
	"os"
	"runtime"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil || os.IsExist(err) {
		return true
	}
	return false
}

func libFile(dir, ext string) string {
	return dir + "/mod_spatialite" + ext
}

// SpatialiteModulePath finds the path of the installed Spatialite
// module, for passing to the SQLite driver.
func SpatialiteModulePath() string {
	var libExtension string
	// FIXME: Does Plan9 use .so? Where does Windows put the library?
	switch runtime.GOOS {
	case "darwin":
		libExtension = ".dylib"
	case "windows":
		libExtension = ".dll"
	default:
		libExtension = ".so"
	}

	var directory string
	if fileExists(libFile("/usr/local/lib", libExtension)) {
		directory = "/usr/local/lib"
	} else if fileExists(libFile("/usr/lib64", libExtension)) {
		directory = "/usr/lib64"
	} else if fileExists(libFile("/usr/lib", libExtension)) {
		directory = "/usr/lib"
	} else {
		log.Fatalf("Unable to find valid path for mod_spatialite.")
	}

	return libFile(directory, libExtension)
}
