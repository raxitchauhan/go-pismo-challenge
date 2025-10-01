//go:build tools

package pkg

// this file exists to forcefully vendor dependencies used by go:generate

import (
	_ "go.uber.org/mock/mockgen"
)
