package pipeline

import (
  
)
type PipelineConfig struct {
	PipelineName string `yaml:"pipelinename"`
	Application string `yaml:"application"`
	TemplateReference string  `yaml:"pipelinetemplatename"`
    Variables interface{}  `yaml:"variables"`
    Triggers []interface{} `yaml:"triggers"`
    ExpectedArtifacts []interface{} `yaml:"expectedArtifacts"`   
    Notifications []interface{} `yaml:"notifications"`
    Description string `yaml:"description"`
    Stages  []interface{} `yaml:"stages"`
    Parameters  []interface{} `yaml:"parameters"`
}