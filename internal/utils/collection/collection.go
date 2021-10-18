package collection

import (
	"errors"
	"math/rand"
	"reflect"
	"time"
)

func SliceShuffle(slice interface{}) error {
	v := reflect.ValueOf(slice)
	if v.Type().Kind() != reflect.Slice {
		return errors.New("unsupported type to do slice shuffle")
	}

	l := v.Len()
	if l < 2 {
		return nil
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	swap := reflect.Swapper(slice)
	for i := 0; i < l; i++ {
		j := r.Intn(i + 1)
		swap(i, j)
	}
	return nil
}

func GenerateUnRepeatableRandomNumbers(start int, end int, count int) ([]int, error) {
	if start == 0 && end == 0 {
		return nil, errors.New("can not generate numbers by 0 range")
	}

	if end < start || (end-start) < count {
		return nil, errors.New("incorrect random range")
	}

	nums := make([]int, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(nums) < count {
		num := r.Intn(end-start) + start
		isExist := false
		for _, n := range nums {
			if n == num {
				isExist = true
				break
			}
		}

		if !isExist {
			nums = append(nums, num)
		}
	}
	return nums, nil
}
