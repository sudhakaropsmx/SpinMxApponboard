// Copyright (c) 2019, OpsMx, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package pipeline

import (
	"encoding/json"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	//"strings"
	"text/template"
	//"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spinnaker/spin/cmd/gateclient"
	"github.com/spinnaker/spin/util"
	pconfig "github.com/sudhakaropsmx/spinmx/config/pipeline"
	//appconfig "github.com/sudhakaropsmx/spinmx/config/application"
	capi "github.com/sudhakaropsmx/spinmx/complianceapi"
	
	"gopkg.in/yaml.v2"
	
)

type CreatePipelineOptions struct {
	*pipelineOptions 
	output       string
	pipelineConfigFile string
}

var (
	createPipelineShort = "Create the pipeline with the provided config values and then save "
	createPipelineLong  = "Create the pipeline with the provided config values and then save "
)

func NewCreatePipelineCmd(pipelineOptions pipelineOptions) *cobra.Command {
	options := CreatePipelineOptions{
		pipelineOptions: &pipelineOptions,
	}
	cmd := &cobra.Command{
		Use:     "create",
		Short:   createPipelineShort,
		Long:    createPipelineLong,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCreatePipeline(cmd, options)
		},
	}

	cmd.PersistentFlags().StringVarP(&options.pipelineConfigFile, "file", "f", "", "path to the pipeline config values file")

	return cmd
}

