package esgo

import (
  "fmt"
  "strings"
  "encoding/json"
  "net/http"
  "io/ioutil"
)

var server string

func Configure(host string, port int) {
  server = fmt.Sprintf("http://%s:%d", host, port)
}

func Request(method, url, data string) (body []byte, err error) {
  var request *http.Request
  var response *http.Response

  if request, err = buildRequest(method, server + url, data); err != nil {
    return
  }

  response, err = http.DefaultClient.Do(request)
  defer response.Body.Close()

  if err != nil {
    return
  }

  if body, err = ioutil.ReadAll(response.Body); err != nil {
    return
  }

  if response.StatusCode > 304 {

    var v map[string]interface{}

    jsonBody := json.Unmarshal(body, &v)

    if jsonBody == nil {
      if error, ok := v["error"]; ok {
        status, _ := v["status"]
        return body, fmt.Errorf("Error [%s] Status [%v]", error, status)
      }
    }
    return body, jsonBody
  }
  return
}

func buildRequest(method, url, data string) (req *http.Request, err error) {
  req, err = http.NewRequest(method, url, strings.NewReader(data))
  if err != nil {
    return
  }

  req.Header.Set("Content-Type", "application/json")
  req.Header.Add("Accept", "application/json")

  return
}