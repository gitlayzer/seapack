package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFromEnvs(t *testing.T) {
	env, err := FromEnvs([]string{
		"VAR1=value1",
		"VAR2=value2",
		"SEAPACK_APT_PACKAGES=apt1,apt2",
		"COMMA=this has, a comma",
		"SEAPACK_TRUTHY_CASE=True ",
		"SEAPACK_TRUTHY_INT_CASE= 1 ",
		"HELLO+WORLD=boop",
	})

	require.NoError(t, err)
	require.Equal(t, env.GetVariable("VAR1"), "value1")
	require.Equal(t, env.GetVariable("VAR2"), "value2")
	require.Equal(t, env.GetVariable("SEAPACK_APT_PACKAGES"), "apt1,apt2")
	require.Equal(t, env.GetVariable("COMMA"), "this has, a comma")
	require.Equal(t, env.IsConfigVariableTruthy("TRUTHY_CASE"), true)
	require.Equal(t, env.IsConfigVariableTruthy("TRUTHY_INT_CASE"), true)
	require.Equal(t, env.GetVariable("HELLO+WORLD"), "boop")
}
