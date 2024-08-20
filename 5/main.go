package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func main() {}

// task1
func Add(x, y int) int {
	return x + y
}

// task2
// TODO: return error if y is 0
func Divide(x, y int) int {
	return x / y
}

// task3
func IntoFlattened(slice [][]int) (result []int) {
	for _, v := range slice {
		result = append(result, v...)
	}
	return result
}

// task4
// [min, max]
func RandomNumber(min, max int) (int, error) {
	if max <= min-1 || -min+1 > math.MaxInt-max {
		return -1, errors.New("invalid range")
	}
	return rand.IntN(max-min+1) + min, nil
}

// task5
// Req: http://localhost:8080/repeat?word=abc&count=3
// Res: abcabcabc
func repeatHandler(w http.ResponseWriter, r *http.Request) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid request")
		return
	}

	word := query.Get("word")
	if len(word) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "missing word to repeat")
		return
	}

	countStr := query.Get("count")
	if len(countStr) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "missing repeat count")
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "repeat count must be number")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, strings.Repeat(word, count))
}

// task6
type User struct {
	Id        int
	Email     string
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string
}

func getUserData(url string) (*User, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get user data")
	}

	var user User
	if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// task7
func IsWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// task8
func concurrentIncrement(amount int, goroutinesAmount int) int64 {
	counter := atomic.Int64{}
	wg := sync.WaitGroup{}
	wg.Add(goroutinesAmount)
	for _ = range goroutinesAmount {
		go func() {
			for _ = range amount / goroutinesAmount {
				counter.Add(1)
			}
			wg.Done()
		}()

	}
	wg.Wait()
	return counter.Load()
}

// task9
var ErrFileNotFound = errors.New("file not found")

func fileContents(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if errors.Is(err, os.ErrNotExist) {
		return "", ErrFileNotFound
	}

	if err != nil {
		return "", err
	}

	return string(content), nil
}

// task10
func NewServer() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	return mux
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{Message: "Hello, World!"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

type Response struct {
	Message string `json:"message"`
}
