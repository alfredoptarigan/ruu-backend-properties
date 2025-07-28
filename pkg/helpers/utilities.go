package helpers

func CheckLengthUUID(uuid string) bool {
	if len(uuid) != 36 {
		return false
	}

	return true
}
