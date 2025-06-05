// internal/core/component.go
package core

// Component 定义可执行的组件接口
type Component interface {
	Execute(ctx *Context) error
}

// ComponentFunc 函数适配器，允许普通函数作为组件
type ComponentFunc func(ctx *Context) error

// Execute 实现Component接口
func (f ComponentFunc) Execute(ctx *Context) error {
	return f(ctx)
}

// DependencyAwareComponent 支持依赖声明的组件
type DependencyAwareComponent interface {
	Component
	GetDependencies() []string
}

// DataAwareComponent 支持数据依赖的组件
type DataAwareComponent interface {
	Component
	GetRequiredData() []string
	GetProvidedData() []string
}