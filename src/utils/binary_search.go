package utils

import (
	i "neon-chat/src/interfaces"
)

// arr is expected to be orderred by item id
func BinarySearch(
	arr []i.Identifiable,
	item_id uint,
) (found i.Identifiable, idx int) {
	// //
	// count := 0
	// start := time.Now()
	// //
	found = nil
	idx = -1
	len := len(arr)
	a := 0
	z := len - 1
	for a <= z && z >= 0 && a <= len {
		// //
		// count++
		// //
		mid_idx := (a + z) / 2
		mid_val := arr[mid_idx]
		mid_val_id := mid_val.GetId()

		if mid_val_id == item_id {
			found = arr[mid_idx]
			idx = mid_idx
			break
		}

		if mid_val_id < item_id {
			a = mid_idx + 1
		} else if item_id < mid_val_id {
			z = mid_idx - 1
		}
	}
	// //
	// elapsed := time.Since(start)
	// log.Printf("BinarySearch for id[%d] isfound:%t check took %s, %d loops in %d items\n",
	// 	item_id, idx >= 0, elapsed, count, len)
	// //

	return found, idx
}
