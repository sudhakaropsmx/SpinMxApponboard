package complianceapi

import (
 "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)
type InputJson struct{
   User string
   Groups []string
}
func CheckApplicationAccess(input interface{}) (interface{}, *http.Response,error)  {
	var successPayload  interface{}
    
    //fmt.Println("Starting the application...")
    
    response, err := http.Get("http://localhost:8000/")
    if err != nil {
        fmt.Printf("The HTTP request failed with error %s\n", err)
        return successPayload, nil, err
    } else {
        data, _ := ioutil.ReadAll(response.Body)
        fmt.Println(string(data))
    }

    jsonValue, _ := json.Marshal(input)
    response, err = http.Post("http://localhost:8000/api/getUserAppAuthorized", "application/json", bytes.NewBuffer(jsonValue))
    
    if err != nil {
        fmt.Println("The HTTP request failed with error %s\n", err)
        return successPayload, response, err
    } 
    	
	if err = json.NewDecoder(response.Body).Decode(&successPayload); err != nil {
		fmt.Println("The HTTP response failed with error %s\n", err)
		return successPayload, response, err
	}	
    
    return successPayload, response, nil
}