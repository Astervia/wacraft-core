package campaign_model

type SearchableCampaignColumn string

var Name SearchableCampaignColumn = "name"

func (t SearchableCampaignColumn) IsValid() bool {
	switch t {
	case Name:
		return true
	default:
		return false
	}
}
