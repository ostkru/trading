package utils

func PtrString(v string) *string { return &v }
func PtrInt64(v int64) *int64 { return &v }
func PtrFloat64(v float64) *float64 { return &v }
func PtrBool(v bool) *bool { return &v } 