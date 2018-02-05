# Clean Architecture in Go

**For Russian version see readme.ru-RU.md file.**

An example of "Clean Architecture" in Go to demonstrate developing a testable
application that can be run on AppEngine with Google Cloud Storage or with 
traditional hosting and MongoDB for storage (but not limited to either).

There are a number of different application architectures that are all simlar
variations on the same theme which is to have clean separation of concerns and 
dependencies that follow the best practices of "the dependency invesion principle":

A. High-level modules should not depend on low-level modules. Both should depend on abstractions.

B. Abstractions should not depend on details. Details should depend on abstractions.

Variations of the approach include:

* [The Clean Architecture](https://blog.8thlight.com/uncle-bob/2012/08/13/the-clean-architecture.html) advocated by Robert Martin ('Uncle Bob')
* Ports & Adapters or [Hexagonal Architecture](http://alistair.cockburn.us/Hexagonal+architecture) by Alistair Cockburn
* [Onion Architecture](http://jeffreypalermo.com/blog/the-onion-architecture-part-1/) by Jeffrey Palermo

From more in-depth practical application of many of the ideas I can strongly 
recommend the excellent book [Implementing Domain-Driven Design](http://www.amazon.com/Implementing-Domain-Driven-Design-Vaughn-Vernon/dp/0321834577)
by Vaughn Vernon that goes into far greater detail.

Besides the clean codebase, the approaches also bring other advantages - significant
parts of the system can be unit tested quickly and easily without having to fire 
up the full web stack, something that is often difficult when the dependencies 
go the wrong way (if you need a database and a web-server running to make your 
tests work, you're doing it wrong).ф

I'd used it before in the world of .NET but forgot about it after moving to coding
more in Python. After switching languages again (yeah, right) to the wonderful 
world of go I came across a blog post which re-ignited my interest in it:
[Applying The Clean Architecture to Go applications](http://manuel.kiessling.net/2012/09/28/applying-the-clean-architecture-to-go-applications/)

It's a great read but I found the example a little overly-complex with too much of
the focus on relational database model parts and at the same time it was light
on some issues I wanted to resolve such as switching between different storage types
and web UI or framework (and Go has so many of those to chose from!).

I've also been looking for a way to make my application usable both standalone
and on AppEngine as well as being easier to test, so this seemed like a good opportunity
to do some experimenting and this is what I came up with as a prototype which I've
hopefully simplified to show the techniques.

## Dependency Rings

We've all heard of n-tier or layered architecture, especially if you've come 
from the world of .NET or Java and it's unfair that it gets a bad rep. Probably
because it was often implemented so poorly with the typical mistake of everything
relying on the database layer at the bottom which made software rigid, difficult
to test and closely tied to the vendor's database implementation (hardly surprising
that they promoted it so hard).

Reversing the dependencies though has a wonderful transformative effect on your 
code. Here is my interpretation of the layers or rings implemented using the Go 
language (or 'Golang' for Google).

### Domain

At the center of the dependencies is the domain. These are the business objects
or entities and should represent and encapsulate the fundamental business rules
such as "can a closed customer account create a new order?". There is usually a
single root object that represents the system and which has the factory methods 
to create other objects (which in turn may have their own methods to create others). 
This is where the domain-driven design lives.

Looking at this should give you an idea of the business model for the application
and what the system is working with. This package allows code such as unit tests 
to excercise the core parts of the app for testing to ensure that basic rules are 
enforced.

### Engine / Use-Cases

The application level rules and use-cases orchestrate the domain model and add richer
rules and logic including persistence. I prefer the term engine for this package 
because it is the engine of what the app actually does. The rules implemented at this
level should not affect the domain model rules but obviously may depend on them. 
The rules of the application also shouldn't rely on the UI or the persistence 
frameworks being used.

Why would the business rules change depending on what UI framework is the new flavour 
of the month or if we want to change from an RDBMS to MongoDB or some cloud datastore?
Those are implementation details that pull the levers of the use cases or are used by
the engine via abstract interfaces.

### Interface Adapters

These are concerned with converting data from a form that the use-cases handle to
whatever the external framework and drivers use. A use-case may expect a request 
struct with a set of parameters and return a response struct with the results. The 
public facing part of the app is more likely to expect to send requests as JSON or 
http form posts and return JSON or rendered HTML. The database may return results 
in a structure that needs to be adapted into something the rest of the app can use.

### Frameworks and Drivers

These are the ports that allow the system to talk to 'outside things' which could be
a database for persistence or a web server for the UI. None of the inner use cases 
or domain entities should know about the implementation of these layers and they may 
change over time because ... well, we used to store data in SQL, then MongoDB and 
now cloud datastores. Changing the storage should not change the application or any 
of the business rules. I tend to call these "providers" because ... well, .NET.

# Run

Within cmd subfolders ...

## App Engine

Install the AppEngine SDK for Go:

    goapp serve

## Standalone

### With MongoDB

Install and start mongodb:

    mongod --config /usr/local/etc/mongod.conf       

Copy templates from cmd\webapp\shared-templates to cmd\webapp\mgo-webapp\
And run:

    cd cmd\webapp\mgo-webapp\
    go run app.go   
    
### With BoltDB

Go get boltdb:

    go get github.com/boltdb/bolt/...

Copy templates from cmd\webapp\shared-templates to cmd\webapp\bolt-webapp\
And run:

    cd cmd\webapp\bolt-webapp\
    go run app_bolted.go
                        
    
### Console Application

Go to `cmd\cli` and compile the `cliapp.go` then try to run it will show help
how to use. Uoy can add and list added greetings.

## Run Tests
Not yet added

    ginkgo watch -cover domain
    go tool cover -html=domain/domain.coverprofile

# Implementation Notes

## Build tags

Go has build tags to control which code is included and when running on AppEngine
the `appengine` tag is automatically applied. This provides an easy way to include
or exclude code that will only work on one platform or the other. i.e. there is no
point building the appengine provider into a standalone build and some code can't
be executed on appengine classic - this provides a way to keep things separated.

## Dependency Injection

Surely it's needed for such a thing? No it isn't. While DI can be a useful tool,
very often it takes over a project and becomes an entangled part of the application
architecture masquerading as the framework. Seriously, you don't need it and it 
often comes with a huge cost in terms of complexity and runtime performance. 
Whatever a DI framework does, you can do yourself with some factories - what we
used before the world went DI crazy and thought Spring was a good idea (oh, how
we laugh about it now).

## Query spec

The Query spec provides a way to pass a query definition to the providers in a
storage agnostic way without depending on any database specific querying language.
This was attempted in .NET with Linq to mixed results - you often ended up coding
for the specifics of certain databases (usually SQL server) but in this case the
query language is much simpler and designed to be more lightweight as it only has
to provide some filtering capability for what is going to be a NoSQL database or
a SQL database being used in a non-relational way.

## Storage providers

I picked AppEngine Datastore and MongoDB because they are kind of similar in that
they are both NoSQL stores but are pretty different in how connections and state
are maintained. The MongoDB storage has the connection passed in through the 
factory setup. The Datastore has no permanent connection and uses the context
from each request.

## Enhancements

There's a lot missing. Some obvious things would be to pass request information
such as authenticated user, request IP etc... in a standard struct embedded within
each interactor request. The responses should also return errors that the web ui
adapter could use in the response.

Some unit tests
would also show how the majority of the system can be tested without having to 
fire up a web server or a database. Test storage instances can be used to test 
the engine and test engine instances can help test the web handlers.  

## What's with the imports?

Why do I separate the imports in Go? I just like it ... I divide the imports into:

* Standard packages (e.g. `strings` or `fmt`)
* Extensions to packages (e.g. `net/http` or `encoding/json`)
* 3rd party packages (e.g. `google.golang.org/appengine/datastore`)
* Application packages (e.g. `github.com/captaincodeman/clean-go/domain`)