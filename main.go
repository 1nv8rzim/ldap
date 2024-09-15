package ldap

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-ldap/ldap/v3"
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
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	var connStr string
	if schema.LDAPS {
		connStr = fmt.Sprintf("ldaps://%s:%d?tls=1&insecure=1", schema.Target, schema.Port)
	} else {
		connStr = fmt.Sprintf("ldap://%s:%d", schema.Target, schema.Port)
	}

	conn, err := ldap.DialURL(connStr)
	if err != nil {
		return err
	}
	defer conn.Close()

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context deadline is not set")
	}

	conn.SetTimeout(time.Until(deadline))
	err = conn.Bind(schema.Username, schema.Password)
	if err != nil {
		return fmt.Errorf("failed to bind: %w", err)
	}

	return nil
}
