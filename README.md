# Jester 🎭  
A lightweight JSON handling library in Go. Built from [github.com/bitly/go-simplejson](https://github.com/bitly/go-simplejson).

## Why?
- Uses [github.com/goccy/go-json](https://github.com/goccy/go-json) instead of `encoding/json`.
- `Get()` supports string as well as int keys to index maps and slices in one call.
- Added `Len()` func to get the length of the underlying data.
- I guess that's all.

## Installation  
```sh
go get github.com/lb-selfbot/go-jester
```
