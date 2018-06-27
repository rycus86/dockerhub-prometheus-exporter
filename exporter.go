package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strconv"
	"time"
)

type RepositoryList struct {
	Count int64
	Next  string

	Results []struct {
		Name      string
		Namespace string
		StarCount int64 `json:"star_count"`
		PullCount int64 `json:"pull_count"`
	}
}

func collectStats(client *http.Client) {
	for _, owner := range owners {
		log.Println("Collecting metrics for", owner)

		url := fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/?page_size=100", owner)

		for url != "" {
			var repos RepositoryList

			if response, err := client.Get(url); err != nil {
				log.Println("Failed to list repositories for ", owner)
				break

			} else if err := json.NewDecoder(response.Body).Decode(&repos); err != nil {
				response.Body.Close()
				log.Println("Failed to list repositories for ", owner)
				break

			} else {
				response.Body.Close()

				url = repos.Next

				repoCount.WithLabelValues(owner).Set(float64(repos.Count))

				for _, repo := range repos.Results {
					starCount.WithLabelValues(repo.Namespace, repo.Name).Set(float64(repo.StarCount))
					pullCount.WithLabelValues(repo.Namespace, repo.Name).Set(float64(repo.PullCount))
				}
			}
		}
	}
}

func main() {
	if len(owners) == 0 {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		fmt.Println()

		log.Fatal("No owners were defined")
	}

	client := &http.Client{Timeout: *timeout}

	go func() {
		firstRun := time.After(0 * time.Second)

		for {
			select {
			case <-firstRun:
				collectStats(client)

			case <-time.Tick(*interval):
				collectStats(client)
			}
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
