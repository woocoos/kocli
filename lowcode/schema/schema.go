package schema

import (
	"bytes"
	"encoding/json"
)

const (
	// RootSchemaNameBlock 区块组件
	RootSchemaNameBlock = "Block"
	// RootSchemaNamePage 页面组件
	RootSchemaNamePage = "Page"
	// RootSchemaNameComponent 低代码业务组件
	RootSchemaNameComponent = "Component"
)

// RootSchema 根业务组件描述
type RootSchema interface {
	// GetRootSchemaName 获取根组件名称
	GetRootSchemaName() string
	GetID() string
}

// ProjectSchema 应用描述
type ProjectSchema struct {
	ID string `json:"id,omitempty"`
	// Version 当前应用协议版本号
	Version string `json:"version"`
	// ComponentsMap 当前应用所有组件映射关系
	ComponentsMap    []Component     `json:"-"`
	ComponentsMapRaw json.RawMessage `json:"componentsMap"`
	// 描述应用所有页面、低代码组件的组件树
	ComponentsTree    []RootSchema    `json:"-"`
	ComponentsTreeRaw json.RawMessage `json:"componentsTree"`
}

// ParseProjectSchema 通过json字符串解析
func ParseProjectSchema(bs []byte) (*ProjectSchema, error) {
	var schema ProjectSchema
	err := json.Unmarshal(bs, &schema)
	if err != nil {
		return nil, err
	}
	if err = schema.parseComponentsMap(); err != nil {
		return nil, err
	}
	if err = schema.parseComponentTree(); err != nil {
		return nil, err
	}

	return &schema, nil
}

func (ps *ProjectSchema) parseComponentsMap() error {
	if ps.ComponentsMapRaw == nil {
		return nil
	}
	var tc []json.RawMessage
	err := json.Unmarshal(ps.ComponentsMapRaw, &tc)
	if err != nil {
		return err
	}
	for _, c := range tc {
		if bytes.Contains(c, []byte(`"devMode"`)) {
			var cm LowCodeComponent
			err = json.Unmarshal(c, &cm)
			if err != nil {
				return err
			}
			ps.ComponentsMap = append(ps.ComponentsMap, &cm)
		} else {
			var cm ProCodeComponent
			err = json.Unmarshal(c, &cm)
			if err != nil {
				return err
			}
			ps.ComponentsMap = append(ps.ComponentsMap, &cm)
		}
	}
	return nil
}

func (ps *ProjectSchema) parseComponentTree() error {
	var tc []json.RawMessage
	err := json.Unmarshal(ps.ComponentsTreeRaw, &tc)
	if err != nil {
		return err
	}
	for _, c := range tc {
		ct := struct {
			Type string `json:"componentName"`
		}{}
		if err = json.Unmarshal(c, &ct); err != nil {
			return err
		}
		switch ct.Type {
		case RootSchemaNamePage:
			var cm = NewPageSchema(ps)
			if err = json.Unmarshal(c, cm); err != nil {
				return err
			}
			ps.ComponentsTree = append(ps.ComponentsTree, cm)
		case RootSchemaNameBlock:
			var cm BlockSchema
			if err = json.Unmarshal(c, &cm); err != nil {
				return err
			}
			ps.ComponentsTree = append(ps.ComponentsTree, &cm)
		case RootSchemaNameComponent:
			var cm ComponentSchema
			if err = json.Unmarshal(c, &cm); err != nil {
				return err
			}
			ps.ComponentsTree = append(ps.ComponentsTree, &cm)
		}
	}
	return nil
}

// FindComponentMap 根据组件名称查找业务组件
func (ps *ProjectSchema) FindComponentMap(component string) *ProCodeComponent {
	for _, cm := range ps.ComponentsMap {
		v, ok := cm.(*ProCodeComponent)
		if !ok {
			continue
		}
		if cm.GetComponentName() == component {
			return v
		}
	}
	return nil
}
