package schema

import (
	"encoding/json"
)

// CVType 组件属性值(CompositeValue)类型
type CVType string

const (
	// 默认类型: CompositeValue
	CVTypeDefault CVType = "CompositeValue"
	// CompositeValueMap
	CVTypeMap CVType = "CompositeValueMap"
	// CompositeValueSlice
	CVTypeSlice CVType = "CompositeValueSlice"
	// String
	CVTypeString CVType = "String"
	// Number
	CVTypeNumber CVType = "Number"
	// Boolean
	CVTypeBoolean CVType = "Boolean"
	// JSExpression
	CVTypeJSExpression CVType = "JSExpression"
	// JSFunction
	CVTypeJSFunction CVType = "JSFunction"
	// JSSlot
	CVTypeJSSlot CVType = "JSSlot"
	// ANY
	CVTypeAny CVType = "any"
)

func (cv CVType) String() string {
	return string(cv)
}

// CompositeValueSlice 组件属性值数组
type CompositeValueSlice []CompositeValue

func (cvs CompositeValueSlice) LoadComponentNames() (cns []string) {
	for _, cv := range cvs {
		cns = append(cns, cv.LoadComponentNames()...)
	}
	return
}

// CompositeValueMap 组件属性值映射
type CompositeValueMap map[string]CompositeValue

func (cvm CompositeValueMap) LoadComponentNames() (cns []string) {
	for _, cv := range cvm {
		cns = append(cns, cv.LoadComponentNames()...)
	}
	return
}

func (cvm CompositeValueMap) HasKey(key string) bool {
	_, ok := cvm[key]
	return ok
}

func (cvm CompositeValueMap) Get(key string) CompositeValue {
	return cvm[key]
}

type DOMText string

// NodeDataType implement NodeData
func (dom *DOMText) NodeDataType() string {
	return "DOMText"
}

func parseNodeDataType(bytes []byte) (any, error) {
	t := struct {
		TypeName      string `json:"type"`
		ComponentName string `json:"componentName"`
	}{}
	if err := json.Unmarshal(bytes, &t); err != nil {
		return nil, err
	}
	switch t.TypeName {
	case "JSExpression":
		var jse JSExpression
		if err := json.Unmarshal(bytes, &jse); err != nil {
			return nil, err
		}
		return &jse, nil
	case "JSFunction":
		var jsf JSFunction
		if err := json.Unmarshal(bytes, &jsf); err != nil {
			return nil, err
		}
		return &jsf, nil
	case "JSSlot":
		var jss JSSlot
		if err := json.Unmarshal(bytes, &jss); err != nil {
			return nil, err
		}
		return &jss, nil
	}
	//
	switch t.ComponentName {
	case "Page", "Block":
	default:
		if t.ComponentName == "" {
			break
		}
		var ns NodeSchema
		if err := json.Unmarshal(bytes, &ns); err != nil {
			return nil, err
		}
		return &ns, nil
	}
	var o CompositeValueMap
	if err := json.Unmarshal(bytes, &o); err != nil {
		return nil, err
	}
	return o, nil
}

func parseNodeDataTypeArray(bytes []byte) (any, error) {
	var a CompositeValueSlice
	if err := json.Unmarshal(bytes, &a); err != nil {
		return nil, err
	}
	return a, nil
}

// CompositeValue 属性值类型
//
// 包括 String, Number, Object, Array, Boolean, JSSlot, JSFunction, JSExpression, NodeSchema及其数组类型
type CompositeValue struct {
	Value any
}

func (cv *CompositeValue) UnmarshalJSON(bytes []byte) error {
	switch bytes[0] {
	case '{':
		v, err := parseNodeDataType(bytes)
		if err != nil {
			return err
		}
		cv.Value = v
	case '[':
		v, err := parseNodeDataTypeArray(bytes)
		if err != nil {
			return err
		}
		cv.Value = v
	case '"':
		var s string
		if err := json.Unmarshal(bytes, &s); err != nil {
			return err
		}
		cv.Value = s
	default:
		var n any
		if err := json.Unmarshal(bytes, &n); err != nil {
			return err
		}
		cv.Value = n
	}
	return nil
}

// LoadComponentNames 读取值内的所有组件名称
func (cv *CompositeValue) LoadComponentNames() (cns []string) {
	if cv.Value == nil {
		return nil
	}
	switch val := cv.Value.(type) {
	case CompositeValueSlice:
		for _, sub := range val {
			cns = append(cns, sub.LoadComponentNames()...)
		}
	case CompositeValueMap:
		for _, sub := range val {
			cns = append(cns, sub.LoadComponentNames()...)
		}
	case Container:
		cns = append(cns, val.LoadComponentNames()...)
	default:
		//fmt.Print(val)
	}
	return
}

// TypeString 返回值类型字符串
func (cv CompositeValue) TypeString() string {
	switch cv.Value.(type) {
	case CompositeValueMap:
		return CVTypeMap.String()
	case CompositeValueSlice:
		return CVTypeSlice.String()
	case CompositeValue:
		return CVTypeDefault.String()
	case int, float64:
		return CVTypeNumber.String()
	case bool:
		return CVTypeBoolean.String()
	case string:
		return CVTypeString.String()
	case *JSSlot:
		return CVTypeJSSlot.String()
	case *JSExpression:
		return CVTypeJSExpression.String()
	case *JSFunction:
		return CVTypeJSFunction.String()
	case Component:
		return "Component"
	default:
		return CVTypeAny.String()
	}
}
