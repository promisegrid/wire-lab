# TE-1: Promise-stack ordering

*Thought experiment, part of the [PromiseGrid Wire Lab](../../harness-spec.md). This file is content-addressable; its hash is its pCID.*

Run two variants — signature outermost vs. signature innermost — across stream and datagram transports. Find the case where one ordering forces buffering, and the case where the other ordering loses information when an intermediate router strips a frame. Outcome: a design rule for which burdens go where.
