# Token Bucket

## Why Token Bucket?

We chose the Token Bucket algorithm because it allows burst traffic while enforcing a long-term average request rate.

Compared to Fixed Window, it avoids the boundary burst problem.

Compared to Sliding Window Log, it requires significantly less memory.

## Advantages

- Burst support
- Constant memory
- Widely used
- Easy to distribute with Redis

## Disadvantages

- Slightly more complex than Fixed Window
- Requires storing timestamps