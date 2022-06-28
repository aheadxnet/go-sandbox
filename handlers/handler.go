package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aheadxnet/go-sandbox/models"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

type RecipesHandler struct {
	collection  *mongo.Collection
	ctx         context.Context
	redisClient *redis.Client
}

func NewRecipesHandler(ctx context.Context, collection *mongo.Collection, redisClient *redis.Client) *RecipesHandler {
	return &RecipesHandler{
		collection:  collection,
		ctx:         ctx,
		redisClient: redisClient,
	}
}

// swagger:operation POST /recipes recipes newRecipe
// Create a new recipe
// ---
// parameters:
// - in: body
//   description: data for the new recipe
//   required: true
//   type: Recipe
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//     '400':
//         description: Invalid input
func (handler *RecipesHandler) NewRecipeHandler(ctx *gin.Context) {
	var recipe models.Recipe
	if err := ctx.ShouldBindJSON(&recipe); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	_, err := handler.collection.InsertOne(ctx, recipe)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new recipe"})
		return
	}
	log.Println("Remove data from Redis")
	handler.redisClient.Del(ctx, "recipes")
	ctx.JSON(http.StatusCreated, recipe)
}

// swagger:operation GET /recipes recipes listRecipes
// Returns list of recipes
// ---
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//         schema:
//           type: array
//           items: Recipe
func (handler *RecipesHandler) ListRecipesHandler(ctx *gin.Context) {
	val, err := handler.redisClient.Get(ctx, "recipes").Result()
	if err == redis.Nil {
		log.Printf("Request to MongoDB")
		cur, err := handler.collection.Find(handler.ctx,
			bson.M{})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,
				gin.H{"error": err.Error()})
			return
		}
		defer cur.Close(handler.ctx)
		recipes := make([]models.Recipe, 0)
		for cur.Next(handler.ctx) {
			var recipe models.Recipe
			cur.Decode(&recipe)
			recipes = append(recipes, recipe)
		}
		data, _ := json.Marshal(recipes)
		handler.redisClient.Set(ctx, "recipes", string(data), 0)
		ctx.JSON(http.StatusOK, recipes)
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	} else {
		log.Printf("Request to Redis")
		recipes := make([]models.Recipe, 0)
		json.Unmarshal([]byte(val), &recipes)
		ctx.JSON(http.StatusOK, recipes)
	}
}

// swagger:operation GET /recipes/{id} recipes getRecipe
// Get an existing recipe
// ---
// parameters:
// - name: id
//   in: path
//   description: ID of the recipe
//   required: true
//   type: string
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//     '404':
//         description: Invalid recipe ID
func (handler *RecipesHandler) GetRecipeHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	cur := handler.collection.FindOne(ctx, bson.M{
		"_id": objectId,
	})
	var recipe models.Recipe
	err := cur.Decode(&recipe)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, recipe)
}

// swagger:operation PUT /recipes/{id} recipes updateRecipe
// Update an existing recipe
// ---
// parameters:
// - name: id
//   in: path
//   description: ID of the recipe
//   required: true
//   type: string
// - in: body
//   description: new date of the recipe
//   required: true
//   type: Recipe
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//     '400':
//         description: Invalid input
//     '404':
//         description: Invalid recipe ID
func (handler *RecipesHandler) UpdateRecipeHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	var recipe models.Recipe
	if err := ctx.ShouldBindJSON(&recipe); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.collection.UpdateOne(ctx, bson.M{
		"_id": objectId,
	}, bson.D{{"$set", bson.D{
		{"name", recipe.Name},
		{"instructions", recipe.Instructions},
		{"ingredients", recipe.Ingredients},
		{"tags", recipe.Tags},
	}}})
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("Remove data from Redis")
	handler.redisClient.Del(ctx, "recipes")
	ctx.JSON(http.StatusOK, recipe)
}

// swagger:operation DELETE /recipes/{id} recipes deleteRecipe
// Deletes an existing recipe
// ---
// parameters:
// - name: id
//   in: path
//   description: ID of the recipe
//   required: true
//   type: string
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//     '404':
//         description: Invalid recipe ID
func (handler *RecipesHandler) DeleteRecipeHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.collection.DeleteOne(ctx, bson.M{
		"_id": objectId,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Recipe has been deleted"})
}

// swagger:operation GET /recipes/search recipes findRecipe
// Search recipes based on tags
// ---
// produces:
// - application/json
// parameters:
//   - name: tag
//     in: query
//     description: recipe tag
//     required: true
//     type: string
// responses:
//     '200':
//         description: Successful operation
func (handler *RecipesHandler) SearchRecipesHandler(c *gin.Context) {
	/*
		tag := c.Query("tag")
		listOfRecipes := make([]Recipe, 0)
		for i := 0; i < len(oldRecipes); i++ {
			found := false
			for _, t := range oldRecipes[i].Tags {
				if strings.EqualFold(t, tag) {
					found = true
				}
			}
			if found {
				listOfRecipes = append(listOfRecipes,
					oldRecipes[i])
			}
		}
		c.JSON(http.StatusOK, listOfRecipes)
	*/
}