func getCreatePipeline(cmd *cobra.Command, options CreatePipelineOptions) error { 
	
	gateClient, err := gateclient.NewGateClient(cmd.InheritedFlags())
	
	if err != nil {
	return err
	}	
	pipelineConfig := &pconfig.PipelineConfig{}
	
	yamlFile, err := ioutil.ReadFile(options.pipelineConfigFile)

	if yamlFile != nil {
		
		err = yaml.UnmarshalStrict([]byte(os.ExpandEnv(string(yamlFile))), &pipelineConfig)
		
		if err != nil {
			util.UI.Error(fmt.Sprintf("Could not deserialize config file with contents: %s, failing.", yamlFile))
			return err
		}
	}
	
	valid := true
	application := pipelineConfig.Application
	pipelineName := pipelineConfig.PipelineName
	pipelinetemplatename := pipelineConfig.TemplateReference

	if  pipelineName == "" {
		util.UI.Error("Required pipeline key 'name' missing...\n")
		valid = false
	}
    
	if application == "" {
		util.UI.Error("Required pipeline key 'application' missing...\n")
		valid = false
	}
	if pipelinetemplatename == "" {
		util.UI.Error("Required pipeline key pipeline template name missing...\n")
		valid = false
	}
	if !valid {
		return fmt.Errorf("Submitted pipeline is invalid data: \n")
	}
	
	flag, err := checkApplicationAuthorized(gateClient, application)
	if flag {
	  return err
	}
    //return nil
	foundPipeline, queryResp, _ := gateClient.ApplicationControllerApi.GetPipelineConfigUsingGET(gateClient.Context, application, pipelineName)

	if queryResp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error querying pipeline, status code: %d\n", queryResp.StatusCode)
	}
	if len(foundPipeline) > 0 {
		return fmt.Errorf("Encountered an error saving pipeline, pipeline already exist in the application: %d\n",application)
	}
	
	queryParams := map[string]interface{}{}
	queryParams["tag"] = "latest"
	_, resp, queryErr := gateClient.V2PipelineTemplatesControllerApi.GetUsingGET2(gateClient.Context, pipelinetemplatename, queryParams)
   
    if resp.StatusCode == http.StatusNotFound {
	    return fmt.Errorf("Encountered an error saving pipeline, pipeline template does't exist : %d\n",pipelinetemplatename)
    }
    if queryErr != nil {
    	fmt.Errorf("Encountered an error saving pipeline : %s\n",queryErr)
        return queryErr
     }
	
	pipelineJsonStr := `{"schema":"v2","name":"{{.PipelineName}}","application":"{{.Application}}","template":{"artifactAccount": "front50ArtifactCredentials","reference":"spinnaker://{{.TemplateReference}}","type": "front50/pipelineTemplate"},"type":"templatedPipeline"}`
	
	// Create a new template and parse the pipelineJsonStr into it.
	t := template.Must(template.New("pipelineJsonStr").Parse(pipelineJsonStr))

	// Execute the template for each recipient.
	var pipelineJsonbuf bytes.Buffer
	err = t.Execute(&pipelineJsonbuf, pipelineConfig)
	if err != nil {
		fmt.Println("executing template:", err)
		return err
	}
	
	pipelineJsonMap := make(map[string]interface{})
	
	convertErr := json.Unmarshal(pipelineJsonbuf.Bytes(), &pipelineJsonMap)

	if convertErr != nil {
	  return fmt.Errorf("Encountered an error saving pipeline, Json to Map Convert Error: %d\n",convertErr)
	}
	
	if pipelineConfig.Variables != nil  {
	  pipelineJsonMap["variables"] = pipelineConfig.Variables
	}
	 
	if len(pipelineConfig.Triggers) >0{
		pipelineJsonMap["triggers"] = pipelineConfig.Triggers
	} 
	if len(pipelineConfig.ExpectedArtifacts) >0{
		pipelineJsonMap["expectedArtifacts"] = pipelineConfig.ExpectedArtifacts
	} 
	if len(pipelineConfig.Notifications) >0{
		pipelineJsonMap["notifications"] = pipelineConfig.Notifications
	} 
	if len(pipelineConfig.Parameters) >0{
		pipelineJsonMap["parameters"] = pipelineConfig.Parameters
	}
	//for key, value := range pipelineJsonMap {
    //  fmt.Println("index : ", key, " value : ", value)
	//}
	
	//saveResp, saveErr := gateClient.PipelineControllerApi.SavePipelineUsingPOST(gateClient.Context, pipelineJsonMap)

	//return nil
	jsonString, err := json.Marshal(pipelineJsonMap)
    fmt.Println("Pipeline :\n",string(jsonString))
    
	saveResp, saveErr := gateClient.PipelineControllerApi.SavePipelineUsingPOST(gateClient.Context, pipelineJsonMap)

	if saveErr != nil {
		fmt.Printf("   s err: %v", saveErr)
		return saveErr
	}
	if saveResp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error saving pipeline, status code: %d\n", saveResp.StatusCode)
	}
	
	util.UI.Info(util.Colorize().Color(fmt.Sprintf("[reset][bold][green]Pipeline save succeeded")))
  	return nil
}
func checkApplicationAuthorized(gateClient *gateclient.GatewayClient,application string) (bool,error) { 
  
    _, resp, err := gateClient.ApplicationControllerApi.GetApplicationUsingGET(gateClient.Context, application, map[string]interface{}{"expand": false})
    //fmt.Errorf(resp)
	if resp != nil {
		if resp.StatusCode == http.StatusNotFound {
			//fmt.Printf("Application '%s' not found\n", app)
			return true, fmt.Errorf("Application '%s' not found\n", application)
		} else if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		   return true, fmt.Errorf("Does not have acces to create pipeline in the application '%s' \n",application)
		}
	}

	if err != nil {
		return true, fmt.Errorf("Encountered an error getting application %d\n",err)
	} 
	/*
	prettyStr, _ := json.MarshalIndent(app["attributes"], "", " ")
    appdata := make(map[string]interface{})
    err = json.Unmarshal(prettyStr, &appdata)
    applicationConfig := &appconfig.AplicationConfig{}
    mapstructure.Decode(appdata, &applicationConfig)
    var write []string  
    write = applicationConfig.Permissions.Write
    var execute []string  
    execute = applicationConfig.Permissions.Execute
    groups := removeDuplicates(append(write, execute...))
    fmt.Printf("Application Groups: %s \n",groups)    
    input := make(map[string]interface{})
    input["User"] = gateClient.Config.Auth.Basic.Username
    input["Groups"] = groups
    */
    input := make(map[string]interface{})
    input["UserName"] = gateClient.Config.Auth.Basic.Username
    input["Application"] = application
    data, resp, err := capi.CheckApplicationAccess(input)
    //util.UI.JsonOutput(data, util.UI.OutputFormat)
    if err != nil {
		return true, fmt.Errorf("Error encountered in compliance API %s \n",err)
	} 

   if resp.StatusCode != http.StatusOK {
		return true, fmt.Errorf("Encountered an error in Compliance API pipeline")
	}
    type myResult struct{
      UserName string
      Applicaiton string
      Groups []string
    } 
    var mydata myResult
    prettyStra, _ := json.MarshalIndent(data, "", " ")
    err = json.Unmarshal(prettyStra, &mydata)
   // compData := make(map[string]interface{})
   // err = json.Unmarshal(prettyStra, &respData)
   // jsonString, err := json.Marshal((compData["Groups"]))
    //fmt.Println(string(jsonString))
    if len(mydata.Groups)  == 0 {
      return true, fmt.Errorf("Does not have acces to create pipeline in the application %s \n",application)
    }
   
    return false, nil
}
func removeDuplicates(elements []string) []string {
    // Use map to record duplicates as we find them.
    encountered := map[string]bool{}
    result := []string{}

    for v := range elements {
        if encountered[elements[v]] == true {
            // Do not add duplicate.
        } else {
            // Record this element as an encountered element.
            encountered[elements[v]] = true
            // Append to result slice.
            result = append(result, elements[v])
        }
    }
    // Return the new slice.
    return result
}
//func pipelienPlanExecution(pipeline interface{}) (interface{},error) { 
	
	 
//}
	

