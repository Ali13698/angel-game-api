package hub

func checkWin(b [9]int, p int) bool {
	lines := [8][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // ردیف‌ها
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // ستون‌ها
		{0, 4, 8}, {2, 4, 6},            // اریب‌ها
	}
	for _, l := range lines {
		if b[l[0]] == p && b[l[1]] == p && b[l[2]] == p { return true }
	}
	return false
}

func isFull(b [9]int) bool {
	for _, v := range b { if v == 0 { return false } }
	return true
}

func evaluate(b [9]int, botP, humanP int) int {
	if checkWin(b, botP) { return 10 }
	if checkWin(b, humanP) { return -10 }
	return 0
}

// الگوریتم Minimax برای انتخاب بهترین حرکت
func minimax(b [9]int, depth int, isMax bool, botP, humanP int) int {
	score := evaluate(b, botP, humanP)
	if score == 10 { return score - depth }
	if score == -10 { return score + depth }
	if isFull(b) { return 0 }

	if isMax {
		best := -1000
		for i := 0; i < 9; i++ {
			if b[i] == 0 {
				b[i] = botP
				val := minimax(b, depth+1, false, botP, humanP)
				b[i] = 0
				if val > best { best = val }
			}
		}
		return best
	} else {
		best := 1000
		for i := 0; i < 9; i++ {
			if b[i] == 0 {
				b[i] = humanP
				val := minimax(b, depth+1, true, botP, humanP)
				b[i] = 0
				if val < best { best = val }
			}
		}
		return best
	}
}

func getBestMove(b [9]int, botP, humanP int) int {
	bestVal := -1000
	bestMove := -1
	for i := 0; i < 9; i++ {
		if b[i] == 0 {
			b[i] = botP
			moveVal := minimax(b, 0, false, botP, humanP)
			b[i] = 0
			if moveVal > bestVal {
				bestMove = i
				bestVal = moveVal
			}
		}
	}
	return bestMove
}
