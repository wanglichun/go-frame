// internal/core/stage.go
package core

import (
	"fmt"
)

// ExecutionMode 阶段执行模式
type ExecutionMode string

const (
	ParallelMode   ExecutionMode = "parallel"   // 所有组件并发执行
	SequentialMode ExecutionMode = "sequential" // 组件按注册顺序执行
	DependencyMode ExecutionMode = "dependency" // 基于依赖关系执行
)

// Stage 表示处理流程中的一个阶段
type Stage struct {
	Name         string
	Components   map[string]Component
	ExecutionMode ExecutionMode
	order        []string // 执行顺序
}

// NewStage 创建新的阶段
func NewStage(name string, mode ExecutionMode) *Stage {
	return &Stage{
		Name:         name,
		Components:   make(map[string]Component),
		ExecutionMode: mode,
	}
}

// AddComponent 添加组件到阶段
func (s *Stage) AddComponent(name string, component Component) {
	s.Components[name] = component
	s.order = append(s.order, name)
}

// SetExecutionOrder 设置组件执行顺序
func (s *Stage) SetExecutionOrder(order []string) error {
	// 验证顺序
	for _, name := range order {
		if _, exists := s.Components[name]; !exists {
			return fmt.Errorf("组件 %s 不存在", name)
		}
	}

	s.order = order
	return nil
}