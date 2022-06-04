package uilist

import "fmt"

type ListItem struct {
	KeyType string
	Key     string
	Err     string
}

func (i ListItem) Title() string { return i.Key }

func (i ListItem) Description() string {
	if i.Err != "" {
		return "got an error"
	}
	return fmt.Sprintf("key: %d bytes, value: %d bytes", len(i.Key))
}

func (i ListItem) FilterValue() string { return i.Key }
