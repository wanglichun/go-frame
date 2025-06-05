// internal/config/registry.go
package config

import (
	"fmt"
	"reflect"
	"sync"
)

// Registry 组件注册表
type Registry struct {
	components map[string]reflect.Type
	mu         sync.RWMutex
}

// NewRegistry 创建新的注册表
func NewRegistry() *Registry {
	return &Registry{
		components: make(map[string]reflect.Type),
	}
}

// Register 注册组件类型
func (r *Registry) Register(name string, componentType interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()

	t := reflect.TypeOf(componentType)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	r.components[name] = t
}

// CreateComponent 根据配置创建组件实例
func (r *Registry) CreateComponent(config ComponentCfg) (interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	componentType, exists := r.components[config.Type]
	if !exists {
		return nil, fmt.Errorf("未注册的组件类型: %s", config.Type)
	}

	// 创建实例
	instance := reflect.New(componentType).Elem()

	// 设置参数
	if config.Params != nil {
		for paramName, paramValue := range config.Params {
			field := instance.FieldByName(paramName)
			if !field.IsValid() {
				continue
			}

			if !field.CanSet() {
				continue
			}

			// 类型转换
			value := reflect.ValueOf(paramValue)
			if value.Type().ConvertibleTo(field.Type()) {
				field.Set(value.Convert(field.Type()))
			}
		}
	}

	return instance.Interface(), nil
}