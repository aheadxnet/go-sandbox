package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// swagger:model Recipe
// A recipe used in this application.
type Recipe struct {
	// the id for this recipe
	//
	// required: true
	// min: 1
	ID primitive.ObjectID `json:"id" bson:"_id"`

	// the name for this recipe
	// required: true
	// min length: 3
	Name string `json:"name" bson:"name"`

	// tags for this recipe
	Tags []string `json:"tags" bson:"tags"`

	// ingredients for this recipe
	Ingredients []string `json:"ingredients" bson:"ingredients"`

	// instructions for preparing this recipe
	Instructions []string `json:"instructions" bson:"instructions"`

	// the publication date for this recipe
	// required: true
	PublishedAt time.Time `json:"publishedAt" bson:"publishedAt"`
}
