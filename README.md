# gowiki
Extension of golang.org's [Writing Web Applications](http://golang.org/doc/articles/wiki/) codewalk

This is my first Go project. After I completed [A Tour of Go](http://tour.golang.org/welcome/1) I worked through the [Writing Web Applications](http://golang.org/doc/articles/wiki/) codewalk.

To learn more about Go I have decided to extend the codewalk solution.

The **current** extensions are:
- Dependency management with [nut](https://github.com/jingweno/nut)
- Using [MongoDb](http://www.mongodb.org/) (via [mgo](http://labix.org/mgo)) instead of using files

Future planned extensions are:
- Middleware for logging, 404's etc
- Third party router to make url analysis simpler (i.e. parameterizing url path)

## Setup
- Setup mongoDB
- Install nut dependency manager: `go get github.com/jingweno/nut`
- Install dependencies: `nut install` from project directory

## Usage
- Start mongo: `mongod`
- Compile and run app: `go install gowiki && gowiki`
