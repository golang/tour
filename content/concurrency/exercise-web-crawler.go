// +build OMIT

package main

import (
	"fmt"
)

type Fetcher interface {
	// Fetch vrací základ URL a
	// slice odkazů URL nalezených na té stránce.
	Fetch(url string) (body string, urls []string, err error)
}

// Crawl používá fetcher aby rekurzivn proparsoval
// stránky začínající s url až do maximální hloubky.
func Crawl(url string, depth int, fetcher Fetcher) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	if depth <= 0 {
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("našel: %s %q\n", url, body)
	for _, u := range urls {
		Crawl(u, depth-1, fetcher)
	}
	return
}

func main() {
	Crawl("https://golang.org/", 4, fetcher)
}

// fakeFetcher je Fetcher, který vrací předpřipravené výsledky.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("nenalezeno: %s", url)
}

// fetcher je naplněný fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"Programovací jazyk Go",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Balíčky",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Balíček fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Balíček os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
