package utils

import "prplchat/src/utils/interfaces"

// arr is expected to be orderred by item id
func BinarySearch(
	arr []interfaces.Identifiable,
	item_id uint,
) interfaces.Identifiable {
	a := 0
	z := len(arr) - 1
	for a <= z && z >= 0 && a <= len(arr) {
		mid_idx := (a + z) / 2
		mid_val := arr[mid_idx]
		mid_val_id := mid_val.GetId()

		if mid_val_id == item_id {
			return arr[mid_idx]
		}

		if mid_val_id < item_id {
			a = mid_idx + 1
		} else if item_id < mid_val_id {
			z = mid_idx - 1
		}
	}
	return nil
}
