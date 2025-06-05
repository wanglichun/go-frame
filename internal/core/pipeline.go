// internal/core/pipeline.go
package core

import (
	"sync"
)

// Pipeline 管理整个处理流程
type Pipeline struct {
	Name   string
	Stages []*Stage
	mu     sync.RWMutex
}

// NewPipeline 创建新的管道
func NewPipeline(name string) *Pipeline {
	return &Pipeline{
		Name:   name,
		Stages: make([]*Stage, 0),
	}
}

// AddStage 添加阶段到管道
func (p *Pipeline) AddStage(stage *Stage) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Stages = append(p.Stages, stage)
}

// Execute 执行整个管道
func (p *Pipeline) Execute(ctx *Context) error {
	for _, stage := range p.Stages {
		if ctx.HasError() {
			return ctx.GetFirstError()
		}

		executor := NewExecutor(stage)
		if err := executor.Execute(ctx); err != nil {
			ctx.AddError(err)
			return err
		}
	}
	return nil
}