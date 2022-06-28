# Go sandbox
Sandbox for learning go, based on the book _Building Distributed Applications in Gin_.

## Starting with go

### First steps

_TBD_

### Initializing a module and adding dependencies

To initialize a module just issue the following.
```
go mod init github.com/aheadxnet/go-sandbox
```
Now we can add gin using the following command:
```
go get github.com/gin-gonic/gin
```
To list all dependencies issue the following command.
```
go list -m all
```
And if you prefer JSON format just say so.
```
go list -m -json all
```

### Vendoring support in go

If you prefer to store the dependencies together with your application in the remote repository,
Go modules support this feature.
```
go mod vendor
```
This creates a directory ``vendor`` and places all the dependencies beneath this directory so you can add them
to your Git project and push them to the remote repository.

### Showing the dependency graph of your application

There seems to be no official tool to accomplish this. and the way shown in the above book also seems not to work - 
at least I didn't manage to et it working.

One tools seems to be _Kyle Banks_ ``depth`` command.
```
go install github.com/KyleBanks/depth/cmd/depth
```
If the go packages are in your ``PATH`` (for example due to issueing ``./do-env.sh`` in this project's root),
just walk your way ti your ``main.go`` (by cd-ing to ``src/hello-world``) and type the following (pipe ``|`` it to
``less``) to easily navigate through the result).
```
depth .
```

## Providing an API for recipes

### POSTing a new dataset

Post at ``localhost:8080/recipes`` a recipe, for example.
```
{
  "name": "Stefans Käsekuchen",
  "tags": ["Kuchen", "Nachtisch", "Freiburg"],
  "ingredients": [
    "500g Sahnequark",
    "200g Sahne",
    "2 Eier"
  ],
  "instructions": [
    "Zunächst den Mürbteig zubereiten",
    "Puddingmasse zubereiten",
    "Teig in Form bringen",
    "Puddingmasse einfüllen",
    "Kuchen backen"
  ]
}
```
Alternatively you may use ``cURL`` to do the POST:
```
curl --location --request POST 'http://localhost:8080/recipes' \
--header 'Content-Type: application/json' \
--data-raw '{
   "name": "Homemade Pizza",
   "tags" : ["italian", "pizza", "dinner"],
   "ingredients": [
       "1 1/2 cups (355 ml) warm water (105°F-115°F)",
       "1 package (2 1/4 teaspoons) of active dry yeast",
       "3 3/4 cups (490 g) bread flour",
       "feta cheese, firm mozzarella cheese, grated"
   ],
   "instructions": [
       "Step 1.",
       "Step 2.",
       "Step 3."
   ]
}' | jq -r
``` 

### GETting a list of recipes

Get a list of recipes by ``http://localhost:8080/recipes``. Or using ``cURL``:
``` 
curl -s --location --request GET 'http://localhost:8080/recipes' \
--header 'Content-Type: application/json'
``` 

### PUTting a recipe to update its state
Put a recipe by ``http://localhost:8080/recipes/{id}``.

### DELETEing a recipe
Delete a recipe by ``http://localhost:8080/recipes/{id}``.
```
curl -v -sX DELETE http://localhost:8080/recipes/c0283p3d0cvuglq85log | jq -r
```

### GETting all recipes with a given tag
Searching for recipes with a given tag by ``http://localhost:8080/recipes/search?tag=mytag``.

## Documenting the API with swagger

Get the latest release of ``go-swagger`` to work with. There are binaries for several OSses, choose the one meeting 
your environment.

### For Linux-style OS:
Copy the binary to a folder of your liking and place a symlink in ``/usr/local/bin``:
```
ln -s /opt/go-swagger/current/go-swagger /usr/local/bin/swagger
```
You might have to do this as ``root`` or by ``sudo``-ing it.

To generate the documentation use 
```
swagger generate spec -o ./swagger.json
```
in the projects root directory.

And you can serve this as with
```
swagger serve ./swagger.json
```
Or in Swagger UI look
```
swagger serve -F swagger ./swagger.json
```

## Working with a MongoDB

It's cool to have data in your service, but it's even cooler to persist this with a database.
For now let's work with a MongoDB.

### Prerequisites

For a simple setup we assume that you use a local database,
for example a MongoDB in a local _docker_ container.

#### Creating and running a MongoDB container with docker

Run 

```
docker run -d --name mongodb -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=password -p 27017:27017 mongo:4.4.3
```

to create and start a new MongoDB container.

Use _docker swarm_ to hide the secrets.

```
openssl rand -base64 12 | docker secret create mongodb_password -
```

Modify start command to use

```
-e MONGO_INITDB_ROOT_PASSWORD_FILE=/run/secrets/mongodb_password
```

For further infos about ```docker```, swarm and secrets see https://docs.docker.com/engine/swarm/secrets/
or https://newbedev.com/how-to-use-docker-secrets-without-a-swarm-cluster

