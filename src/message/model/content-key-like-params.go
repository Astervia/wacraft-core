package message_model

type ContentKeyLikeParams struct {
	KeyName  JsonMessageKey `params:"keyName" validate:"required,json_message_key"`
	LikeText string         `params:"likeText" validate:"required,min=1,max=512"`
}
