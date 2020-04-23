package models

import "go.mongodb.org/mongo-driver/bson/primitive"

//Account Json request payload is as follows,
//{
// "id": "1",
// "firstname": "wian",
// "lastname":  "vos",
// "company":  "Red Hat",
// "status": "internal"
// }
//

type Player struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" schema:"_id"`
	FirstName string             `json:"firstname,omitempty" bson:"firstname,omitempty" schema:"firstname"`
	LastName  string             `json:"lastname,omitempty" bson:"lastname,omitempty" schema:"lastname"`
	Company   string             `json:"company,omitempty" bson:"company,omitempty" schema:"company"`
	Status    string             `json:"status,omitempty" bson:"status,omitempty" schema:"status"`
	TelNumber string             `json:"telnumber,omitempty" bson:"telnumber,omitempty" schema:"telnumber"`
	Linkedin  string             `json:"linkedin,omitempty" bson:"linkedin,omitempty" schema:"linkedin"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty" schema:"email"`
}

func PlayerDefaults() Player {
	return Player{
		Company: "Red Hat",
		Status:  "Internal",
	}
}

func (p Player) checkRequired() bool {
	if p.FirstName == "" {
		return false
	}
	if p.LastName == "" {
		return false
	}
	if p.Company == "" {
		return false
	}
	if p.Status == "" {
		return false
	}
	return true
}
