package utils

var DefaultStringUtil = &StringUtil{}

type StringUtil struct {
}

// GetEmojiForBool 获取布尔值的emoji
func (s *StringUtil) GetEmojiForBool(value bool) string {
	if value {
		return "✅"
	}
	return "❌"
}
