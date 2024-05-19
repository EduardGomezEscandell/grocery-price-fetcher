package formatter

type euro struct {
	float32
}

//nolint:revive // We keep the internal types opaque
func Euro(value float32) euro {
	return euro{value}
}
