package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

func main() {
    hello()
    variables()
    r := even(12)
    fmt.Println(r)
    r = even(11)
    fmt.Println(r)
    loop()
    added := add(45, 90)
    fmt.Println(added)
    s := sum([]int{1, 23, 42})
    fmt.Println(s)
    structs()
    pointers()
    div, err := Errors(12, 3)
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println(div)
    }
    div, err = Errors(12, 0)
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println(div)
    }
    goroutinesAndChannels()
}

func hello() {
    fmt.Println("Hello, World!")
}

func variables() {
    var v1 int = 100
    var v2 string = "golang"
    var v3 bool = true
    fmt.Println(v1)
    fmt.Println(v2)
    fmt.Println(v3)
}

func even(n int) bool {
    if n % 2 == 0 {
        return true
    } else {
        return false
    }
}

func loop() {
    for i := range 5 {
        fmt.Println(i+1)
    }
}

func add(a, b int) int {
    return a + b
}

func sum(s []int) int {
    var sum int
    for n := range s {
        sum += n
    }
    return sum
}

func structs() {
    p1 := Person {
        name: "Nikita",
        age: 20,
        city: "Moscow",
    }
    fmt.Println(p1)
}

type Person struct {
    name string
    age int
    city string
}

func pointers() {
    p := &Person {
        name: "Ivan",
        age: 25,
        city: "Paris",
    }
    fmt.Println(p)
    p.age = p.age + 1
    fmt.Println(p)
}

func Errors(a, b int) (int, error) {
    if b == 0 {
        return 0, DivisionByZero
    }
    return a / b, nil
}

var DivisionByZero = errors.New("Division by zero")

type MockFetcher map[string]*mockResult

type mockResult struct {
	body string
	urls []string
}

func (f MockFetcher) Fetch(url string) (string, []string, error) {
	fetchSignalInstance() <- true
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

var fetcher = MockFetcher{
	"http://golang.org/": &mockResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &mockResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &mockResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &mockResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}

var fetchSignal chan bool

func fetchSignalInstance() chan bool {
	if fetchSignal == nil {
		fetchSignal = make(chan bool, 1000)
	}
	return fetchSignal
}

func Crawl(url string, depth int, wg *sync.WaitGroup, throttle <-chan time.Time) {
	defer wg.Done()

	if depth <= 0 {
		return
	}

	<-throttle

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("found: %s %q\n", url, body)

	wg.Add(len(urls))
	for _, u := range urls {
		go Crawl(u, depth-1, wg, throttle)
	}
	return
}

func goroutinesAndChannels() {
	var wg sync.WaitGroup
	throttle := time.Tick(time.Second)
	wg.Add(1)
	Crawl("http://golang.org/", 4, &wg, throttle)
	wg.Wait()
}
