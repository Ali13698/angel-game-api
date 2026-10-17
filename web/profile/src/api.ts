export async function fetchToken(userId: number): Promise<string | null> {
  try { const r = await fetch('/api/auth/token?id=' + userId); const d = await r.json(); return d.ok ? d.token : null; } catch { return null; }
}
export async function updateProfile(patch: any): Promise<boolean> {
  const p = new URLSearchParams(location.search).get('id'); if (!p) return false;
  const token = await fetchToken(parseInt(p, 10)); if (!token) return false;
  try { const r = await fetch('/api/me', { method: 'PATCH', headers: { 'Authorization': 'Bearer ' + token, 'Content-Type': 'application/json' }, body: JSON.stringify(patch) }); return r.ok; } catch { return false; }
}