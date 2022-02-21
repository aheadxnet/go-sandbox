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