package campaign_model

import "sync"

type CampaignResults struct {
	Total     int64 `json:"total"`
	Sent      int64 `json:"sent"`
	Successes int64 `json:"successes"`
	Errors    int64 `json:"errors"`

	mu *sync.Mutex
}

func (s *CampaignResults) HandleError(
	err error,
	callback func(*CampaignResults),
) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Sent++
	if err != nil {
		s.Errors++
	} else {
		s.Successes++
	}

	callback(s)
}

func CreateCampaignResults(total int64) *CampaignResults {
	var mu sync.Mutex
	return &CampaignResults{
		Total:     total,
		Sent:      0,
		Successes: 0,
		Errors:    0,
		mu:        &mu,
	}
}
