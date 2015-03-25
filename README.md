# gowiki
Extension of golang.org's [Writing Web Applications](http://golang.org/doc/articles/wiki/) codewalk

This is my first Go project. After I completed [A Tour of Go](http://tour.golang.org/welcome/1) I worked through the [Writing Web Applications](http://golang.org/doc/articles/wiki/) codewalk.

To learn more about Go I have decided to extend the codewalk solution.

The **current** extensions are:
- Dependency management with [nut](https://github.com/jingweno/nut)
- Using [MongoDb](http://www.mongodb.org/) (via [mgo](http://labix.org/mgo)) instead of using files
- Using [mux](http://www.gorillatoolkit.org/pkg/mux) router to parameterise url and enforce http method types
- Using [negroni](https://github.com/codegangsta/negroni) to mount middleware

Future planned extensions are:
- Use nested routers with mux so that negroni can mount validateURL middleware only on edit, save and viewHandler's
- Middleware for logging, 404's etc
- make home page display list of wiki pages

## Setup
- Setup mongoDB
- Install nut dependency manager: `go get github.com/jingweno/nut`
- Install dependencies: `nut install` from project directory

## Usage
- Start mongo: `mongod`
- Compile and run app: `go install gowiki && gowiki`
