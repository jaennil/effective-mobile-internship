package main

import (
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"runtime"
	"testing"
	"testing/quick"
	"time"

	"github.com/stretchr/testify/assert"
)

// task1
type addTest struct {
	arg1, arg2, expected int
}

func TestAdd(t *testing.T) {
	tests := []addTest{{2, 3, 5}, {4, 8, 12}, {123, 5734, 5857}}
	for _, test := range tests {
		actual := Add(test.arg1, test.arg2)
		if actual != test.expected {
			t.Errorf("Add(%v, %v) == %v, want %v", test.arg1, test.arg2, actual, test.expected)
		}
	}
}

// task2
type divideTest struct {
	arg1, arg2, expected int
}

func TestDivide(t *testing.T) {
	tests := []divideTest{{8, 2, 4}, {6, 3, 2}, {9, 2, 4}, {3, 0, -1}}
	for _, test := range tests {
		if test.arg2 == 0 {
			if checkDivByZero(func() { Divide(test.arg1, test.arg2) }) {
				continue
			} else {
				t.Errorf("expected div by zero for Divide(%v, %v), got no error\n", test.arg1, test.arg2)
			}
		}
		actual := Divide(test.arg1, test.arg2)
		if actual != test.expected {
			t.Errorf("Divide(%v, %v) == %v, want %v", test.arg1, test.arg2, actual, test.expected)
		}
	}
}

func checkDivByZero(f func()) (divByZero bool) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(runtime.Error); ok && e.Error() == "runtime error: integer divide by zero" {
				divByZero = true
			}
		}
	}()

	f()

	return false
}

// task3
func TestIntoFlattened(t *testing.T) {
	tests := []struct {
		arg      [][]int
		expected []int
	}{
		{
			[][]int{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			},
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
	}
	for _, test := range tests {
		actual := IntoFlattened(test.arg)
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("IntoFlattened(%v) == %v, want %v", test.arg, actual, test.expected)
		}
	}
}

// task4
func TestRandomNumber(t *testing.T) {
	f := func(min, max int) bool {
		n, err := RandomNumber(min, max)
		if err != nil {
			if max <= min-1 || -min+1 > math.MaxInt-max {
				return true
			}
			return false
		}
		return min <= n && n <= max
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

// task5
func TestRepeatHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/repeat?word=aoeu&count=3", nil)
	w := httptest.NewRecorder()

	repeatHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

    if res.StatusCode != http.StatusOK {
        t.Errorf("got response status code = %v; want %v", res.StatusCode, http.StatusOK)
    }

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	got := string(data)
	expected := "aoeuaoeuaoeu"
	if got != expected {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

// task6
func TestGetUserData(t *testing.T) {
    resource := "/api/users/2"
    wantUser := &User{
        Id:        3,
        Email:     "example@mail.fld",
        FirstName: "James",
        LastName:  "Johnson",
        Avatar:    "https://example.com/avatar.png",
    }

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != resource {
            t.Errorf("got url `%v`, want `%v`", r.URL.Path, resource)
		}
        w.WriteHeader(http.StatusOK)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(wantUser)
	}))
	defer mockServer.Close()

    url := mockServer.URL + resource
    gotUser, err := getUserData(url)
    if gotUser == nil {
        t.Errorf("getUserData(%v): got nil user, want non-nil user", url)
    }
    if err != nil {
        t.Errorf("getUserData(%v): got `%v` error, want no error", url, err)
    }
    if !reflect.DeepEqual(gotUser, wantUser) {
        t.Errorf("getUserData(%v): got `%+v` user, want `%+v`", url, gotUser, wantUser)
    }
}

// task7
func TestIsWeekend(t *testing.T) {
	moscowLocation, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		t.Errorf("failed to load location: %v", err)
	}
	friday := time.Date(2024, time.August, 2, 0, 0, 0, 0, moscowLocation)
	saturday := time.Date(2024, time.August, 3, 0, 0, 0, 0, moscowLocation)
	sunday := time.Date(2024, time.August, 4, 0, 0, 0, 0, moscowLocation)
	monday := time.Date(2024, time.August, 5, 0, 0, 0, 0, moscowLocation)
	if IsWeekend(friday) == true {
		t.Error("IsWeekend(saturday) == true, want false")
	}
	if IsWeekend(saturday) == false {
		t.Error("IsWeekend(saturday) == false, want true")
	}
	if IsWeekend(sunday) == false {
		t.Error("IsWeekend(saturday) == false, want true")
	}
	if IsWeekend(monday) == true {
		t.Error("IsWeekend(saturday) == true, want false")
	}
}

// task8
func TestIncrement(t *testing.T) {
	amount := 1000
	goroutinesAmount := 10
	var want int64 = 1000
	got := concurrentIncrement(amount, goroutinesAmount)
	if want != got {
		t.Errorf("increment(%v, %v) = %v, want %v", amount, goroutinesAmount, got, want)
	}
}

// task9
func TestFileContents(t *testing.T) {
	filename := "nonexistent_filename"
	got, err := fileContents(filename)
	if err == nil {
		t.Errorf("fileContents(%v): got nil error, want ErrFileNotFound error", filename)
	}
	if !errors.Is(err, ErrFileNotFound) {
		t.Errorf("fileContents(%v): got `%v` error, want ErrFileNotFound error", filename, err)
	}
	if got != "" {
		t.Errorf("fileContents(%v): got non-empty content, want empty sting", filename)
	}

	file, err := os.Create("file.txt")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	want := "hello, world!"
	_, err = file.WriteString(want)
	if err != nil {
		t.Error(err)
	}
	got, err = fileContents(file.Name())
	if err != nil {
		t.Errorf("fileContents(%v): got `%v` error, want no error", file.Name(), err)
	}
	if got != want {
		t.Errorf("fileContents(%v) == %v, want %v", file.Name(), got, want)
	}
}

// task10
func TestHelloHandler(t *testing.T) {
	server := NewServer()

	request, err := http.NewRequest(http.MethodGet, "/hello", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()

	server.ServeHTTP(rr, request)

	want := `{"message":"Hello, World!"}`
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, want, rr.Body.String())
}
