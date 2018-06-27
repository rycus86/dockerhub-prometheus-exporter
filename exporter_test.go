package main

import (
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"gopkg.in/jarcoal/httpmock.v1"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestCollectStatsForOwner(t *testing.T) {
	repoCount.Reset()
	starCount.Reset()
	pullCount.Reset()

	owners = multiVar([]string{"rycus86"})

	data, err := ioutil.ReadFile("testdata/repos_p1.json")
	if err != nil {
		t.Fatal(err)
	}
	pageOneResponse := httpmock.NewBytesResponse(200, data)

	data, err = ioutil.ReadFile("testdata/repos_p2.json")
	if err != nil {
		t.Fatal(err)
	}
	pageTwoResponse := httpmock.NewBytesResponse(200, data)

	httpmock.Activate()
	defer httpmock.Deactivate()

	httpmock.RegisterResponder(
		"GET", "https://hub.docker.com/v2/repositories/rycus86/",
		func(req *http.Request) (*http.Response, error) {
			if req.URL.Query().Get("page") == "2" {
				return pageTwoResponse, nil
			} else {
				return pageOneResponse, nil
			}
		})

	collectStats(http.DefaultClient)

	gathered, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]int{}

	for _, g := range gathered {
		name := g.GetName()

		if !strings.HasPrefix(name, "dockerhub_") {
			continue
		}

		for _, m := range g.GetMetric() {
			labels := m.GetLabel()
			value := m.GetGauge().GetValue()

			if name == "dockerhub_repo_count" {
				if value != 38 {
					t.Error("Unexpected value:", name, m.String())
				}

				tests["count"] = 1
			}

			if labelMatches(labels, "repository", "prometheus") {
				if name == "dockerhub_pull_count" && value != 339913.0 {
					t.Error("Unexpected value:", name, m.String())
				}

				if name == "dockerhub_star_count" && value != 2.0 {
					t.Error("Unexpected value:", name, m.String())
				}

				tests["prometheus"] = 1
			}

			if labelMatches(labels, "repository", "podlike") {
				if name == "dockerhub_pull_count" && value != 2379.0 {
					t.Error("Unexpected value:", name, m.String())
				}

				if name == "dockerhub_star_count" && value != 0.0 {
					t.Error("Unexpected value:", name, m.String())
				}

				tests["podlike"] = 1
			}
		}
	}

	if len(tests) != 3 {
		t.Error("Only checked", len(tests), "metrics, but expected 3")
	}
}

func labelMatches(labels []*dto.LabelPair, name, value string) bool {
	for _, label := range labels {
		if label.GetName() == name {
			return label.GetValue() == value
		}
	}

	return false
}
