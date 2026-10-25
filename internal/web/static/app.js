(()=>{
  async function fetchProfile(id) {
    try {
      let t = await fetch(`/api/profile/${id}`);
      return t.ok ? (await t.json()).user : null;
    } catch { return null; }
  }
  
  var heartSvg = (filled) => `<svg width="28" height="28" viewBox="0 0 24 24" fill="${filled?"#ef4444":"none"}" stroke="${filled?"#ef4444":"#ffffff"}" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path></svg>`;

  function renderProfile(e, t) {
    return `
    <div class="profile-card">
      <div class="banner" style="background: ${e.banner_id||"linear-gradient(135deg, #7c3aed 0%, #ec4899 100%)"}; background-size: cover;">
        <button class="plus-badge banner-plus" id="btn-edit-banner">+</button>
      </div>
      <div class="avatar-container">
        <div class="avatar"><img id="profile-avatar-img" src="${e.avatar||"https://api.dicebear.com/7.x/bottts/svg?seed=1"}" /></div>
        <button class="plus-badge avatar-plus" id="btn-change-avatar">+</button>
      </div>
      <div class="profile-info">
        <div class="name" style="color:${e.profile_color||"#c084fc"}">${e.display_name||e.game_username}</div>
        <div class="username">@${e.game_username||"unknown"}</div>
        <div class="bio-wrapper">
            <span class="bio" id="profile-bio">${e.bio||"بیوگرافی..."}</span>
            <button class="pencil-btn" id="btn-edit-bio">✏️</button>
        </div>
      </div>
      <div class="main-stats">
        <div class="stat-box" id="like-box" style="cursor:pointer;">
          <div class="stat-value" id="heart-icon">${heartSvg(t)}</div>
          <div class="stat-label" id="likes-count">${e.likes_count||0}</div>
        </div>
        <div class="stat-box"><div class="stat-value">${e.age_days||1}</div><div class="stat-label">روز</div></div>
        <div class="stat-box"><div class="stat-value">${e.wins||0}</div><div class="stat-label">برد</div></div>
      </div>
      <div class="action-buttons">
        <button class="btn-square-small green">💎 ${e.gems||0}</button>
        <button class="btn-square-small">🪙 ${e.coins||0}</button>
      </div>
    </div>`;
  }

  function renderGames() {
    return `
    <div style="padding: 20px;">
        <div style="text-align:center; background:#ef4444; color:white; padding:10px; border-radius:8px; font-weight:bold; margin-bottom:20px;">
          پنل بازی‌ها (نسخه 4.1)
        </div>
        <h2 style="text-align:center; margin-bottom: 20px;">بازی‌های آنلاین</h2>
        <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 15px;">
            
            <div class="game-card" onclick="window.joinGame('tic-tac-toe', 'tictactoe.html')" style="background: var(--inner-card-bg); padding: 20px; border-radius: 12px; text-align: center; cursor: pointer; border: 2px solid #3b82f6;">
                <div style="font-size: 50px; margin-bottom: 10px;">⭕</div>
                <div style="font-size: 16px; font-weight: bold;">دوز (Tic-Tac-Toe)</div>
            </div>

            <div class="game-card" onclick="window.joinGame('ludo', 'ludo.html')" style="background: var(--inner-card-bg); padding: 20px; border-radius: 12px; text-align: center; cursor: pointer;">
                <div style="font-size: 50px; margin-bottom: 10px;">🎲</div>
                <div style="font-size: 16px; font-weight: bold;">منچ ساده</div>
            </div>

            <div class="game-card" onclick="window.joinGame('ludo-premium', 'ludo-pro.html')" style="background: var(--inner-card-bg); padding: 20px; border-radius: 12px; text-align: center; cursor: pointer; border: 2px solid #8b5cf6;">
                <div style="font-size: 50px; margin-bottom: 10px;">🎯</div>
                <div style="font-size: 16px; font-weight: bold;">منچ پرمیوم</div>
            </div>

        </div>
    </div>`;
  }

  var currentUser = null, isLiked = false;

  function getUserId() {
    let t = new URLSearchParams(location.search).get("id");
    if (t) return parseInt(t, 10);
    let i = window.Telegram?.WebApp;
    return Number(i?.initDataUnsafe?.user?.id) || 1001;
  }

  function renderLayout(content, activeTab) {
    return `
    <div class="app-shell">
        <div class="top-header" style="background: #0c0f14; padding: 10px; display: flex; justify-content: space-between; align-items: center;">
            <div style="font-weight: bold; font-size: 20px;">Angel</div>
            <img src="${currentUser?.avatar||"https://api.dicebear.com/7.x/bottts/svg?seed="+(currentUser?.id||1)}" style="width: 36px; height: 36px; border-radius: 50%; border: 2px solid var(--accent-blue);" />
        </div>
        <div class="page-content" style="padding-bottom: 80px;">
            ${content}
        </div>
        <div class="bottom-nav">
            <div class="nav-item ${activeTab==="profile"?"active":""}" onclick="switchTab('profile')">
                <span class="nav-icon">👤</span>
            </div>
            <div class="nav-center-btn" onclick="switchTab('games')">
                <span class="nav-icon">🎮</span>
            </div>
            <div class="nav-item" onclick="window.location.href='https://t.me'">
                <span class="nav-icon">💬</span>
            </div>
        </div>
    </div>`;
  }

  window.switchTab = (tab) => {
    let appEl = document.getElementById("app");
    if (!appEl || !currentUser) return;
    if (tab === "profile") {
        appEl.innerHTML = renderLayout(`<div class="container">${renderProfile(currentUser, isLiked)}</div>`, "profile");
        bindProfileEvents();
    } else if (tab === "games") {
        appEl.innerHTML = renderLayout(renderGames(), "games");
    }
  };

  function bindProfileEvents() {
    let getEl = (id) => document.getElementById(id);
    getEl("like-box")?.addEventListener("click", () => {
        isLiked = !isLiked;
        let icon = getEl("heart-icon")?.querySelector("svg");
        let countEl = getEl("likes-count");
        if (icon) {
            icon.setAttribute("fill", isLiked ? "#ef4444" : "none");
            icon.setAttribute("stroke", isLiked ? "#ef4444" : "#ffffff");
        }
        if (countEl && currentUser) {
            currentUser.likes_count += isLiked ? 1 : -1;
            countEl.innerText = currentUser.likes_count.toString();
        }
    });
  }

  async function init() {
    let appEl = document.getElementById("app");
    if (!appEl) return;
    appEl.innerHTML = '<div style="text-align:center; padding:50px; color:white;">در حال بارگذاری...</div>';
    
    let userId = getUserId();
    let user = await fetchProfile(userId);
    
    if (!user) {
        user = { id: userId, game_username: "user_" + userId, display_name: "User " + userId, avatar: "", bio: "", likes_count: 0, coins: 0, gems: 0, wins: 0, age_days: 1 };
    }
    currentUser = user;
    window.switchTab("profile");
  }

  window.joinGame = async function(gameId, htmlFile) {
    console.log("Attempting to join:", gameId, "Redirecting to:", htmlFile);
    let userId = getUserId();
    let statusEl = document.getElementById("app");
    statusEl.innerHTML = '<div style="text-align:center; padding:50px; color:white;">در حال ساخت اتاق...</div>';
    
    try {
        let res = await fetch("/rooms", {
            method: "POST", headers: {"Content-Type": "application/json"},
            body: JSON.stringify({game_id: gameId, players: [userId.toString(), "bot"]})
        });
        let data = await res.json();
        if(!res.ok) throw new Error(data.error || "خطای سرور");
        window.location.href = `/${htmlFile}?user=${userId}&room=${data.room_id}`;
    } catch(e) {
        alert("خطا در ساخت اتاق: " + e.message);
    }
  };

  document.addEventListener("DOMContentLoaded", init);
})();
