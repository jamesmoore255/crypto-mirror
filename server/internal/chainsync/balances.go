package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	url := "https://eth-mainnet.g.alchemy.com/v2/docs-demo"

	payload := strings.NewReader(`{\"id\":1,\"jsonrpc\":\"2.0\",\
	"method\":\"alchemy_getTokenBalances\",\"params\":[\"0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5\",\"erc20\"]}`)

	req, err := http.NewRequestWithContext(context.Background(), "POST", url, payload)

	if err != nil {
		panic(err) // don't do this
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err) // don't do this
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}
