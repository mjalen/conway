# Conway's Game of Life 

https://github.com/mjalen/conway/assets/54542654/402ff04c-2e81-4501-962d-114471583426

A simple game of life written in Go.

# Why?

Because I wanted to. Or better yet, because I can :). In all seriousness though, I just wanted to learn Go and I have always wanted to implement my own Conway's Game of Life. I learned a lot about Go, HTTP, and concurrency in this project. I also learned how fun Go is to program.

# How does it work?

I intentionally tried not to look at other implementations of Conway's Game of Life. The implementation I wrote tracks the system state using a slice of alive cells. At each iteration, only the neighborhoods of the currently alive cells are checked and updated. Checks are done concurrently. I have implemented wrapping at the system edges; the left wraps to the right and the top wraps to the bottom. Wrapping was implemented as the better alternative to hard borders. A *true* Conway's Game of Life would track an infinite plane. This can easily be added to the current implementation, although it is very slow and typically blows up in population. Hence, I opted for wrapping at the edges. Perhaps if I feel like it, I will make a branch with an infinite plane.

# Stack

The following technologies were utilized in this project:

- [HTMX](https://htmx.org/). For handling HTTP requests on the frontend and for establishing a [SSE](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events) connection.
- [Templ](https://github.com/a-h/templ). For templating HTML components in Go.
- [Go](https://go.dev/). The language of choice. This project was an excuse to learn Go in the first place.

# Install and Usage 

Simply clone this repository to install. Also make sure Templ is installed:

``` bash
$ go install github.com/a-h/templ/cmd/templ@latest
```

To build and run, execute the following

``` bash
$ templ generate # Only run if any .templ files were changed. 
$ go run . # Run with --help for usage.
```
