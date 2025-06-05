// internal/config/parser.go
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/yourusername/pipeline-framework/internal/core"
)

// Parser 配置解析器
type Parser struct {
	registry *Registry
}

// NewParser 创建新的解析器
func NewParser(registry *Registry) *Parser {
	return &Parser{
		registry: registry,
	}
}

// ParseFromFile 从文件解析配置
func (p *Parser) ParseFromFile(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return p.ParseFromJSON(data)
}

// ParseFromJSON 从JSON数据解析配置
func (p *Parser) ParseFromJSON(data []byte) (*Config, error) {
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

// BuildPipeline 根据配置构建Pipeline
func (p *Parser) BuildPipeline(config *Config) (*core.Pipeline, error) {
	pipeline := core.NewPipeline(config.Name)

	for _, stageCfg := range config.Stages {
		stage := core.NewStage(stageCfg.Name, core.ExecutionMode(stageCfg.ExecutionMode))

		for _, compCfg := range stageCfg.Components {
			component, err := p.registry.CreateComponent(compCfg)
			if err != nil {
				return nil, fmt.Errorf("创建组件 %s 失败: %w", compCfg.Name, err)
			}

			// 类型断言
			coreComponent, ok := component.(core.Component)
			if !ok {
				return nil, fmt.Errorf("组件 %s 未实现 core.Component 接口", compCfg.Name)
			}

			// 添加依赖信息
			if depComponent, ok := coreComponent.(core.DependencyAwareComponent); ok {
				depComponent = &dependencyWrapper{
					Component:  coreComponent,
					dependencies: compCfg.Dependencies,
				}
				coreComponent = depComponent
			}

			// 添加数据依赖信息
			if dataComponent, ok := coreComponent.(core.DataAwareComponent); ok {
				dataComponent = &dataWrapper{
					Component:    coreComponent,
					requiredData: compCfg.Requires,
					providedData: compCfg.Provides,
				}
				coreComponent = dataComponent
			}

			stage.AddComponent(compCfg.Name, coreComponent)
		}

		// 设置执行顺序（如果是dependency模式）
		if stage.ExecutionMode == core.DependencyMode {
			// 这里应该实现依赖解析逻辑
			// 简化版：按配置文件中的顺序执行
			order := make([]string, len(stageCfg.Components))
			for i, comp := range stageCfg.Components {
				order[i] = comp.Name
			}

			if err := stage.SetExecutionOrder(order); err != nil {
				return nil, err
			}
		}

		pipeline.AddStage(stage)
	}

	return pipeline, nil
}

// dependencyWrapper 实现DependencyAwareComponent接口
type dependencyWrapper struct {
	core.Component
	dependencies []string
}

// GetDependencies 实现DependencyAwareComponent接口
func (d *dependencyWrapper) GetDependencies() []string {
	return d.dependencies
}

// dataWrapper 实现DataAwareComponent接口
type dataWrapper struct {
	core.Component
	requiredData []string
	providedData []string
}

// GetRequiredData 实现DataAwareComponent接口
func (d *dataWrapper) GetRequiredData() []string {
	return d.requiredData
}

// GetProvidedData 实现DataAwareComponent接口
func (d *dataWrapper) GetProvidedData() []string {
	return d.providedData
}