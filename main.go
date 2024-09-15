package ldap

import (
	"context"
	"encoding/json"
)

type Schema struct {
	Target   string `json:"target"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	LDAPS    bool   `json:"ldaps"`
}

func Run(ctx context.Context, config string) error {
	schema := Schema{}

	err := json.Unmarshal([]byte(config), &schema)
	if err != nil {
		return err
	}

	return nil
}
