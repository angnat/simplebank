package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	/*untuk mengubah gin supaya ke mode test*/
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
