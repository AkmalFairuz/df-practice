package helper

import "github.com/df-mc/dragonfly/server/item"

func SetItemAsUnbreakable(stack item.Stack) item.Stack {
	return stack.WithValue("__unbreakable__", true)
}

func IsItemUnbreakable(stack item.Stack) bool {
	v, _ := stack.Value("__unbreakable__")
	return v != nil && v.(bool)
}
