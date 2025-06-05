// internal/core/context.go
package core

import (
	"context"
	"sync"
)

// Context 包含请求处理过程中的所有信息
type Context struct {
	context.Context
	Request    interface{}
	Params     map[string]interface{}
	Data       map[string]interface{}
	Result     interface{}
	Errors     []error
	mu         sync.RWMutex
	cancelFunc context.CancelFunc
}

// NewContext 创建新的上下文
func NewContext(ctx context.Context, req interface{}) *Context {
	ctx, cancel := context.WithCancel(ctx)
	return &Context{
		Context:    ctx,
		Request:    req,
		Params:     make(map[string]interface{}),
		Data:       make(map[string]interface{}),
		cancelFunc: cancel,
	}
}

// AddError 添加错误到上下文中
func (c *Context) AddError(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Errors = append(c.Errors, err)
}

// HasError 检查上下文中是否有错误
func (c *Context) HasError() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.Errors) > 0
}

// GetFirstError 获取第一个错误
func (c *Context) GetFirstError() error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.Errors) == 0 {
		return nil
	}
	return c.Errors[0]
}