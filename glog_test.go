package glog

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLevel(t *testing.T) {
	buff := bytes.NewBuffer(nil)
	l := NewLogger(buff, LevelDebug|LevelInfo)
	require.NoError(t, l.write(LevelInfo, "Hello %q", "George"))
	require.Equal(t, "[INFO] Hello \"George\"\n", buff.String())
	buff.Reset()
	require.NoError(t, l.write(LevelInfo|LevelDebug, "Hello %q", "George"))
	require.Equal(t, "[DEBUG|INFO] Hello \"George\"\n", buff.String())

	buff.Reset()
	require.NoError(t, l.write(LevelDebug, "Hello %q", "George"))
	require.Equal(t, "[DEBUG] Hello \"George\"\n", buff.String())

	l.SetPrefix("K8s")

	buff.Reset()
	require.NoError(t, l.write(LevelDebug, "Hello %q", "George"))
	require.Equal(t, "[DEBUG|K8s] Hello \"George\"\n", buff.String())

	buff.Reset()
	require.NoError(t, l.write(LevelError, "Hello %q", "George"))
	require.Equal(t, "", buff.String())

	const COMETS = 0x10000000
	RegisterLevel(COMETS, "COMETS")

	buff.Reset()
	require.NoError(t, l.write(COMETS, "Hello %q", "George"))
	require.Equal(t, "", buff.String())

	l.SetLevel(COMETS)
	buff.Reset()
	require.NoError(t, l.write(COMETS, "Hello %q", "George"))
	require.Equal(t, "[COMETS|K8s] Hello \"George\"\n", buff.String())

	l.SetLevel(COMETS | LevelDebug)
	buff.Reset()
	require.NoError(t, l.write(COMETS, "Hello %q", "George"))
	require.Equal(t, "[COMETS|K8s] Hello \"George\"\n", buff.String())

	buff.Reset()
	require.NoError(t, l.write(COMETS|LevelDebug, "Hello %q", "George"))
	require.Equal(t, "[COMETS|DEBUG|K8s] Hello \"George\"\n", buff.String())

	l.SetLevel(COMETS)
	buff.Reset()
	require.NoError(t, l.write(COMETS|LevelDebug, "Hello %q", "George"))
	require.Equal(t, "[COMETS|K8s] Hello \"George\"\n", buff.String())

	UnRegisterLevel(COMETS)
	buff.Reset()
	require.NoError(t, l.write(COMETS|LevelDebug, "Hello %q", "George"))
	require.Equal(t, "", buff.String())

	l.SetLevel(COMETS | LevelDebug)
	buff.Reset()
	require.NoError(t, l.write(COMETS|LevelDebug, "Hello %q", "George"))
	require.Equal(t, "[DEBUG|K8s] Hello \"George\"\n", buff.String())

	UnRegisterLevel(LevelDebug)
	buff.Reset()
	require.NoError(t, l.write(COMETS|LevelDebug, "Hello %q", "George"))
	require.Equal(t, "", buff.String())

	RegisterLevel(COMETS, "COMETS")
	buff.Reset()
	require.NoError(t, l.write(COMETS|LevelDebug, "Hello %q", "George"))
	require.Equal(t, "[COMETS|K8s] Hello \"George\"\n", buff.String())
}
