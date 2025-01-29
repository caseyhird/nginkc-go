# nginkc-go

A simple HTTP server in Go. Just made for fun and my own learning. Inspired by [build your own x](https://github.com/codecrafters-io/build-your-own-x).


Usage:
Can be used following the demonstration in the app/ module.
- require & import the nginkc module
- create an app that implementes the nginkc App interface, i.e. has a "Call" method that accepts a "Request" and returns a "Response"
- serve this app with nginkc.Serve
That's it!


Future work:
- Add support for later HTTP versions
- Maybe setup a web framework or something else more interesting for the app
