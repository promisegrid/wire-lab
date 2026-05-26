# poc5 fibonacci app draft

`poc5-fibonacci` models delayed computation as reciprocal promises:

- The client promises to receive a future `fibonacci_result_v1`.
- The server promises to calculate `fib(n)` only after receiving a
  `fibonacci_request_v1` that names where the result should be returned.
- The client judges keep/break locally from `n`, `result`, and the exact request
  hash.

The server does not become a remote procedure endpoint. It is an app-level
promiser whose result happens to traverse relay promises.

