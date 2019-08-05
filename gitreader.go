package main

import (
 "fmt"
 

)


repo, err := gitreader.OpenRepo("/OpsMx/Onboarding-Spinnaker-API/blob/master/spinnaker-user-onboarding/pipeline-template")
if err != nil {
  panic(err)
}

blob, err := repo.CatFile("HEAD", "OpsMx/Onboarding-Spinnaker-API/blob/master/spinnaker-user-onboarding/pipeline-template")
if err != nil {
  panic(err)
}

// WARNING: use Blob as an io.Reader instead if you can!
bytes, err := blob.Bytes()
if err != nil {
  panic(err)
}

fmt.Printf("%s", bytes)

repo.Close()
