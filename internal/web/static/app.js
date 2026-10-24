(()=>{async function d(e){try{let i=await fetch(`/api/profile/${e}`);return i.ok?(await i.json()).user:null}catch{return null}}var l=e=>`<svg width="28" height="28" viewBox="0 0 24 24" fill="${e?"#ef4444":"none"}" stroke="${e?"#ef4444":"#ffffff"}" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path></svg>`;function r(e,i){return`
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
            <span class="bio" id="profile-bio">${e.bio||"\u0628\u06CC\u0648\u06AF\u0631\u0627\u0641\u06CC..."}</span>
            <button class="pencil-btn" id="btn-edit-bio">\u270F\uFE0F</button>
        </div>
      </div>
      <div class="main-stats">
        <div class="stat-box" id="like-box" style="cursor:pointer;">
          <div class="stat-value" id="heart-icon">${l(i)}</div>
          <div class="stat-label" id="likes-count">${e.likes_count||0}</div>
        </div>
        <div class="stat-box"><div class="stat-value">${e.age_days||1}</div><div class="stat-label">\u0631\u0648\u0632</div></div>
        <div class="stat-box"><div class="stat-value">${e.wins||0}</div><div class="stat-label">\u0628\u0631\u062F</div></div>
      </div>
      <div class="action-buttons">
        <button class="btn-square-small green">\u{1F48E} ${e.gems||0}</button>
        <button class="btn-square-small">\u{1FA99} ${e.coins||0}</button>
      </div>
    </div>`}function o(){return`
    <div style="padding: 20px;">
        <h2 style="text-align:center; margin-bottom: 20px;">\u0628\u0627\u0632\u06CC\u200C\u0647\u0627\u06CC \u0622\u0646\u0644\u0627\u06CC\u0646</h2>
        <div style="display: grid; grid-template-columns: 1fr; gap: 15px;">
            <div class="game-card" onclick="window.joinGame('ludo')" style="background: var(--inner-card-bg); padding: 20px; border-radius: 12px; text-align: center; cursor: pointer;">
                <div style="font-size: 50px; margin-bottom: 10px;">\u{1F3B2}</div>
                <div style="font-size: 16px; font-weight: bold;">\u0645\u0646\u0686 (Ludo)</div>
            </div>
        </div>
    </div>`}var t=null,s=!1;function c(){let e=new URLSearchParams(location.search).get("id");if(e)return parseInt(e,10);let i=window.Telegram?.WebApp;return Number(i?.initDataUnsafe?.user?.id)||1001}function n(e,i){return`
    <div class="app-shell">
        <div class="top-header" style="background: #0c0f14; padding: 10px; display: flex; justify-content: space-between; align-items: center;">
            <div style="font-weight: bold; font-size: 20px;">Angel</div>
            <img src="${t?.avatar||"https://api.dicebear.com/7.x/bottts/svg?seed="+(t?.id||1)}" style="width: 36px; height: 36px; border-radius: 50%; border: 2px solid var(--accent-blue);" />
        </div>
        <div class="page-content" style="padding-bottom: 80px;">
            ${e}
        </div>
        <div class="bottom-nav">
            <div class="nav-item ${i==="profile"?"active":""}" onclick="switchTab('profile')">
                <span class="nav-icon">\u{1F464}</span>
            </div>
            <div class="nav-center-btn" onclick="switchTab('games')">
                <span class="nav-icon">\u{1F3AE}</span>
            </div>
            <div class="nav-item" onclick="window.location.href='https://t.me'">
                <span class="nav-icon">\u{1F4AC}</span>
            </div>
        </div>
    </div>`}window.switchTab=e=>{let i=document.getElementById("app");!i||!t||(e==="profile"?(i.innerHTML=n(`<div class="container">${r(t,s)}</div>`,"profile"),v()):e==="games"&&(i.innerHTML=n(o(),"games")))};function v(){let e=i=>document.getElementById(i);e("like-box")?.addEventListener("click",()=>{s=!s;let i=e("heart-icon")?.querySelector("svg"),a=e("likes-count");i&&(i.setAttribute("fill",s?"#ef4444":"none"),i.setAttribute("stroke",s?"#ef4444":"#ffffff")),a&&t&&(t.likes_count+=s?1:-1,a.innerText=t.likes_count.toString())})}async function u(){let e=document.getElementById("app");if(!e)return;e.innerHTML='<div style="text-align:center; padding:50px; color:white;">\u062F\u0631 \u062D\u0627\u0644 \u0628\u0627\u0631\u06AF\u0630\u0627\u0631\u06CC...</div>';let i=c(),a=await d(i);a||(a={id:i,game_username:"user_"+i,display_name:"User "+i,avatar:"",bio:"",likes_count:0,coins:0,gems:0,wins:0,age_days:1}),t=a,window.switchTab("profile")}document.addEventListener("DOMContentLoaded",u);})();
