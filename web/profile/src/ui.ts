import { User } from './types';
const PRESET_EMOJIS = ['🐉','🦊','🌟','🐼','🦁','🐸','🐺','🦄','🐯','🐨','👾','🤖','🧑','👩','🧔','👑'];
export function escapeHtml(s: string): string {
  return String(s || '').replace(/[&<>"']/g, c => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c] as string));
}
export function render(user: User, appEl: HTMLElement, onBioClick: () => void, onEmojiClick: (e: string) => void) {
  const displayName = user.display_name || user.game_username || ('User' + user.id);
  const color = user.profile_color || '#8b5cf6';
  const avatarContent = user.avatar ? `<img src="${escapeHtml(user.avatar)}">` : (user.name_emoji || '🎮');
  appEl.innerHTML = `
    <div class="banner"></div>
    <div class="avatar-wrap"><div class="avatar">${avatarContent}</div></div>
    <div class="name" style="color:${color}">${escapeHtml(displayName)}</div>
    <div class="username">@${escapeHtml(user.game_username || 'unknown')}</div>
    <div class="bio">${escapeHtml(user.bio || '---')}</div>
    <div class="stats">
      <div><div class="stat-value">${user.coins}</div><div class="stat-label">سکه</div></div>
      <div><div class="stat-value">${user.gems}</div><div class="stat-label">الماس</div></div>
      <div><div class="stat-value">${user.level}</div><div class="stat-label">سطح</div></div>
    </div>
    <button class="btn" id="btn-edit-bio">تغییر بیو</button>
    <div class="section">
      <div class="section-title">ایموجی نام</div>
      <div class="emoji-row">
        ${PRESET_EMOJIS.map(e => `<button class="emoji-btn ${user.name_emoji === e ? 'active' : ''}" data-emoji="${e}">${e}</button>`).join('')}
      </div>
    </div>
  `;
  document.getElementById('btn-edit-bio')?.addEventListener('click', onBioClick);
  document.querySelectorAll('.emoji-btn').forEach(btn => {
    btn.addEventListener('click', () => {
      const e = (btn as HTMLElement).dataset.emoji;
      if (e) onEmojiClick(e);
    });
  });
}
