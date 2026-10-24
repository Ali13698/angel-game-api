function getUserId() {
    let t = new URLSearchParams(location.search).get("id");
    if (t && /^\d+$/.test(t)) return parseInt(t, 10);
    let e = window.Telegram?.WebApp?.initDataUnsafe?.user?.id;
    return e ? Number(e) : 1001;
}

async function joinGame(gameId) {
    let userId = getUserId();
    try {
        let res = await fetch("/rooms", {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({game_id: gameId, players: [userId.toString(), "bot"]})
        });
        if (!res.ok) throw new Error(await res.text());
        let data = await res.json();
        
        // ✅ لینک دقیقاً با پارامترهای AL V4
        window.location.href = `/ludo.html?user=${userId}&room=${data.room_id}`;

    } catch (err) {
        alert(`خطا: ${err.message}`);
    }
}
window.joinGame = joinGame;
