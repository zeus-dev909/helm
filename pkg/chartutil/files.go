/*
Copyright 2016 The Kubernetes Authors All rights reserved.
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

package chartutil

import (
	"encoding/base64"
	"path"

	yaml "gopkg.in/yaml.v2"

	"github.com/gobwas/glob"
	"github.com/golang/protobuf/ptypes/any"
)

// Files is a map of files in a chart that can be accessed from a template.
type Files map[string][]byte

// NewFiles creates a new Files from chart files.
// Given an []*any.Any (the format for files in a chart.Chart), extract a map of files.
func NewFiles(from []*any.Any) Files {
	files := map[string][]byte{}
	if from != nil {
		for _, f := range from {
			files[f.TypeUrl] = f.Value
		}
	}
	return files
}

// GetBytes gets a file by path.
//
// The returned data is raw. In a template context, this is identical to calling
// {{index .Files $path}}.
//
// This is intended to be accessed from within a template, so a missed key returns
// an empty []byte.
func (f Files) GetBytes(name string) []byte {
	v, ok := f[name]
	if !ok {
		return []byte{}
	}
	return v
}

// Get returns a string representation of the given file.
//
// Fetch the contents of a file as a string. It is designed to be called in a
// template.
//
//	{{.Files.Get "foo"}}
func (f Files) Get(name string) string {
	return string(f.GetBytes(name))
}

// Glob takes a glob pattern and returns another files object only containing
// matched  files.
//
// This is designed to be called from a template.
//
// {{ range $name, $content := .Files.Glob("foo/**") }}
// {{ $name }}: |
// {{ .Files.Get($name) | indent 4 }}{{ end }}
func (f Files) Glob(pattern string) Files {
	g, err := glob.Compile(pattern, '/')
	if err != nil {
		g, _ = glob.Compile("**")
	}

	nf := NewFiles(nil)
	for name, contents := range f {
		if g.Match(name) {
			nf[name] = contents
		}
	}

	return nf
}

// AsConfig turns a Files group and flattens it to a YAML map suitable for
// including in the `data` section of a kubernetes ConfigMap definition.
// Duplicate keys will be overwritten, so be aware that your filenames
// (regardless of path) should be unique.
//
// This is designed to be called from a template, and will return empty string
// (via ToYaml function) if it cannot be serialized to YAML, or if the Files
// object is nil.
//
// The output will not be indented, so you will want to pipe this to the
// `indent` template function.
//
//   data:
// {{ .Files.Glob("config/**").AsConfig() | indent 4 }}
func (f Files) AsConfig() string {
	if f == nil {
		return ""
	}

	m := map[string]string{}

	// Explicitly convert to strings, and file names
	for k, v := range f {
		m[path.Base(k)] = string(v)
	}

	return ToYaml(m)
}

// AsSecrets returns the value of a Files object as base64 suitable for
// including in the `data` section of a kubernetes Secret definition.
// Duplicate keys will be overwritten, so be aware that your filenames
// (regardless of path) should be unique.
//
// This is designed to be called from a template, and will return empty string
// (via ToYaml function) if it cannot be serialized to YAML, or if the Files
// object is nil.
//
// The output will not be indented, so you will want to pipe this to the
// `indent` template function.
//
//   data:
// {{ .Files.Glob("secrets/*").AsSecrets() }}
func (f Files) AsSecrets() string {
	if f == nil {
		return ""
	}

	m := map[string]string{}

	for k, v := range f {
		m[path.Base(k)] = base64.StdEncoding.EncodeToString(v)
	}

	return ToYaml(m)
}

// ToYaml takes an interface, marshals it to yaml, and returns a string. It will
// always return a string, even on marshal error (empty string).
//
// This is designed to be called from a template.
func ToYaml(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return string(data)
}
