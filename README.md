# httpfetch
--
    import "github.com/rogpeppe/httpfetch"

Package httpfetch is an illustrative package used for a talk on testing in Go.
The talk is at
http://godoc.org/github.com/rogpeppe/talks/testing.talk/testing.slide .

## Usage

#### func  GetURLAsString

```go
func GetURLAsString(url string) (string, error)
```
GetURLAsString makes a GET request to the given URL and returns the result as a
string.
