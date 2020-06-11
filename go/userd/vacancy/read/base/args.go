package base

import "github.com/GDVFox/tenjin/utils/server"

// Arguments represents arguments for person reading
type Arguments struct {
	ID               []int64  `json:"id"`
	DepartmentNames  []string `json:"department_name"`
	AppointmentNames []string `json:"appointment_name"`
}

// Validate checks argument correct
func (a *Arguments) Validate() error {
	return nil
}

// NewArguments creates new ReadPersonArguments instanse
func NewArguments() server.Arguments {
	return &Arguments{}
}
