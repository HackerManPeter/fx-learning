# Learnings.MD

## FastHTTP

I want to read up about fast http because I want to use GoFiber for the fx tutorial instead of net/http.  
GoFiber is built on top of FastHTTP, so I decided to take a look at fastHttp and how it differs from net/http.
Here are my Learnings

### Different APIs

FastHTTP offers a different set of APIs that net/http doesn't which is I think is much simpler for people coming from a ts background.  
It removes the writer and reader implementation and just provides you with a simple ctx that you can use for most, if not all,  
of your request response handling

### ServeMux

A ServeMux is a `Server Multiplexer` I believe, and what this is like a router sort of.  
It takes in a request, and determines what is the best endpoint to route that request through, something like a controller in NestJs.  
FastHTTP does not provider this and the developer would have to implement it themselves, net/http on the other hand implements this

## Building a Test Application

Using the [NestJs](https://nestjs.com) structure as a **guide** has really made building this project much easier. I say guide because I did not follow the  
format strictly, but I recognized some patterns that I found to be helpful.

I tried to do what felt intuitive, by grouping actions according to service. I started of course with the config service.  
Which has the responsibility of reading from environment variables. For my approach, I had a `NewConfigService()` function that returned  
a `*Config` struct.

I did something similar with the AuthService, I used a public `NewAuthService()` function to return an `*AuthService` struct.  
I then had a `NewAuthRoutes(app fiber.Router)` that I used to define the routes for the AuthService service. I then defined  
the handlers as receiver functions(methods) of the authservice => `func (*a AuthService) Login()`.

A major question I have been asking myself is why did I not define the routes as regular functions, and this is the double edged sword  
of Go's flexibility. Inexperienced developers shoot themselves in the foot a lot, but also learn a lot in the process.  
It seems that defining the handlers as a function that has the `*AuthService` as a parameter would work just as well, but here are  
some reasons I feel that this approach is better:

- This is where I default to my experience with [NestJs](https://nestjs.com), very rarely would will I see method that has it's arguments  
  defined like so `private validate(authservice: AuthService, loggerService: LoggerService, cacheService, databaseService DatabaseService)`.  
  I don't see this pattern a lot, and rightfully so as this is not very elegant. I could, of course, have the dependencies in the Authservice  
  and just pass the AuthService as argument. But again, I default to the patterns I learnt from NestJs
- Because I am using the [fx](https://github.com/uber-go/fx) package as a dependency injection system for this project, I would loose  
  a lot of the benefits that come with it if I choose to pass dependencies as arguments

How I chose to handle repositories is quite different from the NestJs pattern. I choose not to have dedicated Repositories for each service  
rather I pass the database connection to all services that might need it. This, I believe would reduce the chances of circular dependencies,  
because there definitely exists a would where a User the User Service might need to search the Notifications table, and the Notifications Service  
needs to search the User Table.

The only thing about this implementation I am unsure about is it's testability. I am quite new to go, and even newer to Tests. My background in Nestjs has  
abstracted a lot of details from me. I am not confident in how testable these systems would be, and it something I would have to take time out to explore
