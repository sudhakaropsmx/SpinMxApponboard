package application

import (

)

type Permissions  struct {
	   Read []string 
	   Write []string 
	   Execute []string 
    } 
type AplicationConfig struct {
	Name string 
	Aliases string 
	Account string 
	CloudProviders string  
	Description string  
	Email string  
    Permissions   struct {
	   Read []string 
	   Write []string 
	   Execute []string 
    } 
    TrafficGuards interface{} 
    RepoSlug string 
    RepoType string 
    RepoProjectKey string 
    LastModifiedBy string 
}

type AplicationOpsConfig struct {
	Application []struct {
	    Name string `yml:"name"`
	    OwnerEmail string `yml:"owneremail"`
	    Permissions string `yml:"permissions"`
	    Accounts string `yml:"accounts"`
	} `yml:"application"`
}