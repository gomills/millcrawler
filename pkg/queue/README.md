# /queue

This package is dedicated to a concurrency and memory safe queue. It spawns goroutines to put urls in queue, which is memory dangerous, but their lifecycle is managed by the general context.

It also provides a vault for found secrets with deduplication and source information.