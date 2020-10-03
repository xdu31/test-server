package service

import (
	"encoding/json"
)

// MarshalJSON overloads Ip's standard MarshalJSON
func (m *Ip) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		IpAddress string `json:"ip_address,omitempty"`
	}{
		IpAddress: m.IpAddress.Value,
	})
}
