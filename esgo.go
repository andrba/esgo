package esgo

import (
  "fmt"
  "strings"
  "encoding/json"
  "net/http"
  "io"
  "io/ioutil"
)

var server string

func Configure(host string, port int) {
  server = fmt.Sprintf("http://%s:%d", host, port)
}

func Request(method, url, data string) (body []byte, err error) {
  var response map[string]interface{}
  var httpStatusCode int
  var req *http.Request

  req, err = http.NewRequest(method, server + url, nil)
  if err != nil {
    return
  }

  req.Header.Set("Content-Type", "application/json")
  req.Header.Add("Accept", "application/json")
  setBody(req, strings.NewReader(data))

  httpStatusCode, body, err = do(req, &response)
  if err != nil {
    return
  }

  if httpStatusCode > 304 {

    jsonErr := json.Unmarshal(body, &response)
    if jsonErr == nil {
      if error, ok := response["error"]; ok {
        status, _ := response["status"]
        return body, fmt.Errorf("Error [%s] Status [%v]", error, status)
      }
    }
    return body, jsonErr
  }
  return
}

func setBody(req *http.Request, body io.Reader) {
  rc, ok := body.(io.ReadCloser)
  if !ok && body != nil {
    rc = ioutil.NopCloser(body)
  }
  req.Body = rc
  if body != nil {
    switch v := body.(type) {
    case *strings.Reader:
      req.ContentLength = int64(v.Len())
    }
  }
}

func do(req *http.Request, v interface{}) (int, []byte, error) {
  response, bodyBytes, err := processResponse(req, v)
  if err != nil {
    return -1, nil, err
  }
  return response.StatusCode, bodyBytes, err
}

func processResponse(req *http.Request, v interface{}) (*http.Response, []byte, error) {
  res, err := http.DefaultClient.Do(req)
  if err != nil {
    return nil, nil, err
  }

  defer res.Body.Close()
  bodyBytes, err := ioutil.ReadAll(res.Body)

  if err != nil {
    return nil, nil, err
  }

  if res.StatusCode > 304 && v != nil {
    jsonErr := json.Unmarshal(bodyBytes, v)
    if jsonErr != nil {
      return nil, nil, jsonErr
    }
  }
  return res, bodyBytes, err
}