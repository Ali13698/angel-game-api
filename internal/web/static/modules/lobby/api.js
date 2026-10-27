export function getUserId() {
    const qs = new URLSearchParams(location.search).get("id");
    if (qs) return qs;
    const tg = window.Telegram?.WebApp?.initDataUnsafe?.user?.id;
    return tg ? String(tg) : "1001";
}

export async function fetchProfileData() {
    try {
        const userId = getUserId();
        const res = await fetch(`/api/profile/${userId}`);
        const data = await res.json();
        
        if (!res.ok || !data.ok) {
            throw new Error(data.error || "User not found");
        }
        return data.user;
    } catch (e) {
        console.error("API Error in fetchProfileData:", e);
        return null;
    }
}
