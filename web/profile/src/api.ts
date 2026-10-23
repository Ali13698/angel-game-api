export async function fetchProfile(id: number): Promise<UserProfile | null> {
    try {
        const res = await fetch(`/api/profile/${id}`);
        if (!res.ok) return null;
        const data = await res.json();
        return data.user; // فرض بر این است که سرور در فیلد user برمی‌گرداند
    } catch {
        return null;
    }
}

export async function updateProfile(id: number, patch: any): Promise<boolean> {
    try {
        // گرفتن توکن احراز هویت (طبق ساختار قبلی پروژه)
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
    } catch {
        return false;
    }
}
