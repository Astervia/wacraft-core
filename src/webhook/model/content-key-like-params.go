package webhook_model

type ContentKeyLikeParams struct {
	KeyName  SearchableWebhookColumn `params:"keyName" validate:"required,searchable_webhook_column"`
	LikeText string                  `params:"likeText" validate:"required,min=1,max=512"`
}
