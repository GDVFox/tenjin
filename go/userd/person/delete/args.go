package delete

import "github.com/GDVFox/tenjin/utils/server"

// Arguments represents arguments for person deleting
type Arguments struct {
	ID int64 `json:"id"`
}

// Validate checks argument correct
func (a *Arguments) Validate() error {
	return nil
}

// NewArguments creates new DeletePersonArguments instanse
func NewArguments() server.Arguments {
	return &Arguments{}
}
