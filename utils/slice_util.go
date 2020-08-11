package utils

import (
	"sort"
)

func RemoveDuplicatesByMap(slc []string) []string {
	var result []string
	tempMap := map[string]byte{}
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l {
			result = append(result, e)
		}
	}
	return result
}

func RemoveDuplicatesByLoop(slc []string) []string {
	var result []string
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false
				break
			}
		}
		if flag {
			result = append(result, slc[i])
		}
	}
	return result
}

func RemoveDuplicatesWithSort(slc []string) []string {
	if len(slc) == 0 {
		return nil
	}

	sort.Strings(slc)
	i, j := 0, 0

	for {
		if i >= len(slc)-1 {
			break
		}
		for j = i + 1; j < len(slc) && slc[i] == slc[j]; j++ {
		}
		slc = append(slc[:i+1], slc[j:]...)
		i++
	}

	return slc
}

func RemoveDuplicates(slc []string) []string {
	if len(slc) < 1024 {
		return RemoveDuplicatesByLoop(slc)
	} else {
		return RemoveDuplicatesByMap(slc)
	}
}
