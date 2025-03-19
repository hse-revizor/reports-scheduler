package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/hse-revizor/reports-scheduler/internal"
	_ "github.com/lib/pq"
)

const (
	pollInterval = 5 * time.Second
)

func main() {
	mainCfg := internal.MakeMainConfigFromENV()

	// Connect to the database
	db, err := sql.Open("postgres", mainCfg.DbConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	schedulerService := internal.MakeSchedulerService(db, mainCfg)

	plannedPolicied := make(map[string]*gocron.Scheduler)

	// Start polling the database
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for range ticker.C {
		policies, err := schedulerService.FetchCheckingPolicies()
		if err != nil {
			log.Printf("Failed to fetch checking policies: %v", err)
			continue
		}

		for _, policy := range policies {
			if policy.Policy.Type == "IntervalChecking" {
				if len(policy.Policy.Params) == 0 {
					log.Printf("No cron expression found in params for policy %s", policy.ID)
					continue
				}

				_, alreadyPlanned := plannedPolicied[policy.ID]
				if alreadyPlanned {
					continue;
				}

				policyScheduler := gocron.NewScheduler(time.UTC)
				plannedPolicied[policy.ID] = policyScheduler

				cronExpression := policy.Policy.Params[0]
				if _, err := policyScheduler.Cron(cronExpression).Do(schedulerService.SendCombination, policy.Combination); err != nil {
					log.Printf("Failed to schedule cron job for policy %s: %v", policy.ID, err)
				} else {
					log.Printf("Scheduled cron job for policy %s with expression %s", policy.ID, cronExpression)
				}

				policyScheduler.StartAsync()
			}
		}
	}

	for _, sch := range plannedPolicied {
		sch.Stop()
	}

	// Keep the program running
	select {}
}
