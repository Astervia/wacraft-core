package status_model

type ContentKeyLikeParams struct {
	KeyName  SearchableStatusColumn `params:"keyName" validate:"required,searchable_status_column"`
	LikeText string                 `params:"likeText" validate:"required,min=1,max=512"`
}
