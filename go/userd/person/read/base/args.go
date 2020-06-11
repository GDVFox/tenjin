package base

import "github.com/GDVFox/tenjin/utils/server"

// Arguments represents arguments for person reading
type Arguments struct {
	ID            []int64 `json:"id"`
	EmployeesOnly bool    `json:"employees_only"`
}

// Validate checks argument correct
func (a *Arguments) Validate() error {
	return nil
}

// NewArguments creates new ReadPersonArguments instanse
func NewArguments() server.Arguments {
	return &Arguments{EmployeesOnly: false}
}
