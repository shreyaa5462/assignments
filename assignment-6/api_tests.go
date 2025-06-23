package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestPost(t *testing.T) {

	testput := []struct {
		url      string
		expected string
	}{
		{
			url:      "http://localhost:8080/task?task=dance111",
			expected: "dance111",
		},
		{
			url:      "http://localhost:8080/task",
			expected: "Missing id parameter",
		}, {
			url:      "http://localhost:8080/task?task",
			expected: "Missing id parameter",
		},
	}
	for _, test := range testput {
		resp, err := http.Post(test.url, "application/json", nil)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

		if string(bodyBytes) != test.expected {
			t.Errorf("FAIL: Expected %s, got %s", test.expected, string(bodyBytes))
		} else {
			t.Logf("PASS: Test %s", test.expected)
		}

	}
}

func TestGet(t *testing.T) {

	url := "http://localhost:8080/task"

	// Make the GET request
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		t.Logf("PASS: Test Get")
	} else {
		t.Errorf("Expected http.StatusOK , got %d", resp.StatusCode)
	}
}

func TestPut(t *testing.T) {

	testput := []struct {
		url      string
		expected string
	}{
		{
			url:      "http://localhost:8080/task?id=1",
			expected: "1",
		},
		{
			url:      "http://localhost:8080/task?id=2",
			expected: "2",
		}, {
			url:      "http://localhost:8080/task?id",
			expected: "give the ID",
		},
		{
			url:      "http://localhost:8080/task?id=111111",
			expected: "Please Enter the valid ID",
		},
		{
			url:      "http://localhost:8080/task?id=-100",
			expected: "Please Enter the valid ID",
		},
	}
	//tracker := NewTaskTracker()
	for _, test := range testput {

		req, err := http.NewRequest("POST", test.url, nil)

		if err != nil {
			log.Fatal(err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		// Read the response
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode == http.StatusOK {
			t.Logf("PASS: TestPut")
		} else if string(bodyBytes) == test.expected {
			t.Logf("PASS: TestPut")
		} else {
			t.Errorf("FAIL: %s", bodyBytes)
		}
	}
}

func TestDelete(t *testing.T) {
	url := "http://localhost:8080/task"
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	} else {
		t.Logf("PASS: TestDelete %s", string(bodyBytes))
	}

}
