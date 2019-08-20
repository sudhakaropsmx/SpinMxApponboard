package pipeline

import (
  
)
type PipelineConfig struct {
	PipelineName string `yaml:"pipelinename"`
	Application string `yaml:"application"`
	TemplateReference string  `yaml:"pipelinetemplatename"`
    Variables map[string]interface{}  `yaml:"variables"`
    Triggers []map[string]interface{} `yaml:"triggers"`
    ExpectedArtifacts []map[string]interface{} `yaml:"expectedArtifacts"`   
    Notifications []map[string]interface{} `yaml:"notifications"`
    Description string `yaml:"description"`
    Stages  []map[string]interface{} `yaml:"stages"`
    Parameters  []map[string]interface{} `yaml:"parameters"`
}