package api

import (
	"github.com/gin-gonic/gin"
	"os"
	"testing"
)

// WE ONLY CREATED THIS FILE SO THE LOGS LOOK PRETTIER IN CONSOLE.
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
