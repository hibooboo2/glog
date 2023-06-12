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
}
