(()=>{var y=Object.defineProperty;var w=(t,e,i)=>()=>{if(i)throw i[0];try{return t&&(e=t(t=0)),e}catch(n){throw i=[n],n}};var h=(t,e)=>{for(var i in e)y(t,i,{get:e[i],enumerable:!0})};var r={};h(r,{initModals:()=>L});function L(t){let e=document.getElementById("modal-bio"),i=document.getElementById("input-bio"),n=document.getElementById("btn-edit-bio"),s=document.getElementById("btn-close-bio");n&&e&&i&&n.addEventListener("click",()=>{i.value=t?.bio||"",e.classList.add("active")}),s&&e&&s.addEventListener("click",()=>e.classList.remove("active"))}var c=w(()=>{});function x(){let t=new URLSearchParams(window.location.search).get("id");if(t)return t;let e=window.Telegram?.WebApp?.initDataUnsafe?.user?.id;return e?String(e):"1001"}async function o(){try{let t=x(),e=await fetch(`/api/profile/${t}`),i=await e.json();if(!e.ok||!i.ok)throw new Error(i.error||"User not found");return i.user}catch(t){return console.error("API Error in fetchProfileData:",t),null}}async function d(){let t=document.getElementById("screen-profile");if(!t)return;t.innerHTML='<div style="text-align:center;padding:50px;color:#888;">\u062F\u0631 \u062D\u0627\u0644 \u062F\u0631\u06CC\u0627\u0641\u062A \u0627\u0637\u0644\u0627\u0639\u0627\u062A...</div>';let e=await o();if(!e){t.innerHTML='<div style="text-align:center;padding:50px;color:#ef4444;">\u062E\u0637\u0627 \u062F\u0631 \u062F\u0631\u06CC\u0627\u0641\u062A \u0627\u0637\u0644\u0627\u0639\u0627\u062A \u067E\u0631\u0648\u0641\u0627\u06CC\u0644</div>';return}let i=e.wins||0,n=e.coins||0,s=e.gems||0,v=e.likes_count||0,m=e.age_days||1,u=e.display_name||e.game_username||"\u06A9\u0627\u0631\u0628\u0631 \u0645\u0647\u0645\u0627\u0646",f=e.game_username||"unknown",p=e.bio||"\u0628\u06CC\u0648\u06AF\u0631\u0627\u0641\u06CC...",g=e.avatar||"https://api.dicebear.com/7.x/bottts/svg?seed="+e.id;t.innerHTML=`
        <div class="profile-card">
            <div class="banner" style="background: ${e.banner_id||"linear-gradient(135deg, #7c3aed 0%, #ec4899 100%)"}; background-size: cover;">
                <button class="plus-badge banner-plus" id="btn-edit-banner">+</button>
            </div>
            <div class="avatar-container">
                <div class="avatar"><img id="profile-avatar-img" src="${g}" /></div>
                <button class="plus-badge avatar-plus" id="btn-change-avatar">+</button>
            </div>
            <div class="profile-info">
                <div class="name" style="color:${e.profile_color||"#c084fc"}">${u}</div>
                <div class="username">@${f}</div>
                <div class="bio-wrapper">
                    <span class="bio" id="profile-bio">${p}</span>
                    <button class="pencil-btn" id="btn-edit-bio">\u270F\uFE0F</button>
                </div>
            </div>
            <div class="main-stats">
                <div class="stat-box" id="like-box" style="cursor:pointer;">
                    <div class="stat-value" id="heart-icon">
                        <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="#ffffff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path></svg>
                    </div>
                    <div class="stat-label" id="likes-count">${v}</div>
                </div>
                <div class="stat-box"><div class="stat-value">${m}</div><div class="stat-label">\u0631\u0648\u0632</div></div>
                <div class="stat-box"><div class="stat-value">${i}</div><div class="stat-label">\u0628\u0631\u062F</div></div>
            </div>
            <div class="action-buttons">
                <button class="btn-square-small green">\u{1F48E} ${s}</button>
                <button class="btn-square-small">\u{1FA99} ${n}</button>
            </div>
        </div>
    `;let{initModals:b}=await Promise.resolve().then(()=>(c(),r));b(e)}function l(){let t=document.getElementById("screen-games");if(!t)return;let e=[{id:"ludo",name:"\u0645\u0646\u0686 \u0641\u0636\u0627\u06CC\u06CC",icon:"\u{1F3B2}",route:"/games/ludo/lobby.html"}];t.innerHTML=`
        <div class="games-grid">
            ${e.map(i=>`
                <div class="game-card" onclick="window.location.href='${i.route}'">
                    <div class="game-icon">${i.icon}</div>
                    <div class="game-name">${i.name}</div>
                </div>
            `).join("")}
        </div>
    `}var a=class{static switch(e,i){document.querySelectorAll(".screen").forEach(s=>s.classList.remove("active")),document.querySelectorAll(".nav-item").forEach(s=>s.classList.remove("active"));let n=document.getElementById(`screen-${e}`);n&&n.classList.add("active"),e==="home"?(document.querySelectorAll(".screen").forEach(s=>s.classList.remove("active")),document.getElementById("nav-home")?.classList.add("active")):i&&i.classList.add("active"),e==="profile"?d():e==="games"?l():e==="store"?n&&(n.innerHTML='<div style="text-align:center;padding:50px;color:#888;">\u0641\u0631\u0648\u0634\u06AF\u0627\u0647 \u0628\u0647 \u0632\u0648\u062F\u06CC...</div>'):e==="friends"?n&&(n.innerHTML='<div style="text-align:center;padding:50px;color:#888;">\u062F\u0648\u0633\u062A\u0627\u0646 \u0648 \u06AF\u0631\u0648\u0647\u200C\u0647\u0627 \u0628\u0647 \u0632\u0648\u062F\u06CC...</div>'):e==="rooms"&&n&&(n.innerHTML='<div style="text-align:center;padding:50px;color:#888;">\u0631\u0648\u0645\u200C\u0647\u0627\u06CC \u0639\u0645\u0648\u0645\u06CC \u0628\u0647 \u0632\u0648\u062F\u06CC...</div>')}};window.Telegram&&window.Telegram.WebApp&&(window.Telegram.WebApp.ready(),window.Telegram.WebApp.expand());window.Router=a;document.addEventListener("DOMContentLoaded",()=>{a.switch("home",document.getElementById("nav-home"))});})();
