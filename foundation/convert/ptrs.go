package convert

func ToIntPtr(i int) *int {
	return &i
}

func ToStringPtr(s string) *string {
	return &s
}

func ToBoolPtr(a bool) *bool {
	return &a
}
