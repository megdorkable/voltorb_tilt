package main

import "fmt"

func (b *Board) update() bool {
	updated := false

	b.refresh_stats()
	for row := range b.Grid {
		updated = b.update_poss(row) || updated
	}
	b.rotate(90)
	for row := range b.Grid {
		updated = b.update_poss(row) || updated
	}
	b.rotate(-90)

	updated = b.set_known() || updated

	return updated
}

func (b *Board) refresh_stats() {
	b.RowStats = [num_rows]Stats{}
	b.ColStats = [num_cols]Stats{}

	for row := range b.Grid {
		for col, val := range b.Grid[row] {
			if val.Value == -1 {
				b.RowStats[row].Vcount += 1
				b.ColStats[col].Vcount += 1
			} else if val.Value == 0 {
				b.RowStats[row].Unkcount += 1
				b.ColStats[col].Unkcount += 1
			} else {
				b.RowStats[row].Sum += val.Value
				b.ColStats[col].Sum += val.Value
			}
		}
	}
}

func (b *Board) update_poss(row int) bool {
	updated := false

	row_total_sum := b.Vertical[row][0]
	row_total_vcount := b.Vertical[row][1]
	row_rem_sum := row_total_sum - b.RowStats[row].Sum
	row_rem_vcount := row_total_vcount - b.RowStats[row].Vcount
	row_rem_unknown := b.RowStats[row].Unkcount

	// if remaining sum + remaining number voltorb == remaining number unknown, must be VOLTORB OR VAL_ONE
	if row_rem_sum+row_rem_vcount == row_rem_unknown {
		for col := range b.Grid[row] {
			if b.Grid[row][col].Value == UNKNOWN {
				new_poss := []int{}
				for _, pval := range b.Grid[row][col].Poss {
					if pval == VOLTORB || pval == VAL_ONE {
						new_poss = append(new_poss, pval)
					}
				}
				(*b).Grid[row][col].Poss = new_poss
				updated = true
			}
		}
	}

	// if no remaining voltorbs, can't be VOLTORB
	if row_rem_vcount == 0 {
		for col := range b.Grid[row] {
			if b.Grid[row][col].Value == UNKNOWN {
				new_poss := []int{}
				for _, pval := range b.Grid[row][col].Poss {
					if pval != -1 {
						new_poss = append(new_poss, pval)
					}
				}
				(*b).Grid[row][col].Poss = new_poss
				updated = true
			}
		}
	}

	return updated
}

func (b *Board) set_known() bool {
	updated := false

	for row := range b.Grid {
		for col := range b.Grid[row] {
			if b.Grid[row][col].Value == 0 && len(b.Grid[row][col].Poss) == 1 {
				b.Grid[row][col].Value = b.Grid[row][col].Poss[0]
				updated = true
			}
		}
	}

	return updated
}

func (b *Board) rotate(direction int) {
	grid_len := len(b.Grid)
	num_layers_to_rotate := grid_len / 2

	for layer := 0; layer < num_layers_to_rotate; layer++ {
		first := layer
		last := grid_len - first - 1

		for i := first; i < last; i++ {
			offset := i - first

			top_left := b.Grid[first][i]
			top_right := b.Grid[i][last]
			bottom_right := b.Grid[last][last-offset]
			bottom_left := b.Grid[last-offset][first]

			if direction == 90 {
				b.Grid[first][i] = bottom_left
				b.Grid[i][last] = top_left
				b.Grid[last][last-offset] = top_right
				b.Grid[last-offset][first] = bottom_right

			} else if direction == -90 {
				b.Grid[first][i] = top_right
				b.Grid[i][last] = bottom_right
				b.Grid[last][last-offset] = bottom_left
				b.Grid[last-offset][first] = top_left
			}
		}
	}

	v := b.Vertical
	h := b.Horizontal
	if direction == 90 {
		b.Vertical = h
		b.Horizontal = reverse(v)
	} else if direction == -90 {
		b.Vertical = reverse(h)
		b.Horizontal = v
	}

	b.Grid = b.Grid
}

func reverse(arr [num_cols][num_inputs]int) [num_cols][num_inputs]int {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func remove(s *[]int, i int) {
	fmt.Println(s)
	(*s)[i] = (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
}
