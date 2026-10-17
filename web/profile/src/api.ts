import { logInfo, logError } from './logger';

export async function fetchToken(userId: number): Promise<string | null> {
  try { 
    logInfo('Fetching token for user', userId);
    const r = await fetch('/api/auth/token?id=' + userId); 
    const d = await r.json(); 
    return d.ok ? d.token : null; 
  } catch(e) { logError('Fetch token failed', e); return null; }
}

export async function updateProfile(patch: any): Promise<{ok: boolean, error?: string}> {
  const p = new URLSearchParams(location.search).get('id'); 
  if (!p) return {ok: false, error: 'No ID'};
  const token = await fetchToken(parseInt(p, 10)); 
  if (!token) return {ok: false, error: 'Auth error'};
  
  logInfo('Updating profile with patch:', patch);
  try { 
    const r = await fetch('/api/me', { method: 'PATCH', headers: { 'Authorization': 'Bearer ' + token, 'Content-Type': 'application/json' }, body: JSON.stringify(patch) }); 
    const d = await r.json();
    if (!r.ok) logError('Update profile failed on server', d);
    return { ok: r.ok, error: d.error };
  } catch(e) { logError('Update profile network error', e); return {ok: false, error: 'Network error'}; }
}

export async function uploadImage(file: File): Promise<{ok: boolean, url?: string, error?: string}> {
  const p = new URLSearchParams(location.search).get('id'); 
  if (!p) return {ok: false, error: 'No ID'};
  const token = await fetchToken(parseInt(p, 10)); 
  if (!token) return {ok: false, error: 'Auth error'};
  
  logInfo('Uploading image:', file.name, 'Size:', file.size);
  const formData = new FormData(); 
  formData.append('image', file);
  
  try {
    const r = await fetch('/api/me/upload', { method: 'POST', headers: { 'Authorization': 'Bearer ' + token }, body: formData });
    const d = await r.json(); 
    if (!r.ok) logError('Upload failed on server', d);
    else logInfo('Upload successful', d.url);
    return { ok: r.ok, url: d.url, error: d.error };
  } catch(e) { logError('Upload network error', e); return {ok: false, error: 'Network error'}; }
}