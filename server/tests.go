package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	checks := []string{}
	testUrls := []string{
		"http://localhost:8080/api/auctions",
		"http://localhost:8080/api/auction",
		"http://localhost:8080/api/user",
		"http://localhost:8080/api/user/login",
	}

	runTestR(testUrls, &checks)

	fmt.Printf("\n--- Test Results ---\n")
	for i, c := range checks {
		fmt.Printf("%s: %s\n", c, testUrls[i])
	}
}

func runTestR(urls []string, checks *[]string) int {
	if len(urls) == 0 {
		return 0
	}

	fmt.Printf("\nRequesting %s...\n", urls[0])
	res, err := http.Get(urls[0])
	if err != nil {
		*checks = append(*checks, "FAILED")
	}
	printResponse(res)
	*checks = append(*checks, "PASSED")
	return runTestR(append(urls[:0], urls[1:]...), checks)
}

func printResponse(res *http.Response) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	fmt.Printf("Status Code: %d\n", res.StatusCode)
	fmt.Printf("Body: %s\n", string(body))
	for key, value := range res.Header {
		fmt.Printf("%s: %s\n", key, value[0])
	}
}
