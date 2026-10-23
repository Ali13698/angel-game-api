(()=>{async function l(e){try{let t=await fetch(`/api/profile/${e}`);return t.ok?(await t.json()).user:null}catch{return null}}var c=e=>`<svg width="28" height="28" viewBox="0 0 24 24" fill="${e?"#ef4444":"none"}" stroke="${e?"#ef4444":"#ffffff"}" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path></svg>`;function o(e,t){let i="";for(let a=1;a<=20;a++)i+=`<img src="https://api.dicebear.com/7.x/bottts/svg?seed=${a}" class="avatar-option" data-url="https://api.dicebear.com/7.x/bottts/svg?seed=${a}" />`;let n=["linear-gradient(135deg, #7c3aed 0%, #ec4899 100%)","linear-gradient(135deg, #06b6d4 0%, #3b82f6 100%)","linear-gradient(135deg, #f59e0b 0%, #ef4444 100%)"],s=n.map(a=>`<div class="banner-option" style="background: ${a};" data-bg="${a}"></div>`).join("");return`
    <div class="app-shell">
        <div style="background:#22c55e; color:white; padding:8px; text-align:center; font-size:13px; font-weight:bold;">\u2705 \u0646\u0633\u062E\u0647 TypeScript \u0641\u0639\u0627\u0644 \u0634\u062F!</div>
        <div class="page-content">
            <div class="profile-card">
                <div class="banner" id="profile-banner" style="background: ${e.banner_id||n[0]}; background-size: cover;">
                    <button class="plus-badge banner-plus" id="btn-edit-banner">+</button>
                </div>
                <div class="avatar-container">
                    <div class="avatar"><img id="profile-avatar-img" src="${e.avatar||"https://api.dicebear.com/7.x/bottts/svg?seed=1"}" /></div>
                    <button class="plus-badge avatar-plus" id="btn-change-avatar">+</button>
                </div>
                <div class="profile-info">
                    <div class="name" style="color:${e.profile_color||"#c084fc"}">${e.display_name||e.game_username}</div>
                    <div class="username">@${e.game_username}</div>
                    <div class="bio-wrapper">
                        <span class="bio" id="profile-bio">${e.bio||"\u0628\u06CC\u0648\u06AF\u0631\u0627\u0641\u06CC..."}</span>
                        <button class="pencil-btn" id="btn-edit-bio">\u270F\uFE0F</button>
                    </div>
                </div>
                <div class="main-stats">
                    <div class="stat-box" id="like-box" style="cursor:pointer;">
                        <div class="stat-value" id="heart-icon">${c(t)}</div>
                        <div class="stat-label" id="likes-count">${e.likes_count||0}</div>
                    </div>
                    <div class="stat-box"><div class="stat-value">${e.age_days||1}</div><div class="stat-label">\u0631\u0648\u0632</div></div>
                    <div class="stat-box"><div class="stat-value">${e.wins||0}</div><div class="stat-label">\u0628\u0631\u062F</div></div>
                </div>
            </div>
        </div>
    </div>

    <!-- Modals -->
    <div id="avatar-modal" class="modal"><div class="modal-content"><h3>\u0627\u0646\u062A\u062E\u0627\u0628 \u0622\u0648\u0627\u062A\u0627\u0631</h3><div class="avatar-grid" id="avatar-grid">${i}</div><button class="modal-btn close" id="close-avatar-btn">\u0627\u0646\u0635\u0631\u0627\u0641</button></div></div>
    <div id="banner-modal" class="modal"><div class="modal-content"><h3>\u0627\u0646\u062A\u062E\u0627\u0628 \u0628\u0646\u0631</h3><div class="banner-grid" id="banner-grid">${s}</div><button class="modal-btn close" id="close-banner-btn">\u0627\u0646\u0635\u0631\u0627\u0641</button></div></div>
    `}var r=null,d=!1;function v(){let t=new URLSearchParams(location.search).get("id");if(t)return parseInt(t,10);let i=window.Telegram?.WebApp;return Number(i?.initDataUnsafe?.user?.id)||1001}async function m(){let e=document.getElementById("app");if(!e)return;e.innerHTML='<div style="text-align:center; padding:50px; color:white;">\u062F\u0631 \u062D\u0627\u0644 \u0628\u0627\u0631\u06AF\u0630\u0627\u0631\u06CC...</div>';let t=v(),i=await l(t);i||(i={id:t,game_username:"user_"+t,display_name:"User "+t,avatar:"",bio:"",likes_count:0,coins:0,gems:0,wins:0,age_days:1}),r=i,e.innerHTML=o(r,d),document.getElementById("like-box")?.addEventListener("click",()=>{d=!d;let n=document.getElementById("heart-icon")?.querySelector("svg"),s=document.getElementById("likes-count");n&&(n.setAttribute("fill",d?"#ef4444":"none"),n.setAttribute("stroke",d?"#ef4444":"#ffffff")),s&&r&&(r.likes_count+=d?1:-1,s.innerText=r.likes_count.toString())}),document.getElementById("btn-change-avatar")?.addEventListener("click",()=>document.getElementById("avatar-modal")?.classList.add("active")),document.getElementById("close-avatar-btn")?.addEventListener("click",()=>document.getElementById("avatar-modal")?.classList.remove("active")),document.getElementById("btn-edit-banner")?.addEventListener("click",()=>document.getElementById("banner-modal")?.classList.add("active")),document.getElementById("close-banner-btn")?.addEventListener("click",()=>document.getElementById("banner-modal")?.classList.remove("active")),document.querySelectorAll(".avatar-option").forEach(n=>{n.addEventListener("click",s=>{let a=s.target.getAttribute("data-url");a&&(document.getElementById("profile-avatar-img").src=a,document.getElementById("avatar-modal")?.classList.remove("active"))})}),document.querySelectorAll(".banner-option").forEach(n=>{n.addEventListener("click",s=>{let a=s.target.getAttribute("data-bg");a&&(document.getElementById("profile-banner").style.background=a,document.getElementById("banner-modal")?.classList.remove("active"))})})}document.addEventListener("DOMContentLoaded",m);})();
