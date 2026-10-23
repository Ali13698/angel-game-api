function getUserId() {
    let t = new URLSearchParams(location.search).get("id");
    if (t && /^\d+$/.test(t)) return parseInt(t, 10);
    let e = window.Telegram?.WebApp?.initDataUnsafe?.user?.id;
    return e ? Number(e) : 1001;
}

async function joinGame(gameId) {
    let userId = getUserId();
    if (!userId) return alert("ابتدا وارد شوید");
    
    try {
        // ۱. درخواست ساخت اتاق به سرور Go
        let res = await fetch("/rooms", {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            // برای تست، یک بازیکن مجازی (bot) هم اضافه می‌کنیم
            body: JSON.stringify({game_id: gameId, players: [userId.toString(), "bot"]})
        });
        
        if (!res.ok) throw new Error(await res.text());
        let data = await res.json();
        
        // ۲. هدایت کاربر به صفحه اختصاصی بازی با آیدی اتاق
        window.location.href = `/ludo.html?room=${data.room_id}&id=${userId}`;

    } catch (err) {
        alert(`خطا در ساخت اتاق: ${err.message}`);
    }
}

window.joinGame = joinGame;
