package activity

import (
	"context"
	"errors"
	"time"

	"github.com/kmu-kcc/buddy-backend/config"
	"github.com/kmu-kcc/buddy-backend/pkg/member"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Activity struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Start        int64              `json:"start,string" bson:"start"`
	End          int64              `json:"end,string" bson:"end"`
	Place        string             `json:"place" bson:"place"`
	Type         string             `json:"type" bson:"type"`
	Description  string             `json:"description" bson:"description"`
	Participants []string           `json:"participants" bson:"participants"`
	Applicants   []string           `json:"applicants" bson:"applicants"`
	Cancelers    []string           `json:"cancelers" bson:"cancelers"`
	Private      bool               `json:"private" bson:"private"`
}

func New(start, end int64, place, typ, description string, participants []string, private bool) Activity {
	return Activity{
		ID:           primitive.NewObjectID(),
		Start:        start,
		End:          end,
		Place:        place,
		Type:         typ,
		Description:  description,
		Participants: participants,
		Applicants:   []string{},
		Cancelers:    []string{},
		Private:      private,
	}
}

// ApplyP applies for an activity of activityID.
//
// NOTE:
//
// It is member-limited operation:
//	Only the authenticated members can access to this operation.
func ApplyP(activityID primitive.ObjectID, memberID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	collection := client.Database("club").Collection("activities")
	activity := new(Activity)

	if err = collection.FindOne(ctx, bson.M{"_id": activityID}).Decode(activity); err != nil {
		return err
	}

	if func() bool {
		for _, p := range activity.Participants {
			if memberID == p {
				return true
			}
		}
		return false
	}() {
		if err = client.Disconnect(ctx); err != nil {
			return err
		}
		return errors.New("already in participants")
	}

	if func() bool {
		for _, a := range activity.Applicants {
			if memberID == a {
				return true
			}
		}
		return false
	}() {
		if err = client.Disconnect(ctx); err != nil {
			return err
		}
		return errors.New("already in applicants")
	}

	if _, err = collection.UpdateByID(ctx, activityID, bson.M{"$push": bson.M{"applicants": memberID}}); err != nil {
		return err
	}
	return client.Disconnect(ctx)
}

// Papplies returns the applicant list of the activity of activityID.
//
// NOTE:
//
// It is privileged operation:
//	Only the club managers can access to this operation.
func Papplies(activityID primitive.ObjectID) (members []member.Member, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}
	activity := new(Activity)
	member := new(member.Member)

	if err = client.Database("club").Collection("activities").FindOne(ctx, bson.M{"_id": activityID}).Decode(activity); err != nil {
		return
	}

	filter := func() bson.D {
		arr := make(bson.A, len(activity.Applicants))
		for idx, applicant := range activity.Applicants {
			arr[idx] = applicant
		}
		return bson.D{bson.E{Key: "id", Value: bson.D{bson.E{Key: "$in", Value: arr}}}}
	}()

	cur, err := client.Database("club").Collection("members").Find(ctx, filter)
	if err != nil {
		return
	}

	for cur.Next(ctx) {
		if err = cur.Decode(member); err != nil {
			return
		}
		members = append(members, *member)
	}
	return members, client.Disconnect(ctx)
}

// ApproveP approve the applicants lis of the activity of activityID.
//
// NOTE:
//
// It is privileged operation:
//	Only the club managers can access to this operation.
func ApproveP(activityID primitive.ObjectID, ids []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	update := func() bson.D {
		arr := make(bson.A, len(ids))
		for idx, id := range ids {
			arr[idx] = id
		}
		return bson.D{
			bson.E{
				Key: "$pull",
				Value: bson.D{
					bson.E{
						Key: "applicants",
						Value: bson.D{
							bson.E{
								Key:   "$in",
								Value: arr,
							},
						},
					},
				},
			},
			bson.E{
				Key: "$push",
				Value: bson.D{
					bson.E{
						Key: "participants",
						Value: bson.D{
							bson.E{
								Key:   "$each",
								Value: arr},
						},
					},
				},
			},
		}
	}()

	if _, err := client.Database("club").
		Collection("activities").
		UpdateByID(ctx, activityID, update); err != nil {
		return err
	}
	return client.Disconnect(ctx)
}

// RejectP reject the applicanst list of the activity of activityID.
//
// NOTE:
//
// It is privileged operation:
//	Only the club managers can access to this operation.
func RejectP(activityID primitive.ObjectID, ids []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	update := func() bson.D {
		arr := make(bson.A, len(ids))
		for idx, id := range ids {
			arr[idx] = id
		}

		return bson.D{
			bson.E{
				Key: "$pull",
				Value: bson.D{
					bson.E{
						Key: "applicants",
						Value: bson.D{
							bson.E{
								Key:   "$in",
								Value: arr,
							},
						},
					},
				},
			},
		}
	}()

	if _, err = client.Database("club").Collection("activities").UpdateByID(ctx, activityID, update); err != nil {
		return err
	}
	return client.Disconnect(ctx)
}

// CancelP cancels the member of memberID's apply request to the activity of activityID.
//
// NOTE:
//
// It is member-limited operation:
//	Only the authenticated members can access to this operation.s
func CancelP(activityID primitive.ObjectID, memberID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	if _, err = client.Database("club").
		Collection("activities").
		UpdateByID(ctx, activityID, bson.M{"$pull": bson.M{"applicants": memberID}}); err != nil {
		return err
	}
	return client.Disconnect(ctx)
}
