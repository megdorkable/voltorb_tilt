package main

import (
	"fmt"
	"reflect"
)

func (b *Board) update() bool {
	// save previous poss to compare later
	var prev_poss [][][]int
	for row := range b.Grid {
		row_poss := [][]int{}
		for _, val := range b.Grid[row] {
			row_poss = append(row_poss, val.Poss)
		}
		prev_poss = append(prev_poss, row_poss)
	}

	// refresh the sum, voltorb count, and unknown count stats
	b.refresh_stats()
	// update the rows
	for row := range b.Grid {
		b.update_poss(row)
	}
	// update the columns
	b.rotate(90)
	for row := range b.Grid {
		b.update_poss(row)
	}
	b.rotate(-90)

	// set the known values
	b.set_known()

	// check if anything has changed
	changed := false
	solved := true
	for row := range b.Grid {
		for col, val := range b.Grid[row] {
			if !reflect.DeepEqual(val.Poss, prev_poss[row][col]) {
				changed = true
			}
			if val.Value == UNKNOWN && (has(val.Poss, VAL_TWO) || has(val.Poss, VAL_THREE)) {
				solved = false
			}
		}
	}

	if changed {
		// fmt.Println("running again")
		solved = b.update()
	} else if !solved {
		flipped := false
		for row := range b.Grid {
			for col, val := range b.Grid[row] {
				if val.Value == UNKNOWN && !has(val.Poss, VOLTORB) && (has(val.Poss, VAL_TWO) || has(val.Poss, VAL_THREE)) {
					fmt.Println(b)
					f := flip(row, col)
					b.Grid[row][col].Value = f
					b.Grid[row][col].Poss = []int{f}
					if f == VOLTORB {
						return false
					}
					flipped = true
					solved = b.update()
					break
				}
			}
		}

		if !flipped {
			for row := range b.Grid {
				for col, val := range b.Grid[row] {
					if val.Value == UNKNOWN && (has(val.Poss, VAL_TWO) || has(val.Poss, VAL_THREE)) {
						fmt.Println(b)
						f := flip(row, col)
						b.Grid[row][col].Value = f
						b.Grid[row][col].Poss = []int{f}
						if f == VOLTORB {
							return false
						}
						solved = b.update()
						break
					}
				}
			}
		}
	}

	return solved
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

	// skip complete rows
	if row_rem_unknown == 0 {
		return updated
	}

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
					if pval != VOLTORB {
						new_poss = append(new_poss, pval)
					}
				}
				(*b).Grid[row][col].Poss = new_poss
				updated = true
			}
		}
	}

	// if remaining number unknown - 1 == remaining number voltorbs, must be either VOLTORB or remaining sum
	if row_rem_unknown-1 == row_rem_vcount {
		for col := range b.Grid[row] {
			if b.Grid[row][col].Value == UNKNOWN {
				new_poss := []int{}
				for _, pval := range b.Grid[row][col].Poss {
					if pval == VOLTORB || pval == row_rem_sum {
						new_poss = append(new_poss, pval)
					}
				}
				(*b).Grid[row][col].Poss = new_poss
				updated = true
			}
		}
	}

	// if remaining number unknown - remaining number voltorbs <= (remaining sum + 1)/3, can't be VAL_ONE
	if row_rem_unknown-row_rem_vcount <= (row_rem_sum+1)/3 {
		for col := range b.Grid[row] {
			if b.Grid[row][col].Value == UNKNOWN {
				new_poss := []int{}
				for _, pval := range b.Grid[row][col].Poss {
					if pval != VAL_ONE {
						new_poss = append(new_poss, pval)
					}
				}
				(*b).Grid[row][col].Poss = new_poss
				updated = true
			}
		}
	}

	// if remaining number unknown - remaining number voltorbs == remaining sum - 1, can't be VAL_THREE
	if row_rem_unknown-row_rem_vcount == row_rem_sum-1 {
		for col := range b.Grid[row] {
			if b.Grid[row][col].Value == UNKNOWN {
				new_poss := []int{}
				for _, pval := range b.Grid[row][col].Poss {
					if pval != VAL_THREE {
						new_poss = append(new_poss, pval)
					}
				}
				(*b).Grid[row][col].Poss = new_poss
				updated = true
			}
		}
	}

	// if remaining number voltorbs == remaining number unknown, must be VOLTORB
	if row_rem_vcount == row_rem_unknown {
		for col := range b.Grid[row] {
			if b.Grid[row][col].Value == UNKNOWN {
				new_poss := []int{VOLTORB}
				(*b).Grid[row][col].Poss = new_poss
				updated = true
			}
		}
	}

	// if remaining sum == remaining number unknown * 3, must be VAL_THREE
	if row_rem_sum == row_rem_unknown*3 {
		for col := range b.Grid[row] {
			if b.Grid[row][col].Value == UNKNOWN {
				new_poss := []int{VAL_THREE}
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
	rs := b.RowStats
	cs := b.ColStats
	if direction == 90 {
		b.Vertical = h
		b.Horizontal = reverse(v)
		b.RowStats = cs
		b.ColStats = reverse_stats(rs)
	} else if direction == -90 {
		b.Vertical = reverse(h)
		b.Horizontal = v
		b.RowStats = reverse_stats(cs)
		b.ColStats = rs
	}
}

func reverse(arr [num_cols][num_inputs]int) [num_cols][num_inputs]int {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func reverse_stats(arr [num_cols]Stats) [num_cols]Stats {
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

func has(arr []int, i int) bool {
	for _, val := range arr {
		if i == val {
			return true
		}
	}
	return false
}

func flip(row int, col int) int {
	var i int

	fmt.Printf("Flip the tile at [%d][%d], and type the result: ", row, col)
	fmt.Scan(&i)

	return i
}
