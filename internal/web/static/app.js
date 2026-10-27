(()=>{var o=class{currentUser=null;currentStats=[];isLiked=!1;setUser(t,r){this.currentUser=t,this.currentStats=r||[]}setLiked(t){this.isLiked=t}getUser(){return this.currentUser}},i=new o;function s(){let t=new URLSearchParams(window.location.search).get("id");if(t)return parseInt(t,10);let r=window.Telegram?.WebApp;return r?.initDataUnsafe?.user?.id?r.initDataUnsafe.user.id:1001}async function d(e){try{let t=e||s(),r=await fetch(`/api/profile/${t}`);if(!r.ok)throw new Error("Network response was not ok");return await r.json()}catch(t){return console.error("Failed to fetch profile",t),null}}function c(e,t){if(!e)return'<div style="text-align:center;padding:50px;">\u062F\u0631 \u062D\u0627\u0644 \u0628\u0627\u0631\u06AF\u0630\u0627\u0631\u06CC...</div>';let r=a=>`<svg width="28" height="28" viewBox="0 0 24 24" fill="${a?"#ef4444":"none"}" stroke="${a?"#ef4444":"#ffffff"}" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path></svg>`;return`
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
          <div class="stat-value" id="heart-icon">${r(t)}</div>
          <div class="stat-label" id="likes-count">${e.likes_count||0}</div>
        </div>
        <div class="stat-box"><div class="stat-value">${e.age_days||1}</div><div class="stat-label">\u0631\u0648\u0632</div></div>
        <div class="stat-box"><div class="stat-value">${e.wins||0}</div><div class="stat-label">\u0628\u0631\u062F</div></div>
      </div>
      <div class="action-buttons">
        <button class="btn-square-small green">\u{1F48E} ${e.gems||0}</button>
        <button class="btn-square-small">\u{1FA99} ${e.coins||0}</button>
      </div>
    </div>`}function u(){return`
    <div style="padding: 20px;">
        <h2>\u062E\u0648\u0634 \u0622\u0645\u062F\u06CC\u062F</h2>
        <p>\u0628\u0647 \u067E\u0644\u062A\u0641\u0631\u0645 \u0628\u0627\u0632\u06CC \u0622\u0646\u062C\u0644 \u062E\u0648\u0634 \u0622\u0645\u062F\u06CC\u062F. \u0627\u0632 \u0645\u0646\u0648\u06CC \u067E\u0627\u06CC\u06CC\u0646 \u0628\u0647 \u0628\u062E\u0634\u200C\u0647\u0627\u06CC \u0645\u062E\u062A\u0644\u0641 \u062F\u0633\u062A\u0631\u0633\u06CC \u062F\u0627\u0634\u062A\u0647 \u0628\u0627\u0634\u06CC\u062F.</p>
    </div>`}function p(){return`
    <div style="padding: 20px;">
        <h2>\u0628\u0627\u0632\u06CC\u200C\u0647\u0627</h2>
        <div class="games-grid">
            ${[{id:"ludo",name:"\u0645\u0646\u0686 \u0641\u0636\u0627\u06CC\u06CC",icon:"\u{1F3B2}",route:"/games/ludo/lobby.html"}].map(t=>`
                <div class="game-card" onclick="window.location.href='${t.route}'">
                    <div class="game-icon">${t.icon}</div>
                    <div class="game-name">${t.name}</div>
                </div>
            `).join("")}
        </div>
    </div>`}var l={log:(...e)=>{console.log("[Angel]",...e)},error:(...e)=>{console.error("[Angel Error]",...e)},warn:(...e)=>{console.warn("[Angel Warn]",...e)}};window.Router={switch:async(e,t)=>{document.querySelectorAll(".screen").forEach(n=>n.classList.remove("active")),document.querySelectorAll(".nav-item").forEach(n=>n.classList.remove("active"));let r=document.getElementById(`screen-${e}`);r&&r.classList.add("active"),t&&t.classList.add("active");let a=document.getElementById("screen-content");a&&(e==="home"?a.innerHTML=u():e==="profile"?a.innerHTML=c(i.currentUser,i.isLiked):e==="games"&&(a.innerHTML=p()))}};async function f(){window.Telegram?.WebApp&&(window.Telegram.WebApp.ready(),window.Telegram.WebApp.expand());let e=s(),t=await d(e);t&&t.ok?(i.setUser(t.user,t.stats),l.log("Profile loaded",t.user.display_name),window.Router.switch("profile",document.getElementById("nav-profile"))):l.error("Failed to load initial data")}document.addEventListener("DOMContentLoaded",f);})();
