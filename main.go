	mux := http.NewServeMux()
	profile.Register(mux)
	profile.RegisterAuth(mux)
	profile.RegisterDevice(mux) // ⬅️ این خط اضافه شود
