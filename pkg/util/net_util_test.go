package utill

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIp4ViaUDP(t *testing.T) {
	ipv4, err := GetIp4ViaUDP()
	assert.Nil(t, err)
	assert.NotEqual(t, "", ipv4)
	t.Log("Local ipv4 address:", ipv4)

}
