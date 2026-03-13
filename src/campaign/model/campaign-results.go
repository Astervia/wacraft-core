package campaign_model

import (
	"sync"

	synch_contract "github.com/Astervia/wacraft-core/src/synch/contract"
)

// CampaignResults tracks Sent / Successes / Errors across goroutines and,
// when a DistributedCounter is provided, across multiple instances.
type CampaignResults struct {
	Total     int64 `json:"total"`
	Sent      int64 `json:"sent"`
	Successes int64 `json:"successes"`
	Errors    int64 `json:"errors"`

	mu          *sync.Mutex
	counter     synch_contract.DistributedCounter // nil in memory-only mode
	campaignKey string                            // key prefix for Redis counters
}

// HandleError increments the counters and calls the progress callback.
// When a DistributedCounter is configured the remote counters are updated
// atomically via INCRBY; local fields are kept in sync for the callback.
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

	if s.counter != nil {
		s.counter.Increment("sent:"+s.campaignKey, 1)
		if err != nil {
			s.counter.Increment("errors:"+s.campaignKey, 1)
		} else {
			s.counter.Increment("successes:"+s.campaignKey, 1)
		}
	}

	callback(s)
}

// CreateCampaignResults creates an in-memory-only CampaignResults (backward compatible).
func CreateCampaignResults(total int64) *CampaignResults {
	var mu sync.Mutex
	return &CampaignResults{
		Total: total,
		mu:    &mu,
	}
}

// CreateCampaignResultsWithCounter creates CampaignResults backed by a
// DistributedCounter so counts are visible across instances.
// campaignID must be unique per campaign (used as Redis key suffix).
func CreateCampaignResultsWithCounter(
	total int64,
	counter synch_contract.DistributedCounter,
	campaignID string,
) *CampaignResults {
	var mu sync.Mutex
	return &CampaignResults{
		Total:       total,
		mu:          &mu,
		counter:     counter,
		campaignKey: campaignID,
	}
}
