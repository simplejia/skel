package api

// Skel 用于示例
type Skel struct {
	ID int64 `json:"id" bson:"_id"`
}

// NewSkel 生成skel对象
func NewSkel() *Skel {
	return &Skel{}
}