For accessing the DB with a GUI you may use [Compass](https://www.mongodb.com/try/download/compass?tck=docs_compass).
Configure the connection string like so:

```
mongodb://admin:s0ky84la2vp8lxmak0fqv91cy@localhost:27017/test?authSource=admin&readPreference=primary&appname=MongoDB%20Compass&directConnection=true&ssl=false
```
> If you are interestet in managing your docker containers more professionally you
> might consider using [Portainer](https://portainer.io).
 
### Update project to user mongo-db driver

Get the mongo db driver

```
go get go.mongodb.org/mongo-driver/mongo
```

And integrate it into the project, like so.

```
package main
import (
   "go.mongodb.org/mongo-driver/mongo"
   "go.mongodb.org/mongo-driver/mongo/options"
   "go.mongodb.org/mongo-driver/mongo/readpref"
)
```

Initialize the database connection in the project.

```
var ctx context.Context
var err error
var client *mongo.Client

func init() {
   ctx = context.Background()
   client, err = mongo.Connect(ctx, 
       options.Client().ApplyURI(os.Getenv("MONGO_URI")))
   if err = client.Ping(context.TODO(), 
                        readpref.Primary()); err != nil {
       log.Fatal(err)
   }
   log.Println("Connected to MongoDB")
}
```

And implement the ```ListRecipesHandler``` function.

```
func ListRecipesHandler(c *gin.Context) {
   cur, err := collection.Find(ctx, bson.M{})
   if err != nil {
       c.JSON(http.StatusInternalServerError, 
              gin.H{"error": err.Error()})
       return
   }
   defer cur.Close(ctx)
   recipes := make([]Recipe, 0)
   for cur.Next(ctx) {
       var recipe Recipe
       cur.Decode(&recipe)
       recipes = append(recipes, recipe)
   }
   c.JSON(http.StatusOK, recipes)
}
```

The other functions have to be updated likewise.

### Project layout

Up until now we stored everything in the file ```main.go```, which is OK for starting with, but which will fail on
larger projects.

So how do we structure a _go_ project?

Basically we will separate the model and the handler code.

#### Separating the ```model```

Move the ```Recipe``` struct to a new folder named ```models```.

```
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
```

#### Placing all handlers in ```handlers```

All _handlers_ go into a package named ```handlers```, though we only have one, yet.

```
package handlers

import (
	"context"
	"fmt"
	"github.com/aheadxnet/go-sandbox/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

type RecipesHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewRecipesHandler(ctx context.Context, collection *mongo.Collection) *RecipesHandler {
	return &RecipesHandler{
		collection: collection,
		ctx:        ctx,
	}
}
```

> We create a struct for the collection and the context, so we don't have those as global variables anymore.

```
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
	ctx.JSON(http.StatusCreated, recipe)
}
```
> Then we move all the handler functions from main to the newly created ```handler.go```.
> These functions now become members of the ```RecipesHandler``` class, so they have access
> to the internal state of the handler - means: they can access the collection.

#### Bringing it all together in ```main.go```

The ```main.go``` file now becomes much smaller and cleaner. There is only some initialization code, that's all.

```
package main

import (
	"context"
	"github.com/aheadxnet/go-sandbox/handlers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
)

var recipesHandler *handlers.RecipesHandler

func init() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")

	recipesHandler = handlers.NewRecipesHandler(ctx, collection)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	router.GET("/recipes/:id", recipesHandler.GetRecipeHandler)
	/*
		router.GET("/recipes/search", SearchRecipesHandler)*/
	router.Run()
}
```

Only the ```init()``` function and the ```main()``` function are kept in the file (and some imports, of course).
Doesn't this look neat?

## Performance issues? Use caching!

To show how easy it is to introduce fast caching mechanism into the project we will use [Redis](https://redis.io/).

Redis is an in-memory database, everything is stored in the main memory - so it is easily and fast accessed.

To start redis in a docker container run the following command.

```
docker run -d --name redis -p 6379:6379 redis
```

For production scenarios consider using adjusted configuration values, to be placed in a file named ```redis.conf```, for example.

```
maxmemory-policy allkeys-lru
maxmemory 512mb
```

Then you create and start the container with

```
docker run -d -v $PWD/conf:/usr/local/etc/redis --name redis -p 6379:6379 redis
```

To access the redis cache from within the application get the driver first.

```
go get github.com/go-redis/redis/v8
```

Add an import to ```main.go``` like

``` 
import redis "github.com/go-redis/redis/v8"
```

Now you can add the following initialization code in the ```ini()``` function.

``` 
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	status := redisClient.Ping(ctx)
	fmt.Println(status)
```

Try starting the application now, it should run (without using the cache yet).

> If go complains about some libs being declared "explicit" in go.mod and not in ```vendor/modules.txt```,
> try fixing this with ```go mod vendor``` as may be suggested by go itself.

Now we can add the ```redisClient``` to the ```RecipesHandler``` class.

``` 
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
```

Now we can optimize the handler function to retrieve all recipes.

``` 
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
```

If you want to check of the redis cache works as designed you can look up the ```recipes``` entry in the container.

``` 
$ docker exec -it redis bash
root@5035bb0e864d:/data# redis-cli 
127.0.0.1:6379> EXISTS recipes
(integer) 1
127.0.0.1:6379> EXISTS blubb
(integer) 0
127.0.0.1:6379> exit
root@5035bb0e864d:/data# exit
exit
```

This issued command should return 1 for an existing entry (recipes) and 0 for a non existing one (blubb).

> To get even more insight you can use [Redis Insight](https://redislabs.com/fr/redis-enterprise/redis-insight/), 
> start a docker container with with redis insight via
> ```docker run -d --name redisinsight --link redis -p 8001:8001 redislabs/redisinsight```.
> Now you can connect this (Host: redis, default port 6379, database local, no username and no password)
> to your redis db container.

Two things have to be considered with using a cache:
1. How long do you wish to read the cached value from the cache? We should use a TTL (time to live) parameter for taht scenario.
2. If we add or update a value in the Database we should invalidate or delete the cached value to prevent reading an outdated state from the cache.

So the ```NewRecipeHandler``` functionshould look something like this.

``` 
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
```

### Analyzing the cache with a benchmark

Using [Apache Benchmark](https://httpd.apache.org/docs/2.4/programs/ab.html) you can test the performance with for example 100 requests in parallel and 2000 requests total.

Without cache:

``` 
ab -n 2000 -c 100 -g without-cache.data http://localhost:8080/recipes
This is ApacheBench, Version 2.3 <$Revision: 1843412 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 200 requests
Completed 400 requests
Completed 600 requests
Completed 800 requests
Completed 1000 requests
Completed 1200 requests
Completed 1400 requests
Completed 1600 requests
Completed 1800 requests
Completed 2000 requests
Finished 2000 requests


Server Software:        
Server Hostname:        localhost
Server Port:            8080

Document Path:          /recipes
Document Length:        679602 bytes

Concurrency Level:      100
Time taken for tests:   12.159 seconds
Complete requests:      2000
Failed requests:        0
Total transferred:      1359410000 bytes
HTML transferred:       1359204000 bytes
Requests per second:    164.49 [#/sec] (mean)
Time per request:       607.929 [ms] (mean)
Time per request:       6.079 [ms] (mean, across all concurrent requests)
Transfer rate:          109186.17 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   1.4      0      14
Processing:    28  592 308.7    556    1830
Waiting:       26  580 299.6    547    1829
Total:         29  593 309.1    556    1830

Percentage of the requests served within a certain time (ms)
  50%    556
  66%    679
  75%    778
  80%    830
  90%    996
  95%   1165
  98%   1393
  99%   1441
 100%   1830 (longest request)
```

With cache:

```
ab -n 2000 -c 100 -g with-cache.data http://localhost:8080/recipes
This is ApacheBench, Version 2.3 <$Revision: 1843412 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 200 requests
Completed 400 requests
Completed 600 requests
Completed 800 requests
Completed 1000 requests
Completed 1200 requests
Completed 1400 requests
Completed 1600 requests
Completed 1800 requests
Completed 2000 requests
Finished 2000 requests


Server Software:        
Server Hostname:        localhost
Server Port:            8080

Document Path:          /recipes
Document Length:        679602 bytes

Concurrency Level:      100
Time taken for tests:   9.026 seconds
Complete requests:      2000
Failed requests:        0
Total transferred:      1359410000 bytes
HTML transferred:       1359204000 bytes
Requests per second:    221.59 [#/sec] (mean)
Time per request:       451.292 [ms] (mean)
Time per request:       4.513 [ms] (mean, across all concurrent requests)
Transfer rate:          147083.02 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.7      0       4
Processing:    22  441 304.0    390    1960
Waiting:       21  434 301.7    384    1791
Total:         22  441 304.1    390    1960

Percentage of the requests served within a certain time (ms)
  50%    390
  66%    530
  75%    632
  80%    696
  90%    857
  95%    989
  98%   1191
  99%   1360
 100%   1960 (longest request)
```

To visually compare those two tests you can use ```gnuplot```.
1. Add a file ```apache-benchmark.p``` with
```
set terminal png
set output "benchmark.png"
set title "Cache benchmark"
set size 1,0.7
set grid y
set xlabel "request"
set ylabel "response time (ms)"
plot "with-cache.data" using 9 smooth sbezier with lines title "with cache", "without-cache.data" using 9 smooth sbezier with lines title "without cache"
```
2. Create the image with ```gnuplot apache-benchmark.p```.

