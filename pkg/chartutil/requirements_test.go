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
	"sort"
	"testing"

	"strconv"

	"k8s.io/helm/pkg/proto/hapi/chart"
)

func TestLoadRequirements(t *testing.T) {
	c, err := Load("testdata/frobnitz")
	if err != nil {
		t.Fatalf("Failed to load testdata: %s", err)
	}
	verifyRequirements(t, c)
}

func TestLoadRequirementsLock(t *testing.T) {
	c, err := Load("testdata/frobnitz")
	if err != nil {
		t.Fatalf("Failed to load testdata: %s", err)
	}
	verifyRequirementsLock(t, c)
}
func TestRequirementsTagsNonValue(t *testing.T) {
	c, err := Load("testdata/subpop")
	if err != nil {
		t.Fatalf("Failed to load testdata: %s", err)
	}
	// tags with no effect
	v := &chart.Config{Raw: "tags:\n  nothinguseful: false\n\n"}
	// expected charts including duplicates in alphanumeric order
	e := []string{"parentchart", "subchart1", "subcharta", "subchartb"}

	verifyRequirementsEnabled(t, c, v, e)
}
func TestRequirementsTagsDisabledL1(t *testing.T) {
	c, err := Load("testdata/subpop")
	if err != nil {
		t.Fatalf("Failed to load testdata: %s", err)
	}
	// tags disabling a group
	v := &chart.Config{Raw: "tags:\n  front-end: false\n\n"}
	// expected charts including duplicates in alphanumeric order
	e := []string{"parentchart"}

	verifyRequirementsEnabled(t, c, v, e)
}
func TestRequirementsTagsEnabledL1(t *testing.T) {
	c, err := Load("testdata/subpop")
	if err != nil {
		t.Fatalf("Failed to load testdata: %s", err)
	}
	// tags disabling a group and enabling a different group
	v := &chart.Config{Raw: "tags:\n  front-end: false\n\n  back-end: true\n"}
	// expected charts including duplicates in alphanumeric order
	e := []string{"parentchart", "subchart2", "subchartb", "subchartc"}

	verifyRequirementsEnabled(t, c, v, e)
}
func TestRequirementsTagsDisabledL2(t *testing.T) {
	c, err := Load("testdata/subpop")
	if err != nil {
		t.Fatalf("Failed to load testdata: %s", err)
	}
	// tags disabling only children
	v := &chart.Config{Raw: "tags:\n  subcharta: false\n\n  subchartb: false\n"}
	// expected charts including duplicates in alphanumeric order
	e := []string{"parentchart", "subchart1"}

	verifyRequirementsEnabled(t, c, v, e)
}
func TestRequirementsTagsDisabledL1Mixed(t *testing.T) {
	c, err := Load("testdata/subpop")
	if err != nil {
		t.Fatalf("Failed to load testdata: %s", err)
	}
	// tags disabling all parents/children with additional tag re-enabling a parent
	v := &chart.Config{Raw: "tags:\n  front-end: false\n\n  subchart1: true\n\n  back-end: false\n"}
	// expected charts including duplicates in alphanumeric order
	e := []string{"parentchart", "subchart1"}

	verifyRequirementsEnabled(t, c, v, e)
}
func TestRequirementsConditionsNonValue(t *testing.T) {
	c, err := Load("testdata/subpop")
	if err != nil {
		t.Fatalf("Failed to load testdata: %s", err)
	}
	// tags with no effect
	v := &chart.Config{Raw: "subchart1:\n  nothinguseful: false\n\n"}
	// expected charts including duplicates in alphanumeric order
	e := []string{"parentchart", "subchart1", "subcharta", "subchartb"}

	verifyRequirementsEnabled(t, c, v, e)
}
func TestRequirementsConditionsEnabledL1Both(t *testing.T) {
	c, err := Load("testdata/subpop")
	if err != nil {
		t.Fatalf("Failed to load testdata: %s", err)
	}
	// conditions enabling the parent charts, effectively enabling children
	v := &chart.Config{Raw: "subchart1:\n  enabled: true\nsubchart2:\n  enabled: true\n"}
	// expected charts including duplicates in alphanumeric order
	e := []string{"parentchart", "subchart1", "subchart2", "subcharta", "subchartb", "subchartb", "subchartc"}

	verifyRequirementsEnabled(t, c, v, e)
}
func TestRequirementsConditionsDisabledL1Both(t *testing.T) {
	c, err := Load("testdata/subpop")
	if err != nil {
		t.Fatalf("Failed to load testdata: %s", err)
	}
	// conditions disabling the parent charts, effectively disabling children
	v := &chart.Config{Raw: "subchart1:\n  enabled: false\nsubchart2:\n  enabled: false\n"}
	// expected charts including duplicates in alphanumeric order
	e := []string{"parentchart"}

	verifyRequirementsEnabled(t, c, v, e)
}

