package utils

import (
	"reflect"
	"testing"
)

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "no duplicates",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "some duplicates",
			input:    []string{"a", "b", "a", "c", "b", "d"},
			expected: []string{"a", "b", "c", "d"},
		},
		{
			name:     "all duplicates",
			input:    []string{"one", "one", "one"},
			expected: []string{"one"},
		},
		{
			name:     "order preservation",
			input:    []string{"c", "b", "a", "c", "b"},
			expected: []string{"c", "b", "a"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RemoveDuplicates(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("RemoveDuplicates(%v) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestCheckMissingKeys(t *testing.T) {
	t.Run("some missing keys", func(t *testing.T) {
		m := map[string]interface{}{
			"apple":  1,
			"banana": 2,
		}
		keys := []string{"apple", "banana", "cherry", "date"}
		err := CheckMissingKeys(m, keys)
		if err == nil {
			t.Error("expected an error for missing keys, got nil")
		} else {
			expectedErr := "missing keys: cherry, date"
			if err.Error() != expectedErr {
				t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
			}
		}
	})

	t.Run("no missing keys", func(t *testing.T) {
		m := map[string]interface{}{
			"apple":  1,
			"banana": 2,
			"cherry": 3,
		}
		keys := []string{"apple", "banana", "cherry"}
		err := CheckMissingKeys(m, keys)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("all keys present - simple", func(t *testing.T) {
		m := map[string]interface{}{
			"apple":  1,
			"banana": 2,
			"cherry": 3,
		}
		keys := []string{"apple", "banana", "cherry"}
		err := CheckMissingKeys(m, keys)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("nested keys present", func(t *testing.T) {
		m := map[string]interface{}{
			"project": map[string]interface{}{
				"name":    "test",
				"version": "1.0.0",
			},
			"maven": map[string]interface{}{
				"groupId": "com.example",
			},
		}
		keys := []string{"project.name", "project.version", "maven.groupId"}
		err := CheckMissingKeys(m, keys)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("missing nested keys", func(t *testing.T) {
		m := map[string]interface{}{
			"project": map[string]interface{}{
				"name": "test",
			},
		}
		keys := []string{"project.name", "project.version", "maven.groupId"}
		err := CheckMissingKeys(m, keys)
		if err == nil {
			t.Error("expected an error for missing nested keys, got nil")
		} else {
			expectedErr := "missing keys: project.version, maven.groupId"
			if err.Error() != expectedErr {
				t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
			}
		}
	})

	t.Run("yaml style map", func(t *testing.T) {
		m := map[string]interface{}{
			"teamcity": map[interface{}]interface{}{
				"repositoryName": "test-repo",
				"projectName":    "test-project",
			},
		}
		keys := []string{"teamcity.repositoryName", "teamcity.projectName"}
		err := CheckMissingKeys(m, keys)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("empty key list", func(t *testing.T) {
		m := map[string]interface{}{
			"apple": 1,
		}
		keys := []string{}
		err := CheckMissingKeys(m, keys)
		if err != nil {
			t.Errorf("expected no error for empty key list, got %v", err)
		}
	})

	t.Run("empty map", func(t *testing.T) {
		m := map[string]interface{}{}
		keys := []string{"a", "b"}
		err := CheckMissingKeys(m, keys)
		if err == nil {
			t.Error("expected an error for missing keys, got nil")
		} else {
			expectedErr := "missing keys: a, b"
			if err.Error() != expectedErr {
				t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
			}
		}
	})
}

func TestHasNestedKey(t *testing.T) {
	tests := []struct {
		name     string
		m        map[string]interface{}
		path     []string
		expected bool
	}{
		{
			name: "simple key exists",
			m: map[string]interface{}{
				"key": "value",
			},
			path:     []string{"key"},
			expected: true,
		},
		{
			name: "nested key exists",
			m: map[string]interface{}{
				"parent": map[string]interface{}{
					"child": "value",
				},
			},
			path:     []string{"parent", "child"},
			expected: true,
		},
		{
			name: "yaml style map",
			m: map[string]interface{}{
				"config": map[interface{}]interface{}{
					"setting": "value",
				},
			},
			path:     []string{"config", "setting"},
			expected: true,
		},
		{
			name: "key does not exist",
			m: map[string]interface{}{
				"key": "value",
			},
			path:     []string{"nonexistent"},
			expected: false,
		},
		{
			name: "nested key does not exist",
			m: map[string]interface{}{
				"parent": map[string]interface{}{},
			},
			path:     []string{"parent", "child"},
			expected: false,
		},
		{
			name:     "empty path",
			m:        map[string]interface{}{},
			path:     []string{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasNestedKey(tt.m, tt.path)
			if got != tt.expected {
				t.Errorf("hasNestedKey() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestApplyOverrides(t *testing.T) {
	tests := []struct {
		name      string
		initial   map[string]interface{}
		overrides []string
		expected  map[string]interface{}
	}{
		{
			name:      "simple key-value",
			initial:   map[string]interface{}{},
			overrides: []string{"name=John"},
			expected: map[string]interface{}{
				"name": "John",
			},
		},
		{
			name: "nested keys",
			initial: map[string]interface{}{
				"existing": "value",
			},
			overrides: []string{"config.port=8080", "config.host=localhost"},
			expected: map[string]interface{}{
				"existing": "value",
				"config": map[string]interface{}{
					"port": "8080",
					"host": "localhost",
				},
			},
		},
		{
			name: "override existing value",
			initial: map[string]interface{}{
				"name": "Old",
			},
			overrides: []string{"name=New"},
			expected: map[string]interface{}{
				"name": "New",
			},
		},
		{
			name:    "deep nesting",
			initial: map[string]interface{}{},
			overrides: []string{
				"a.b.c.d=value",
				"a.b.x=1",
				"a.y=2",
			},
			expected: map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{
						"c": map[string]interface{}{
							"d": "value",
						},
						"x": "1",
					},
					"y": "2",
				},
			},
		},
		{
			name:      "invalid format ignored",
			initial:   map[string]interface{}{},
			overrides: []string{"invalid", "name=John", "=value", "key="},
			expected: map[string]interface{}{
				"name": "John",
				"key":  "",
				"":     "value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ApplyOverrides(tt.initial, tt.overrides)
			if !reflect.DeepEqual(tt.initial, tt.expected) {
				t.Errorf("ApplyOverrides() = %v; want %v", tt.initial, tt.expected)
			}
		})
	}
}
