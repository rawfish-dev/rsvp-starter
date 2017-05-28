package utils

func StringPointer(value string) *string {
	return &value
}

func StringPointerValue(stringPointer *string) string {
	if stringPointer == nil {
		return ""
	}

	return *stringPointer
}
