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
