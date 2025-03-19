package internal

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type SchedulerService struct {
	db *sql.DB
	cfg *MainConfig
}

func MakeSchedulerService(db *sql.DB, cfg *MainConfig) *SchedulerService {
	return &SchedulerService{
		db,
		cfg,
	}
}

// fetchCheckingPolicies retrieves all checking policies from the database
func (s *SchedulerService) FetchCheckingPolicies() ([]CheckingPolicy, error) {
	rows, err := s.db.Query(`SELECT id, combination, "policy_data" FROM public.checking_policies`)
	if err != nil {
		return nil, fmt.Errorf("failed to query checking policies: %w", err)
	}
	defer rows.Close()

	var policies []CheckingPolicy
	for rows.Next() {
		var policy CheckingPolicy
		if err := rows.Scan(&policy.ID, &policy.Combination, &policy.Policy); err != nil {
			return nil, fmt.Errorf("failed to scan policy: %w", err)
		}
		policies = append(policies, policy)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return policies, nil
}

// sendCombination sends the combination to the X service endpoint
func (s *SchedulerService) SendCombination(combination json.RawMessage) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", s.cfg.AnalysisServiceURL, bytes.NewReader(combination))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send combination: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Received non-OK response: %s", resp.Status)
		return
	}

	log.Printf("Analysis successfully compounded")
}
