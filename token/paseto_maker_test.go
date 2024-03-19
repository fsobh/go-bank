package token

import (
	"github.com/fsobh/simplebank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPasetoMaker(t *testing.T) {

	//create a new maker, PASSING IN A RANDOM KEY
	maker, err := NewPasetoMaker(util.RandomString(32))

	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)

}
func TestExpiredPasetoToken(t *testing.T) {

	//create a new maker, PASSING IN A RANDOM KEY
	maker, err := NewPasetoMaker(util.RandomString(32))

	require.NoError(t, err)

	//create a token, passing in a negative duration
	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	//check if the token is valid and if we can decode the token into a payload
	payload, err := maker.VerifyToken(token)
	//Ensure the verification throws an error
	require.Error(t, err)
	//eEnsure the error thrown is of Expired token
	require.EqualError(t, err, ErrExpiredToken.Error())
	//Ensure the payload returned was nil
	require.Nil(t, payload)

}
