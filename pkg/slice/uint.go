package slice

func DoseExist(list []uint, value uint) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false

}

func MapFromUint64ToUint(list []uint64) []uint {
	result := make([]uint, len(list))
	for i, item := range list {
		result[i] = uint(item)
	}
	return result
}

func MapFromUintToUint64(list []uint) []uint64 {
	result := make([]uint64, len(list))
	for i, item := range list {
		result[i] = uint64(item)
	}
	return result
}
