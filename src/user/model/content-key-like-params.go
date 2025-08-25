package user_model

type ContentKeyLikeParams struct {
	KeyName  SearchableUserColumn `params:"keyName" validate:"required,searchable_user_column"`
	LikeText string               `params:"likeText" validate:"required,min=1,max=512"`
}
