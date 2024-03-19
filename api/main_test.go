package api

import (
	"crypto/ed25519"
	"crypto/rand"
	db "github.com/fsobh/simplebank/db/sqlc"
	"github.com/fsobh/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)

	require.NoError(t, err)

	config := util.Config{
		PasetoPublicKey:     publicKey,
		PasetoPrivateKey:    privateKey,
		AccessTokenDuration: time.Minute,
		PasetoSymmetricKey:  util.RandomString(32),
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)
	return server
}

// WE ONLY CREATED THIS FILE SO THE LOGS LOOK PRETTIER IN CONSOLE.
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