func TestRequirementsConditionsSecond(t *testing.T) {
	c, err := Load("testdata/subpop")
	if err != nil {
		t.Fatalf("Failed to load testdata: %s", err)
	}
	// conditions a child using the second condition path of child's condition
	v := &chart.Config{Raw: "subchart1:\n  subcharta:\n    enabled: false\n"}
	// expected charts including duplicates in alphanumeric order
	e := []string{"parentchart", "subchart1", "subchartb"}

	verifyRequirementsEnabled(t, c, v, e)
}
func TestRequirementsCombinedDisabledL2(t *testing.T) {
	c, err := Load("testdata/subpop")
	if err != nil {
		t.Fatalf("Failed to load testdata: %s", err)
	}
	// tags enabling a parent/child group with condition disabling one child
	v := &chart.Config{Raw: "subchartc:\n  enabled: false\ntags:\n  back-end: true\n"}
	// expected charts including duplicates in alphanumeric order
	e := []string{"parentchart", "subchart1", "subchart2", "subcharta", "subchartb", "subchartb"}

	verifyRequirementsEnabled(t, c, v, e)
}
func TestRequirementsCombinedDisabledL1(t *testing.T) {
	c, err := Load("testdata/subpop")
	if err != nil {
		t.Fatalf("Failed to load testdata: %s", err)
	}
	// tags will not enable a child if parent is explicitly disabled with condition
	v := &chart.Config{Raw: "subchart1:\n  enabled: false\ntags:\n  front-end: true\n"}
	// expected charts including duplicates in alphanumeric order
	e := []string{"parentchart"}

	verifyRequirementsEnabled(t, c, v, e)
}

func verifyRequirementsEnabled(t *testing.T, c *chart.Chart, v *chart.Config, e []string) {
	out := []*chart.Chart{}
	err := ProcessRequirementsEnabled(c, v)
	if err != nil {
		t.Errorf("Error processing enabled requirements %v", err)
	}
	out = extractCharts(c, out)
	// build list of chart names
	p := []string{}
	for _, r := range out {
		p = append(p, r.Metadata.Name)
	}
	//sort alphanumeric and compare to expectations
	sort.Strings(p)
	if len(p) != len(e) {
		t.Errorf("Error slice lengths do not match got %v, expected %v", len(p), len(e))
		return
	}
	for i := range p {
		if p[i] != e[i] {
			t.Errorf("Error slice values do not match got %v, expected %v", p[i], e[i])
		}
	}
}

