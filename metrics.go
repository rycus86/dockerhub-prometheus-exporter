package main

import "github.com/prometheus/client_golang/prometheus"

var (
	starCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "dockerhub",
		Name:      "star_count",
		Help:      "Number of Stars",
	}, []string{"owner", "repository"})

	pullCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "dockerhub",
		Name:      "pull_count",
		Help:      "Number of Pulls",
	}, []string{"owner", "repository"})

	repoCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "dockerhub",
		Name:      "repo_count",
		Help:      "Number of Repositories",
	}, []string{"owner"})
)

func init() {
	prometheus.MustRegister(starCount, pullCount, repoCount)
}
