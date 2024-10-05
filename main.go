package ldap

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-ldap/ldap/v3"
)

type Schema struct {
	Server   string `key:"server"`
	Port     int    `key:"port" default:"389"`
	Username string `key:"username"`
	Password string `key:"password"`
	LDAPS    bool   `key:"ldaps" default:"false"`
}

func (s *Schema) Validate(config string) error {
	conf := Schema{}

	err := schema.Unmarshal([]byte(config), &conf)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if conf.Server == "" {
		return fmt.Errorf("server is required; got %q", conf.Server)
	}

	if conf.Port <= 0 || conf.Port > 65535 {
		return fmt.Errorf("port is invalid; got %d", conf.Port)
	}

	if conf.Username == "" {
		return fmt.Errorf("username is required; got %q", conf.Username)
	}

	if conf.Password == "" {
		return fmt.Errorf("password is required; got %q", conf.Password)
	}

	return nil
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
