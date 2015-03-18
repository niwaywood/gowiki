# gowiki
Extension of golang.org "Writing Web Applications" codewalk

This is my first Go project. After I completed "A tour of Go" I worked through the "Writing Web Applications" codewalk.

To learn more about Go I have decided to extend the codewalk solution.

The current extensions are:
- Using MongoDb (via mgo) instead of using files
- Dependency management with nut

## Setup
- Setup mongoDB
- Install nut dependency manager: `go get github.com/jingweno/nut`
- Install dependencies: `nut install` from project directory

## Usage
- Start mongo: `mongod`
- Compile and run app: `go install gowiki && gowiki`
