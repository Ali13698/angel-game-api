(()=>{var y=Object.defineProperty;var w=(t,e,n)=>()=>{if(n)throw n[0];try{return t&&(e=t(t=0)),e}catch(i){throw n=[i],i}};var L=(t,e)=>{for(var n in e)y(t,n,{get:e[n],enumerable:!0})};var u={};L(u,{initModals:()=>h});function h(t){let e=document.getElementById("modal-bio"),n=document.getElementById("input-bio"),i=document.getElementById("btn-edit-bio"),a=document.getElementById("btn-close-bio");i&&e&&n&&i.addEventListener("click",()=>{n.value=t?.bio||"",e.classList.add("active")}),a&&e&&a.addEventListener("click",()=>e.classList.remove("active"));let s=document.getElementById("modal-avatar"),c=document.getElementById("btn-change-avatar"),d=document.getElementById("btn-close-avatar");c&&s&&c.addEventListener("click",()=>s.classList.add("active")),d&&s&&d.addEventListener("click",()=>s.classList.remove("active"));let o=document.getElementById("modal-banner"),l=document.getElementById("btn-edit-banner"),v=document.getElementById("btn-close-banner");l&&o&&l.addEventListener("click",()=>o.classList.add("active")),v&&o&&v.addEventListener("click",()=>o.classList.remove("active"))}var f=w(()=>{});function E(){let t=new URLSearchParams(window.location.search).get("id");if(t)return t;let e=window.Telegram?.WebApp?.initDataUnsafe?.user?.id;return e?String(e):"1001"}async function m(){try{let t=E(),e=await fetch(`/api/profile/${t}`),n=await e.json();if(!e.ok||!n.ok)throw new Error(n.error||"User not found");return n.user}catch(t){return console.error("API Error in fetchProfileData:",t),null}}async function g(){let t=document.getElementById("screen-profile");if(!t)return;t.innerHTML='<div style="text-align:center;padding:50px;color:#888;">\u062F\u0631 \u062D\u0627\u0644 \u062F\u0631\u06CC\u0627\u0641\u062A \u0627\u0637\u0644\u0627\u0639\u0627\u062A...</div>';let e=await m();if(!e){t.innerHTML='<div style="text-align:center;padding:50px;color:#ef4444;">\u062E\u0637\u0627 \u062F\u0631 \u062F\u0631\u06CC\u0627\u0641\u062A \u0627\u0637\u0644\u0627\u0639\u0627\u062A \u067E\u0631\u0648\u0641\u0627\u06CC\u0644</div>';return}let n=e.wins||0,i=e.coins||0,a=e.gems||0,s=e.likes_count||0,c=e.age_days||1,d=e.display_name||e.game_username||"\u06A9\u0627\u0631\u0628\u0631 \u0645\u0647\u0645\u0627\u0646",o=e.game_username||"unknown",l=e.bio||"\u0628\u06CC\u0648\u06AF\u0631\u0627\u0641\u06CC...",v=e.avatar||"https://api.dicebear.com/7.x/bottts/svg?seed="+e.id;t.innerHTML=`
        <div class="profile-card">
            <div class="banner" style="background: ${e.banner_id||"linear-gradient(135deg, #7c3aed 0%, #ec4899 100%)"}; background-size: cover;">
                <button class="plus-badge banner-plus" id="btn-edit-banner">+</button>
            </div>
            <div class="avatar-container">
                <div class="avatar"><img id="profile-avatar-img" src="${v}" /></div>
                <button class="plus-badge avatar-plus" id="btn-change-avatar">+</button>
            </div>
            <div class="profile-info">
                <div class="name" style="color:${e.profile_color||"#c084fc"}">${d}</div>
                <div class="username">@${o}</div>
                <div class="bio-wrapper">
                    <span class="bio" id="profile-bio">${l}</span>
                    <button class="pencil-btn" id="btn-edit-bio">\u270F\uFE0F</button>
                </div>
            </div>
            <div class="main-stats">
                <div class="stat-box" id="like-box" style="cursor:pointer;">
                    <div class="stat-value" id="heart-icon">
                        <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="#ffffff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path></svg>
                    </div>
                    <div class="stat-label" id="likes-count">${s}</div>
                </div>
                <div class="stat-box"><div class="stat-value">${c}</div><div class="stat-label">\u0631\u0648\u0632</div></div>
                <div class="stat-box"><div class="stat-value">${n}</div><div class="stat-label">\u0628\u0631\u062F</div></div>
            </div>
            <div class="action-buttons">
                <button class="btn-square-small green">\u{1F48E} ${a}</button>
                <button class="btn-square-small">\u{1FA99} ${i}</button>
            </div>
        </div>
    `;let{initModals:p}=await Promise.resolve().then(()=>(f(),u));p(e)}function b(){let t=document.getElementById("screen-games");if(!t)return;let e=[{id:"ludo",name:"\u0645\u0646\u0686 \u0641\u0636\u0627\u06CC\u06CC",icon:"\u{1F3B2}",route:"/games/ludo/lobby.html"}];t.innerHTML=`
        <div class="games-grid">
            ${e.map(n=>`
                <div class="game-card" onclick="window.location.href='${n.route}'">
                    <div class="game-icon">${n.icon}</div>
                    <div class="game-name">${n.name}</div>
                </div>
            `).join("")}
        </div>
    `}var r=class{static switch(e,n){document.querySelectorAll(".screen").forEach(a=>a.classList.remove("active")),document.querySelectorAll(".nav-item").forEach(a=>a.classList.remove("active"));let i=document.getElementById(`screen-${e}`);i&&i.classList.add("active"),e==="home"?(document.querySelectorAll(".screen").forEach(a=>a.classList.remove("active")),document.getElementById("nav-home")?.classList.add("active")):n&&n.classList.add("active"),e==="profile"?g():e==="games"?b():e==="store"?i&&(i.innerHTML='<div style="text-align:center;padding:50px;color:#888;">\u0641\u0631\u0648\u0634\u06AF\u0627\u0647 \u0628\u0647 \u0632\u0648\u062F\u06CC...</div>'):e==="friends"?i&&(i.innerHTML='<div style="text-align:center;padding:50px;color:#888;">\u062F\u0648\u0633\u062A\u0627\u0646 \u0648 \u06AF\u0631\u0648\u0647\u200C\u0647\u0627 \u0628\u0647 \u0632\u0648\u062F\u06CC...</div>'):e==="rooms"&&i&&(i.innerHTML='<div style="text-align:center;padding:50px;color:#888;">\u0631\u0648\u0645\u200C\u0647\u0627\u06CC \u0639\u0645\u0648\u0645\u06CC \u0628\u0647 \u0632\u0648\u062F\u06CC...</div>')}};window.Telegram&&window.Telegram.WebApp&&(window.Telegram.WebApp.ready(),window.Telegram.WebApp.expand());window.Router=r;document.addEventListener("DOMContentLoaded",()=>{r.switch("home",document.getElementById("nav-home"))});})();
