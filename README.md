# Counter

[![Counter](https://github.com/jpcercal/counter/actions/workflows/counter.yml/badge.svg?branch=main)](https://github.com/jpcercal/counter/actions/workflows/counter.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jpcercal/counter)](https://goreportcard.com/report/github.com/jpcercal/counter)
[![codecov](https://codecov.io/gh/jpcercal/counter/graph/badge.svg?token=DU7SHVHB6A)](https://codecov.io/gh/jpcercal/counter)

## Requirements

This project is written in golang using the latest version available at the time this was coded, but you only need docker to run it on your machine. =)

## Github Actions (pipeline)

There are some github jobs configured as part of the github action `Counter` configured on this project, they are:

- `docker-security-checker`;
- `go-security-checker`;
- `golangci-lint`;
- `test`. 

`docker-security-checker` validates that no vulnerability critical or high are reaching production.

`go-security-checker` validates that the golang project does not contain known vulnerabilities.

The `golangci-lint` one is there to ensure that all the rules defined by the linter are being followed.

While the `test` is there to ensure that the project works as expected. 100% of coverage was not the goal in here, but the critical parts of the application are tested. By the way, there I had been configured badges to highlight the quality of the deliverables here defined with the coverage report.

## Running the Project

You will have to build the docker image first, you can do so by running the following command:

```
docker build --tag jpcercal/counter .
```

Once it's done, you can spin up the http web server with the following command:

```
docker run -it --rm -v $(pwd)/counter.json:/app/counter.json -p 3000:3000 jpcercal/counter
```

The command above will start the HTTP server on the port 3000, you must make sure that this port is available on your system before proceding, if needed adjust the TCP port according to your needs.

Once the server shuts down, the file informed as a volume parameter to the container on the path `/app/counter.json` will store a collection of timestamps which as a consequence will be reflected to your host machine, this information will them be used to warm up the state of the thread safe variable that stores the state of the counter on the server once it starts, again.

The log below is shown when the server starts running:

```
2024/01/27 23:17:43 configuring server...
2024/01/27 23:17:43 starting server...
2024/01/27 23:17:43 loading state from disk...
2024/01/27 23:17:43 the counter time window is 60 seconds...
2024/01/27 23:17:43 Listening on :3000
```

There are two endpoints registered on the application, they are:

### `GET /healthcheck`

This endpoint exists to ensure that the server is healthy, and it can be used lately by a monitoring application, and by kubernetes to orchestrate the service pods;

Below you can see this endpoint in action:

```
http GET http://localhost:3000/healthcheck

HTTP/1.1 200 OK
Content-Length: 41
Content-Type: application/json
Date: Sat, 27 Jan 2024 23:17:57 GMT

{
    "status": "OK",
    "timestamp": 1706397477678
}
```

### `GET /counter`

This is the main goal of this application, to count each request made to the webserver (affecting this endpoint only) in a time window of 60s;

The first request looks like this:

```
http GET http://localhost:3000/counter

HTTP/1.1 200 OK
Content-Length: 13
Content-Type: application/json
Date: Sat, 27 Jan 2024 23:19:50 GMT

{
    "counter": 1
}
```

The subsequent request looks like this, note that it increaments the value of the counter from `1` to `2`:

```
http GET http://localhost:3000/counter

HTTP/1.1 200 OK
Content-Length: 13
Content-Type: application/json
Date: Sat, 27 Jan 2024 23:19:54 GMT

{
    "counter": 2
}
```

However, after `60 secounds` have passed since the request was fired up, requesting the value of the counter once again from the server would have a similar response like the one shown below:

```
http GET http://localhost:3000/counter

HTTP/1.1 200 OK
Content-Length: 13
Content-Type: application/json
Date: Sat, 27 Jan 2024 23:24:53 GMT

{
    "counter": 1
}
```

As you can see the counter resets automatically. Cool, right!?

If for some reason the server is restarted or terminated, a file will be persisted to keep track of everything that still might count on this time window period, note this log message `saving state to disk on a file named counter.json`.

```
2024/01/27 23:25:09 Shutting down server... Reason: interrupt
2024/01/27 23:25:09 Server gracefully stopped
2024/01/27 23:25:09 saving state to disk on a file named counter.json...
```

The file `counter.json` would look like this once the server got stopped:

```
["2024-01-27T23:24:53.774042545Z"]
```

And then, restarting the server it would come back restoring the data that was on the session and performing some clean up once the endpoint `GET /counter` is reached again:

```
2024/01/27 23:25:36 configuring server...
2024/01/27 23:25:36 starting server...
2024/01/27 23:25:36 loading state from disk...
2024/01/27 23:25:36 the counter time window is 60 seconds...
2024/01/27 23:25:36 Listening on :3000
```

So, by firing up a request immediately before the time window moves out of the 60s period you will have `2` as a value of the counter as it can be seen below:

```
http GET http://localhost:3000/counter

HTTP/1.1 200 OK
Content-Length: 13
Content-Type: application/json
Date: Sat, 27 Jan 2024 23:25:39 GMT

{
    "counter": 2
}
```

Amazing! Ah, the operation is thread-safe, rry it out or just take a look on the tests. =)

### The default handler

Any other endpoint that not the ones specified above will return a `404` HTTP Status Code, as they of course, do not exist.

```
http GET http://localhost:3000/this-does-not-exist

HTTP/1.1 404 Not Found
Content-Length: 0
Date: Sat, 27 Jan 2024 23:32:50 GMT
```

### An experiment

Even though some modern and widely adopted libraries are in being heavily used for handling the web server, routing, logging, etc. I decided that I wanted to stay with only the golang standard packages to experiment with it. 

The only external library used in here is `cobra`, as the cmd package uses it just to have a command line application to start the server.

I would consider using `gin`, `chi`, `viper`, `testcontainers` with `docker`, `gock`, but you will see it around the code tagged with some TODOs too. 

Some other file formats were also considered for persisting the state of the application, but I stayed with json just because of portability and because it's easy to read and not critical to the application, as it is not expected that the server will be restarted all the time.

Defining data contracts by using an openapi definition file would also be a must to be present in any modern REST API, although I left it out of the scope of integrating it with the go standard libraries.
