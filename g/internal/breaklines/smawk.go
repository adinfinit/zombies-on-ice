package breaklines

import (
	"math"
	"strings"
)

var infinity = float32(math.Inf(1))

// Implements SMAWK divide & conquer
// see http://xxyxyz.org/line-breaking/
func SMAWK(text string, frameWidth float32, measure func(text string) float32) []string {
	words := strings.Split(text, " ")

	offsets := make([]float32, 1, len(words)+1)
	next := float32(0.0)
	for _, word := range words {
		next += measure(word)
		offsets = append(offsets, next)
	}

	minima := make([]float32, 1, len(words)+1)
	for range words {
		minima = append(minima, infinity)
	}

	breaks := make([]int, len(words)+1)

	cost := func(i, k int) float32 {
		w := offsets[k] - offsets[i] + float32(k) - float32(i) - 1
		if w > frameWidth {
			return infinity
		}

		return minima[i] + (frameWidth-w)*(frameWidth-w)
	}

	search := func(i0, k0, i1, k1 int) {
		stack := [][4]int{{i0, k0, i1, k1}}
		for len(stack) > 0 {
			t := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			i0, k0, i1, k1 := t[0], t[1], t[2], t[3]
			if k0 < k1 {
				k := (k0 + k1) / 2
				for i := i0; i < i1; i++ { // todo check range
					c := cost(i, k)
					if c <= minima[k] {
						minima[k] = c
						breaks[k] = i
					}
				}
				stack = append(stack, [4]int{breaks[k], k + 1, i1, k1})
				stack = append(stack, [4]int{i0, k0, breaks[k] + 1, k})
			}
		}
	}

	n, i, offset := len(words)+1, 0, 0
	for {
		r := n
		if v := 1 << uint(i+1); v < r {
			r = v
		}

		edge := (1 << uint(i)) + offset
		search(offset, edge, edge, r+offset)
		x := minima[r-1+offset]

		for k := 1 << uint(i); k < r-1; k++ { // todo check range
			y := cost(k+offset, r-1+offset)
			if y <= x {
				n -= k
				i = 0
				offset += k
				break
			}
		}
		if r == n {
			break
		}
		i++
	}

	lines := []string{}
	k := len(words)
	for k > 0 {
		i = breaks[k]
		line := strings.Join(words[i:k], " ")
		lines = append(lines, line)
		k = i
	}

	for i, k := 0, len(lines)-1; i < k; i, k = i+1, k-1 {
		lines[i], lines[k] = lines[k], lines[i]
	}

	return lines
}
