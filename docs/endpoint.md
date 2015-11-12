# endpoint

The endpoint block is a server method. Once the endpoint's `name` is
set, it will receive HTTP requests from the connected server, The
request is emitted from three routes: the request object itself from
`request`, a writer object that can be written to, flushed, and/or
closed from `writer`, and the body of the request from `body`. 
