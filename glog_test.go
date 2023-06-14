package glog

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLevel(t *testing.T) {
	t.Logf("%b", LevelDebug)
	t.Logf("%b", LevelInfo)
	t.Logf("%b", LevelWarn)
	t.Logf("%b", LevelError)
	t.Logf("%b", LevelFatal)

	buff := bytes.NewBuffer(nil)
	l := NewLogger(buff, LevelDebug|LevelInfo)
	require.NoError(t, l.writef(LevelInfo, "Hello %q", "George"))
	require.Equal(t, "[INFO] Hello \"George\"\n", buff.String())
	buff.Reset()
	require.NoError(t, l.writef(LevelInfo|LevelDebug, "Hello %q", "George"))
	require.Equal(t, "[DEBUG|INFO] Hello \"George\"\n", buff.String())

	buff.Reset()
	require.NoError(t, l.writef(LevelDebug, "Hello %q", "George"))
	require.Equal(t, "[DEBUG] Hello \"George\"\n", buff.String())

	l.SetPrefix("K8s")

	buff.Reset()
	require.NoError(t, l.writef(LevelDebug, "Hello %q", "George"))
	require.Equal(t, "[DEBUG][K8s] Hello \"George\"\n", buff.String())

	buff.Reset()
	require.NoError(t, l.writef(LevelError, "Hello %q", "George"))
	require.Equal(t, "", buff.String())

	const COMETS = 0x10000000
	l.RegisterLevel(COMETS, "COMETS")

	buff.Reset()
	require.NoError(t, l.writef(COMETS, "Hello %q", "George"))
	require.Equal(t, "", buff.String())

	l.SetLevel(COMETS)
	buff.Reset()
	require.NoError(t, l.writef(COMETS, "Hello %q", "George"))
	require.Equal(t, "[COMETS][K8s] Hello \"George\"\n", buff.String())

	l.SetLevel(COMETS | LevelDebug)
	buff.Reset()
	require.NoError(t, l.writef(COMETS, "Hello %q", "George"))
	require.Equal(t, "[COMETS][K8s] Hello \"George\"\n", buff.String())

	buff.Reset()
	require.NoError(t, l.writef(COMETS|LevelDebug, "Hello %q", "George"))
	require.Equal(t, "[COMETS|DEBUG][K8s] Hello \"George\"\n", buff.String())

	l.SetLevel(COMETS)
	buff.Reset()
	require.NoError(t, l.writef(COMETS|LevelDebug, "Hello %q", "George"))
	require.Equal(t, "[COMETS][K8s] Hello \"George\"\n", buff.String())

	l.UnRegisterLevel(COMETS)
	buff.Reset()
	require.NoError(t, l.writef(COMETS|LevelDebug, "Hello %q", "George"))
	require.Equal(t, "", buff.String())

	l.SetLevel(COMETS | LevelDebug)
	buff.Reset()
	require.NoError(t, l.writef(COMETS|LevelDebug, "Hello %q", "George"))
	require.Equal(t, "[DEBUG][K8s] Hello \"George\"\n", buff.String())

	l.UnRegisterLevel(LevelDebug)
	buff.Reset()
	require.NoError(t, l.writef(COMETS|LevelDebug, "Hello %q", "George"))
	require.Equal(t, "", buff.String())

	l.RegisterLevel(COMETS, "COMETS")
	buff.Reset()
	require.NoError(t, l.writef(COMETS|LevelDebug, "Hello %q", "George"))
	require.Equal(t, "[COMETS][K8s] Hello \"George\"\n", buff.String())
}
