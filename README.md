[![GoDoc](https://godoc.org/github.com/Ronmi/pastebin?status.svg)](https://godoc.org/github.com/Ronmi/pastebin)

pastebin API binding for golang.

```go
api := pastebin.API{Key: "my_dev_key"}

// post a public paste
uri, err := api.Post(&paste.Post{
    Title: "markdown test",
    Content: `my cool content`,
    Format: "markdown",
    ExpireAt: pastebin.In1Y,
])
if err != nil {
    // handle error
}
log.Printf("new paste url: %s", uri)

// grab all pastes created by user "darius"
ukey, err := api.UserKey("darius", "my_password")
if err != nil {
    // handle error here
}
posts, err := api.List(ukey, 0)
if err != nil {
    // handle error here
}

for _, p := range posts {
    log.Printf("Post#%s: %s", p.Key, p.Title)
}
```

# License

WTFPL
