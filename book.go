package audbk

import "github.com/spf13/cast"

type Meta struct {
	Fields map[string]any
}

func NewMeta() *Meta {
	return &Meta{
		Fields: make(map[string]any),
	}
}

func (m *Meta) Set(key string, val any) *Meta {
	m.Fields[key] = val
	return m
}

func (m *Meta) Del(key string) *Meta {
	delete(m.Fields, key)
	return m
}

func (m *Meta) Has(key string) bool {
	if _, ok := m.Fields[key]; ok {
		return true
	}
	return false
}

func (m *Meta) Get(key string) any {
	if m.Has(key) {
		return m.Fields[key]
	}
	return nil
}

func (m *Meta) GetString(key string) string {
	val := m.Get(key)
	if val != nil {
		return cast.ToString(val)
	}
	return ""
}

func (m *Meta) GetSlice(key string) []any {
	val := m.Get(key)
	if val != nil {
		return cast.ToSlice(val)
	}
	return []any{val}
}

func (m *Meta) GetStringMap(key string) map[string]any {
	val := m.Get(key)
	if val != nil {
		return cast.ToStringMap(val)
	}
	return map[string]any{key: val}
}
