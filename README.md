# redis-go

A simple Redis server implementation in Go. This is a mock example I built to learn about network programming and the Redis protocol.

## What it does

It's a basic Redis clone that handles a few commands:
- `PING` - returns PONG
- `ECHO` - echoes back what you send
- `SET` - stores a key-value pair (with optional expiration)
- `GET` - retrieves a value by key

## Running it

```bash
./your_program.sh
```

The server runs on port 6379 (standard Redis port).

## Testing

You can test it with `redis-cli` or netcat:

```bash
# Using redis-cli
redis-cli ping
redis-cli set mykey "hello"
redis-cli get mykey

# Using netcat
echo -e '*1\r\n$4\r\nPING\r\n' | nc localhost 6379
```

## Notes

This is just a learning project - it implements the bare minimum of the Redis protocol to understand how it works. Don't use this in production! 😄