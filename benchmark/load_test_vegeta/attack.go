package main

import (
  "fmt"
  "io/ioutil"
  "time"
  "github.com/tsenart/vegeta/lib"
  "math"
  "os"
  "bytes"
)

func main() {
  rate := uint64(250) // per second
  duration := 4 * time.Second
  body, err := ioutil.ReadFile("./reqBody.json")
  if err != nil {
    fmt.Println("err:", err)
    return
  }
  //fmt.Println("body:", string(body))
  targeter := vegeta.NewStaticTargeter(vegeta.Target{
    Method: "POST",
    URL:    "http://localhost:8080/url",
    Body:   body,
  })
  attacker := vegeta.NewAttacker()

  var metrics vegeta.Metrics
  var results vegeta.Results
  for res := range attacker.Attack(targeter, rate, duration, "") {
    metrics.Add(res)
    results.Add(res)
  }
  metrics.Close()

  fmt.Println("Max latency:", metrics.Latencies.Max)
  fmt.Println("Requests:", metrics.Requests)
  fmt.Println("Rate:", math.Round(metrics.Rate))
  fmt.Println("BytesIn:", metrics.BytesIn)
  fmt.Println("BytesOut:", metrics.BytesOut)
  fmt.Println("Errors:", metrics.Errors)
  fmt.Println("StatusCodes: ", metrics.StatusCodes)
  fmt.Println("Success: ", metrics.Success)

  //configure to generate a html report
  var buf bytes.Buffer
  rep := vegeta.NewPlotReporter("Plot Report", &results )
  rep.Report(&buf)
  file, err := os.Create(`./rep.html`)
  defer file.Close()
  file.Write(buf.Bytes())
}
