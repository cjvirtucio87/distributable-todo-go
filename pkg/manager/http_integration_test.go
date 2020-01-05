// +build integration

package manager

import (
	"cjvirtucio87/distributed-todo-go/pkg/config"
	"testing"
)

func TestIntegrationStart(t *testing.T) {
	if m, err := NewHttpManager(
		config.NewViperLoader(
			"app_test",
			"yaml",
		),
	); err != nil {
		t.Fatal(err)
	} else {
		defer m.Stop()

		m.Start()
	}
}