// extractCharts recursively searches chart dependencies returning all charts found
func extractCharts(c *chart.Chart, out []*chart.Chart) []*chart.Chart {

	if len(c.Metadata.Name) > 0 {
		out = append(out, c)
	}
	for _, d := range c.Dependencies {
		out = extractCharts(d, out)
	}
	return out
}
func TestProcessRequirementsImportValues(t *testing.T) {
	c, err := Load("testdata/subpop")
	if err != nil {
		t.Fatalf("Failed to load testdata: %s", err)
	}

	v := &chart.Config{Raw: ""}

	e := make(map[string]string)

	e["imported-chart1.SC1bool"] = "true"
	e["imported-chart1.SC1float"] = "3.14"
	e["imported-chart1.SC1int"] = "100"
	e["imported-chart1.SC1string"] = "dollywood"
	e["imported-chart1.SC1extra1"] = "11"
	e["imported-chart1.SPextra1"] = "helm rocks"
	e["imported-chart1.SC1extra1"] = "11"

	e["imported-chartA.SCAbool"] = "false"
	e["imported-chartA.SCAfloat"] = "3.1"
	e["imported-chartA.SCAint"] = "55"
	e["imported-chartA.SCAstring"] = "jabba"
	e["imported-chartA.SPextra3"] = "1.337"
	e["imported-chartA.SC1extra2"] = "1.337"
	e["imported-chartA.SCAnested1.SCAnested2"] = "true"

	e["imported-chartA-B.SCAbool"] = "false"
	e["imported-chartA-B.SCAfloat"] = "3.1"
	e["imported-chartA-B.SCAint"] = "55"
	e["imported-chartA-B.SCAstring"] = "jabba"

	e["imported-chartA-B.SCBbool"] = "true"
	e["imported-chartA-B.SCBfloat"] = "7.77"
	e["imported-chartA-B.SCBint"] = "33"
	e["imported-chartA-B.SCBstring"] = "boba"
	e["imported-chartA-B.SPextra5"] = "k8s"
	e["imported-chartA-B.SC1extra5"] = "tiller"

	e["overridden-chart1.SC1bool"] = "false"
	e["overridden-chart1.SC1float"] = "3.141592"
	e["overridden-chart1.SC1int"] = "99"
	e["overridden-chart1.SC1string"] = "pollywog"
	e["overridden-chart1.SPextra2"] = "42"

	e["overridden-chartA.SCAbool"] = "true"
	e["overridden-chartA.SCAfloat"] = "41.3"
	e["overridden-chartA.SCAint"] = "808"
	e["overridden-chartA.SCAstring"] = "jaberwocky"
	e["overridden-chartA.SPextra4"] = "true"

	e["overridden-chartA-B.SCAbool"] = "true"
	e["overridden-chartA-B.SCAfloat"] = "41.3"
	e["overridden-chartA-B.SCAint"] = "808"
	e["overridden-chartA-B.SCAstring"] = "jaberwocky"
	e["overridden-chartA-B.SCBbool"] = "false"
	e["overridden-chartA-B.SCBfloat"] = "1.99"
	e["overridden-chartA-B.SCBint"] = "77"
	e["overridden-chartA-B.SCBstring"] = "jango"
	e["overridden-chartA-B.SPextra6"] = "111"
	e["overridden-chartA-B.SCAextra1"] = "23"
	e["overridden-chartA-B.SCBextra1"] = "13"
	e["overridden-chartA-B.SC1extra6"] = "77"

	// `exports` style
	e["SCBexported1B"] = "1965"
	e["SC1extra7"] = "true"
	e["SCBexported2A"] = "blaster"
	e["global.SC1exported2.all.SC1exported3"] = "SC1expstr"

	verifyRequirementsImportValues(t, c, v, e)
}
func verifyRequirementsImportValues(t *testing.T, c *chart.Chart, v *chart.Config, e map[string]string) {

	err := ProcessRequirementsImportValues(c, v)
	if err != nil {
		t.Errorf("Error processing import values requirements %v", err)
	}
	cv := c.GetValues()
	cc, err := ReadValues([]byte(cv.Raw))
	if err != nil {
		t.Errorf("Error reading import values %v", err)
	}
	for kk, vv := range e {
		pv, err := cc.PathValue(kk)
		if err != nil {
			t.Fatalf("Error retrieving import values table %v %v", kk, err)
			return
		}

		switch pv.(type) {
		case float64:
			s := strconv.FormatFloat(pv.(float64), 'f', -1, 64)
			if s != vv {
				t.Errorf("Failed to match imported float value %v with expected %v", s, vv)
				return
			}
		case bool:
			b := strconv.FormatBool(pv.(bool))
			if b != vv {
				t.Errorf("Failed to match imported bool value %v with expected %v", b, vv)
				return
			}
		default:
			if pv.(string) != vv {
				t.Errorf("Failed to match imported string value %v with expected %v", pv, vv)
				return
			}
		}

	}
}
