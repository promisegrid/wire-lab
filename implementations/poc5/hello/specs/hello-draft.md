# poc5 hello app draft

`poc5-hello` is intentionally tiny. In the five-container demo Alice's hello
app asks Dave's signed app for signed evidence over a hello text. Ellen's hello
app exists as a second colocated app and makes no network promise in this
bounded run.

The point is not chat. The point is that a simple app can promise to receive a
result, carry an app-level request through relays, and judge the returned
signed evidence locally.

