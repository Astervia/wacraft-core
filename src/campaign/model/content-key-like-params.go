package campaign_model

type ContentKeyLikeParams struct {
	KeyName  SearchableCampaignColumn `params:"keyName" validate:"required,searchable_campaign_column"`
	LikeText string                   `params:"likeText" validate:"required,min=1,max=512"`
}
