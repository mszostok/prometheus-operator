// Copyright The prometheus-operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Holds the CLI related version functions that unifies handling version printing in all CLIs binaries.
package version

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"strings"
	"text/template"
)

var (
	printVer   bool
	printShort bool
)

// RegisterParseFlags registers and parses version related flags.
func RegisterParseFlags() {
	RegisterFlags()
	flag.Parse()
}

// RegisterFlags parses version related flags.
func RegisterFlags() {
	flag.BoolVar(&printVer, "version", false, "Prints current version.")
	flag.BoolVar(&printShort, "short-version", false, "Print just the version number.")
}

// ShouldPrintVersion returns true if version should be printed.
// Use Print function to print version information.
func ShouldPrintVersion() bool {
	return printVer || printShort
}

// versionInfoTmpl contains the go template used by Print.
// Same printing template that Prometheus has: https://github.com/prometheus/common/blob/317b7b125e8fddda956d0c9574e5f03f438ed5bc/version/info.go#L58-L65
var versionInfoTmpl = `
{{.program}}, version {{.version}} (branch: {{.branch}}, revision: {{.revision}})
  build user:       {{.buildUser}}
  build date:       {{.buildDate}}
  go version:       {{.goVersion}}
`

// Print version information to a given out writer.
func Print(out io.Writer, program string) {
	if printShort {
		fmt.Fprint(out, Version)
		return
	}

	m := map[string]string{
		"program":   program,
		"version":   Version,
		"revision":  Revision,
		"branch":    Branch,
		"buildUser": BuildUser,
		"buildDate": BuildDate,
		"goVersion": GoVersion,
	}
	t := template.Must(template.New("version").Parse(versionInfoTmpl))

	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "version", m); err != nil {
		panic(err)
	}

	fmt.Fprintln(out, strings.TrimSpace(buf.String()))
}
