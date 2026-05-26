# poc4 storage app draft

`poc4-storage` models storage as reciprocal promises, not a remote filesystem
authority:

- The client promises to receive a store confirmation before asking another app
  to store a key-value pair.
- The storage app promises to remember the value locally and return a
  `store_confirm_v1` if the store promise is kept.
- The client later asks for the value by key while promising to receive a
  `read_result_v1`.

The storage app's memory is intentionally process-local for the proof. It is
evidence about app promise shape, not a durable storage design.

