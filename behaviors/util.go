package behaviors

func pairHash(id1, id2 int) int {
	switch {
	case id1 > id2:
		return id1<<16 | id2&0xFFFF
	case id1 < id2:
		return id2<<16 | id1&0xFFFF
	default:
		return -1
	}
}
