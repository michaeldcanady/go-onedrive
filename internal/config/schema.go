package config

import (
	_ "embed"
	"encoding/json"
	"sort"
	"strings"
)

//go:embed schema.json
var schemaJSON []byte

// SchemaProperty represents a property in the configuration schema.
type SchemaProperty struct {
	Type        string                    `json:"type"`
	Enum        []string                  `json:"enum,omitempty"`
	AnyOf       []schemaAnyOf             `json:"anyOf,omitempty"`
	Properties  map[string]SchemaProperty `json:"properties,omitempty"`
}

type schemaAnyOf struct {
	Const string `json:"const,omitempty"`
	Type  string `json:"type,omitempty"`
}

type schemaRoot struct {
	Properties map[string]SchemaProperty `json:"properties"`
}

// IsStrictEnum returns true if the key must strictly match one of the allowed values.
func IsStrictEnum(key string) bool {
	prop := getProperty(key)
	if prop == nil {
		return false
	}

	// If Enum is present, it's usually strict.
	if len(prop.Enum) > 0 {
		return true
	}

	// If AnyOf is present, it's strict ONLY if all branches are constants.
	if len(prop.AnyOf) > 0 {
		for _, ao := range prop.AnyOf {
			if ao.Const == "" && ao.Type != "" {
				// We found a branch that allows any value of a certain type.
				return false
			}
		}
		return true
	}

	return false
}

func getProperty(key string) *SchemaProperty {
	var root schemaRoot
	if err := json.Unmarshal(schemaJSON, &root); err != nil {
		return nil
	}

	parts := strings.Split(key, ".")
	currentProp := SchemaProperty{Properties: root.Properties}

	for _, part := range parts {
		prop, ok := currentProp.Properties[part]
		if !ok {
			return nil
		}
		currentProp = prop
	}
	return &currentProp
}

// GetAvailableKeys returns all valid configuration keys in dotted format.
func GetAvailableKeys() []string {
	var root schemaRoot
	if err := json.Unmarshal(schemaJSON, &root); err != nil {
		return nil
	}

	var keys []string
	for k, p := range root.Properties {
		keys = append(keys, getKeysRecursive(k, p)...)
	}
	sort.Strings(keys)
	return keys
}

func getKeysRecursive(prefix string, prop SchemaProperty) []string {
	if len(prop.Properties) == 0 {
		return []string{prefix}
	}

	var keys []string
	for k, p := range prop.Properties {
		keys = append(keys, getKeysRecursive(prefix+"."+k, p)...)
	}
	return keys
}

// GetAllowedValues returns the enum or const values for a given key.
func GetAllowedValues(key string) []string {
	prop := getProperty(key)
	if prop == nil {
		return nil
	}

	var values []string
	values = append(values, prop.Enum...)
	for _, ao := range prop.AnyOf {
		if ao.Const != "" {
			values = append(values, ao.Const)
		}
	}
	sort.Strings(values)
	return values
}
