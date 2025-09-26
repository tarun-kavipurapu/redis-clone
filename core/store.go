package core

//key-->string
//val -->
// can we make this into a singleton instance

var store map[string]*Value

type Value struct {
	value  interface{}
	expiry int64
}

func init() {
	store = make(map[string]*Value)
}

func NewValue(value interface{}, expiry int64) *Value {
	val := &Value{
		value:  value,
		expiry: expiry,
	}
	return val
}
func Put(key string, val *Value) {
	store[key] = val
}

func Get(key string) *Value {
	val := store[key]

	return val
}
