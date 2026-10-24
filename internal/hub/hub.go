// ... (همان کدهای قبلی تا خط ۴۰)
	if len(players) == 1 {
		// ✅ پاس دادن gameID به ربات
		go h.runBot(rm, types.UserID("bot_1"), gameID, 3*time.Second)
	}
// ... (ادامه کد)
