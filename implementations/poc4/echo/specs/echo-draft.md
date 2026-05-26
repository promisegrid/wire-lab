# poc4 echo app draft

`poc4-echo` carries a same-protocol request and response pair across relays.
Bob's echo client promises to receive an echo result before asking Ellen's echo
server to echo text. Ellen's app returns the result to Bob through the relay
mesh.

The app judges the echo result locally by comparing exact text and request hash.

