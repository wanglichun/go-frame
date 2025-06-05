// internal/config/config.go
package config

import (
	"fmt"
)

// Config 定义配置结构
type Config struct {
	Name   string      `json:"name"`
	Stages []*StageCfg `json:"stages"`
}

// StageCfg 定义阶段配置
type StageCfg struct {
	Name         string         `json:"name"`
	ExecutionMode string         `json:"executionMode"`
	Components   []ComponentCfg `json:"components"`
}

// ComponentCfg 定义组件配置
type ComponentCfg struct {
	Type         string                 `json:"type"`
	Name         string                 `json:"name"`
	Params       map[string]interface{} `json:"params,omitempty"`
	Dependencies []string               `json:"dependencies,omitempty"`
	Requires     []string               `json:"requires,omitempty"`
	Provides     []string               `json:"provides,omitempty"`
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("管道名称不能为空")
	}

	if len(c.Stages) == 0 {
		return fmt.Errorf("管道必须至少包含一个阶段")
	}

	for _, stage := range c.Stages {
		if stage.Name == "" {
			return fmt.Errorf("阶段名称不能为空")
		}

		if len(stage.Components) == 0 {
			return fmt.Errorf("阶段 %s 必须至少包含一个组件", stage.Name)
		}

		// 验证ExecutionMode
		mode := stage.ExecutionMode
		if mode == "" {
			stage.ExecutionMode = "parallel"
		} else if mode != "parallel" && mode != "sequential" && mode != "dependency" {
			return fmt.Errorf("阶段 %s 的执行模式 %s 无效", stage.Name, mode)
		}

		// 验证组件名称唯一
		names := make(map[string]bool)
		for _, comp := range stage.Components {
			if comp.Name == "" {
				return fmt.Errorf("阶段 %s 中的组件名称不能为空", stage.Name)
			}

			if names[comp.Name] {
				return fmt.Errorf("阶段 %s 中的组件名称 %s 重复", stage.Name, comp.Name)
			}

			names[comp.Name] = true
		}
	}

	return nil
}