package fee

import "go.mongodb.org/mongo-driver/bson/primitive"

type Log struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	MemberID  string             `json:"member_id" bson:"member_id"`
	Amount    int                `json:"amount,string" bson:"amount"`
	Type      string             `json:"type" bson:"type"`
	CreatedAt int64              `json:"created_at,string" bson:"created_at"`
	UpdatedAt int64              `json:"updated_at,string" bson:"updated_at"`
}

type Logs []Log
