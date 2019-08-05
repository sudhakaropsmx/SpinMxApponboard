package model

import (
    "encoding/json"
    "fmt"
)
type artifacttemplate struct {
   artifactAccount string
   reference string
 }
type pipeline struct {
    application  string
    keepWaitingPipelines bool
    limitConcurrent bool
    name string
    schema string
    template artifacttemplate
}

func main() {
    // Create an instance of the Box struct.
    artifact := artifacttemplate{
    } 
    pipelinejson := pipeline{
       
    }
    // Create JSON from the instance data.
    // ... Ignore errors.
    pipelinejson, _ := json.Marshal(pipelinejson)
    // Convert bytes to string.
    s := string(pipelinejson)
    fmt.Println(pipelinejson)
}