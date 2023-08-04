package model

// 对话信息结构体
type Conversation struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	CreatorID string   `json:"creator_id"`
	Members   []string `json:"members"`
}

// 创建对话请求参数结构体
type CreateConversationRequest struct {
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

// 添加成员请求参数结构体
type AddMemberRequest struct {
	ConversationID string   `json:"conversation_id"`
	Members        []string `json:"members"`
}
