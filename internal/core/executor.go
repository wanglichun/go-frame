// internal/core/executor.go
package core

import (
	"fmt"
	"sync"
)

// Executor 阶段执行器
type Executor struct {
	stage *Stage
}

// NewExecutor 创建新的执行器
func NewExecutor(stage *Stage) *Executor {
	return &Executor{
		stage: stage,
	}
}

// Execute 执行阶段
func (e *Executor) Execute(ctx *Context) error {
	switch e.stage.ExecutionMode {
	case ParallelMode:
		return e.executeParallel(ctx)
	case SequentialMode:
		return e.executeSequential(ctx)
	case DependencyMode:
		return e.executeWithDependencies(ctx)
	default:
		return fmt.Errorf("未知的执行模式: %s", e.stage.ExecutionMode)
	}
}

// executeParallel 并行执行所有组件
func (e *Executor) executeParallel(ctx *Context) error {
	var wg sync.WaitGroup
	errorCh := make(chan error, len(e.stage.Components))

	for name, component := range e.stage.Components {
		wg.Add(1)
		go func(c Component, name string) {
			defer wg.Done()
			if err := c.Execute(ctx); err != nil {
				errorCh <- fmt.Errorf("组件 %s 执行失败: %w", name, err)
				ctx.AddError(err)
				ctx.cancelFunc()
			}
		}(component, name)
	}

	wg.Wait()
	close(errorCh)

	if len(ctx.Errors) > 0 {
		return fmt.Errorf("阶段 %s 执行失败: %w", e.stage.Name, ctx.GetFirstError())
	}

	return nil
}

// executeSequential 按顺序执行组件
func (e *Executor) executeSequential(ctx *Context) error {
	for _, name := range e.stage.order {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			component, exists := e.stage.Components[name]
			if !exists {
				return fmt.Errorf("组件 %s 不存在", name)
			}

			if err := component.Execute(ctx); err != nil {
				return fmt.Errorf("组件 %s 执行失败: %w", name, err)
			}
		}
	}

	return nil
}

// executeWithDependencies 按依赖关系执行组件
func (e *Executor) executeWithDependencies(ctx *Context) error {
	// 这里实现依赖关系处理逻辑
	// 简化版：假设顺序已在stage.order中设置好
	return e.executeSequential(ctx)
}