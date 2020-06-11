package delete

import "github.com/GDVFox/tenjin/utils/server"

// Arguments represents arguments for task deleting
type Arguments struct {
	ID int64 `json:"id"`
}

// Validate checks argument correct
func (a *Arguments) Validate() error {
	return nil
}

// NewArguments creates arguments for task delete
func NewArguments() server.Arguments {
	return &Arguments{}
}
