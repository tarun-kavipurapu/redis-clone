package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadSimpleString(t *testing.T) {
	cases := map[string]string{
		"+OK\r\n": "OK",
	}

	for k, v := range cases {
		v_out, _ := Decode([]byte(k))
		assert.Equal(t, v, v_out)
	}
}
func TestReadError(t *testing.T) {
	cases := map[string]string{
		"-ERR unknown command 'foobar'\r\n": "ERR unknown command 'foobar'",
	}

	for k, v := range cases {
		v_out, _ := Decode([]byte(k))
		assert.Equal(t, v, v_out)

	}
}

func TestReadBulkString(t *testing.T) {
	cases := map[string]string{
		"$6\r\nfoobar\r\n":       "foobar",
		"$0\r\n\r\n":             "",
		"$11\r\nhello world\r\n": "hello world",
	}

	for input, expected := range cases {
		out, _ := Decode([]byte(input))
		assert.Equal(t, expected, out)
	}
}

func TestReadArray(t *testing.T) {
	cases := map[string]interface{}{
		"*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n": []interface{}{"foo", "bar"},
		"*3\r\n:1\r\n:2\r\n:3\r\n":         []interface{}{int64(1), int64(2), int64(3)},
		"*0\r\n":                           []interface{}{},
		// "*-1\r\n":                          nil,
		"*2\r\n*2\r\n:1\r\n:2\r\n*2\r\n+OK\r\n$5\r\nhello\r\n": []interface{}{
			[]interface{}{int64(1), int64(2)},
			[]interface{}{"OK", "hello"},
		},
	}

	for input, expected := range cases {
		out, _ := Decode([]byte(input))
		assert.Equal(t, expected, out)
	}
}
