export async function fetchProfile(id: number): Promise<User | null> {
    try {
        const res = await fetch(`/api/profile/${id}`);
        if (!res.ok) return null;
        const data = await res.json();
        return data.user;
    } catch (e) {
        return null;
    }
}

export async function updateProfile(id: number, patch: any): Promise<boolean> {
    try {
        const tokenRes = await fetch(`/api/auth/token?id=${id}`);
        const tokenData = await tokenRes.json();
        if (!tokenData.ok) return false;

        const res = await fetch('/api/me', {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${tokenData.token}`
            },
            body: JSON.stringify(patch)
        });
        return res.ok;
    } catch (e) {
        return false;
    }
}
