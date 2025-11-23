package entities

import "math/rand"

func makeReviewerRandGen(col []User, except []UserID) func() (int, error) {
	buff := make(map[int]struct{})
	return func() (int, error) {
		for len(buff) != len(col) {
			idx := rand.Intn(len(col))

			if _, ok := buff[idx]; !ok {
				flagExcept := false
				for _, elem := range except {
					if elem == col[idx].ID {
						flagExcept = true
						break
					}
				}
				buff[idx] = struct{}{}

				if col[idx].IsActive && !flagExcept {
					return idx, nil
				}
			}
		}
		return -1, ErrReviewerAssign
	}
}
