package http_client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/JetBlink/orderbook/base"
)

type Response struct {
	request  *http.Request
	response *http.Response
	body     []byte
}

type ResponseBody struct {
	Code    string          `json:"code"`
	RawData json.RawMessage `json:"data"` // delay parsing
	Message string          `json:"msg"`
}

func NewResponse(request *http.Request, response *http.Response) *Response {
	resp := &Response{request, response, nil}
	resp.body, _ = resp.read()

	return resp
}

func (resp *Response) success() bool {
	if resp.response.StatusCode != http.StatusOK {
		return false
	}

	return true
}

func (resp *Response) read() ([]byte, error) {
	defer resp.response.Body.Close()
	return ioutil.ReadAll(resp.response.Body)
}

func (resp *Response) StatusCode() int {
	return resp.response.StatusCode
}

func (resp *Response) Body() []byte {
	return resp.body
}

func (resp *Response) ReadJson(v interface{}) error {
	if !resp.success() {
		return resp.error(fmt.Sprintf("http code is not %d", http.StatusOK))
	}

	var responseBody ResponseBody
	if err := json.Unmarshal(resp.body, &responseBody); err != nil {
		return resp.error("ResponseBody json unmarshal failure")
	}

	const ApiSuccess = "200000"

	if responseBody.Code != ApiSuccess {
		return resp.error("api code is not " + ApiSuccess)
	}

	if os.Getenv("API_DEBUG_MODE") == "true" {
		log.Println(resp.request.URL.String() + " #####:### " + base.ToJsonString(responseBody.RawData))
	}
	if err := json.Unmarshal(responseBody.RawData, v); err != nil {
		return resp.error(fmt.Sprintf("responseBody.RawData json unmarshal failure: %v", err))
	}

	return nil
}

func (resp *Response) error(error string) error {
	return fmt.Errorf(
		"http request failure, error: %s\nstatus code: %d, %s %s, body:\n%s",
		error,
		resp.response.StatusCode,
		resp.request.Method,
		resp.request.URL.String(),
		string(resp.body),
	)
}
