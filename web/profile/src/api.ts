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
