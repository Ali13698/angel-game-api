(()=>{function y(){let a=new URLSearchParams(window.location.search).get("id");if(a)return a;let e=window.Telegram?.WebApp?.initDataUnsafe?.user?.id;return e?String(e):"1001"}async function m(){try{let a=y(),e=await fetch(`/api/profile/${a}`),t=await e.json();if(!e.ok||!t.ok)throw new Error(t.error||"User not found");return t.user}catch(a){return console.error("API Error in fetchProfileData:",a),null}}function h(){let a=new URLSearchParams(window.location.search).get("id");if(a)return parseInt(a,10);let e=window.Telegram?.WebApp?.initDataUnsafe?.user?.id;return Number(e)||1001}async function v(){try{let a=h();return(await(await fetch(`/api/auth/token?id=${a}`)).json()).token||""}catch{return""}}async function c(a){try{let e=await v();if(!e)throw new Error("No token");return await(await fetch("/api/me",{method:"PATCH",headers:{"Content-Type":"application/json",Authorization:"Bearer "+e},body:JSON.stringify(a)})).json()}catch{return{ok:!1,error:"\u062E\u0637\u0627 \u062F\u0631 \u0627\u062D\u0631\u0627\u0632 \u0647\u0648\u06CC\u062A"}}}async function u(a,e){try{let t=await v();return t?(await(await fetch("/api/profile/like",{method:"POST",headers:{"Content-Type":"application/json",Authorization:"Bearer "+t},body:JSON.stringify({target_id:a,is_liked:e})})).json()).ok:!1}catch{return!1}}function f(a){let e=document.getElementById("modal-bio"),t=document.getElementById("input-bio");document.getElementById("btn-edit-bio")?.addEventListener("click",()=>{t&&(t.value=a?.bio||""),e?.classList.add("active")}),document.getElementById("btn-close-bio")?.addEventListener("click",()=>e?.classList.remove("active")),document.getElementById("btn-save-bio")?.addEventListener("click",async()=>{if(t){let n=await c({bio:t.value});n.ok?window.location.reload():alert("\u062E\u0637\u0627: "+(n.error||"\u0627\u0631\u062A\u0628\u0627\u0637 \u0628\u0627 \u0633\u0631\u0648\u0631 \u0628\u0631\u0642\u0631\u0627\u0631 \u0646\u0634\u062F"))}}),document.getElementById("btn-change-avatar")?.addEventListener("click",()=>document.getElementById("modal-avatar")?.classList.add("active")),document.getElementById("btn-close-avatar")?.addEventListener("click",()=>document.getElementById("modal-avatar")?.classList.remove("active")),document.getElementById("btn-edit-banner")?.addEventListener("click",()=>document.getElementById("modal-banner")?.classList.add("active")),document.getElementById("btn-close-banner")?.addEventListener("click",()=>document.getElementById("modal-banner")?.classList.remove("active")),window.selectAvatar=async n=>{let i=await c({avatar:n});i.ok?window.location.reload():alert("\u062E\u0637\u0627: "+(i.error||""))},window.selectBanner=async n=>{let i=await c({banner_id:n});i.ok?window.location.reload():alert("\u062E\u0637\u0627: "+(i.error||""))}}function b(){let a=["https://api.dicebear.com/7.x/bottts/svg?seed=1","https://api.dicebear.com/7.x/bottts/svg?seed=2","https://api.dicebear.com/7.x/bottts/svg?seed=3","https://api.dicebear.com/7.x/bottts/svg?seed=4","https://api.dicebear.com/7.x/bottts/svg?seed=5"],e=["linear-gradient(135deg, #7c3aed 0%, #ec4899 100%)","linear-gradient(135deg, #3b82f6 0%, #2dd4bf 100%)","linear-gradient(135deg, #f59e0b 0%, #ef4444 100%)"];return`
    <div id="modal-bio" class="modal">
        <div class="modal-content">
            <h3>\u0648\u06CC\u0631\u0627\u06CC\u0634 \u0628\u06CC\u0648\u06AF\u0631\u0627\u0641\u06CC</h3>
            <textarea id="input-bio" placeholder="\u062F\u0631\u0628\u0627\u0631\u0647 \u0645\u0646..."></textarea>
            <button class="modal-btn save" id="btn-save-bio">\u0630\u062E\u06CC\u0631\u0647</button>
            <button class="modal-btn close" id="btn-close-bio">\u0644\u063A\u0648</button>
        </div>
    </div>
    <div id="modal-avatar" class="modal">
        <div class="modal-content">
            <h3>\u062A\u063A\u06CC\u06CC\u0631 \u0622\u0648\u0627\u062A\u0627\u0631</h3>
            <div class="avatar-grid">
                ${a.map(t=>`<img src="${t}" class="avatar-option" data-url="${t}" onclick="window.selectAvatar('${t}')">`).join("")}
            </div>
            <button class="modal-btn close" id="btn-close-avatar">\u0628\u0633\u062A\u0646</button>
        </div>
    </div>
    <div id="modal-banner" class="modal">
        <div class="modal-content">
            <h3>\u062A\u063A\u06CC\u06CC\u0631 \u0628\u0646\u0631</h3>
            <div class="banner-grid">
                ${e.map(t=>`<div class="banner-option" style="background:${t}" onclick="window.selectBanner('${t}')"></div>`).join("")}
            </div>
            <button class="modal-btn close" id="btn-close-banner">\u0628\u0633\u062A\u0646</button>
        </div>
    </div>`}async function p(){let a=document.getElementById("screen-profile");if(!a)return;a.innerHTML='<div style="text-align:center;padding:50px;color:#888;">\u062F\u0631 \u062D\u0627\u0644 \u062F\u0631\u06CC\u0627\u0641\u062A...</div>';let e=await m();if(!e){a.innerHTML='<div style="text-align:center;padding:50px;color:#ef4444;">\u062E\u0637\u0627 \u062F\u0631 \u062F\u0631\u06CC\u0627\u0641\u062A \u0627\u0637\u0644\u0627\u0639\u0627\u062A</div>';return}let t=e.avatar||"https://api.dicebear.com/7.x/bottts/svg?seed="+e.id,n=document.querySelector("#header-avatar svg");n&&(n.outerHTML=`<img src="${t}" style="width:100%; height:100%; object-fit:cover;">`),a.innerHTML=`
        <div class="profile-card">
            <div class="banner" style="background: ${e.banner_id||"linear-gradient(135deg, #7c3aed 0%, #ec4899 100%)"}; background-size: cover;">
                <button class="plus-badge banner-plus" id="btn-edit-banner">+</button>
            </div>
            <div class="avatar-container">
                <div class="avatar"><img id="profile-avatar-img" src="${t}" /></div>
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
                    <div class="stat-value" id="heart-icon">
                        <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="#ffffff" stroke-width="2"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path></svg>
                    </div>
                    <div class="stat-label" id="likes-count">${e.likes_count||0}</div>
                </div>
                <div class="stat-box"><div class="stat-value">${e.age_days||1}</div><div class="stat-label">\u0631\u0648\u0632</div></div>
                <div class="stat-box"><div class="stat-value">${e.wins||0}</div><div class="stat-label">\u0628\u0631\u062F</div></div>
            </div>
            <div class="action-buttons">
                <button class="btn-square-small green">\u{1F48E} ${e.gems||0}</button>
                <button class="btn-square-small">\u{1FA99} ${e.coins||0}</button>
            </div>
        </div>
        ${b()}
    `,f(e);let i=!1;document.getElementById("like-box")?.addEventListener("click",async()=>{let o=document.getElementById("heart-icon")?.querySelector("svg"),s=document.getElementById("likes-count"),r=!i;if(o&&(o.setAttribute("fill",r?"#ef4444":"none"),o.setAttribute("stroke",r?"#ef4444":"#ffffff")),s){let l=parseInt(s.innerText);s.innerText=(r?l+1:Math.max(0,l-1)).toString()}await u(e.id,r)?i=r:(o&&(o.setAttribute("fill",i?"#ef4444":"none"),o.setAttribute("stroke",i?"#ef4444":"#ffffff")),s&&(s.innerText=(e.likes_count||0).toString()),alert("\u062E\u0637\u0627 \u062F\u0631 \u062B\u0628\u062A \u0644\u0627\u06CC\u06A9"))})}function g(){let a=document.getElementById("screen-games");if(!a)return;let e=[{id:"ludo",name:"\u0645\u0646\u0686 \u0641\u0636\u0627\u06CC\u06CC",icon:"\u{1F3B2}",route:"/games/ludo/lobby.html"},{id:"tic-tac-toe",name:"\u062F\u0648\u0632",icon:"\u2B55",route:"/games/tic-tac-toe/index.html"}];a.innerHTML=`
        <div class="games-grid">
            ${e.map(t=>`
                <div class="game-card" onclick="window.location.href='${t.route}'">
                    <div class="game-icon">${t.icon}</div>
                    <div class="game-name">${t.name}</div>
                </div>
            `).join("")}
        </div>
    `}var d=class{static switch(e,t){document.querySelectorAll(".screen").forEach(i=>i.classList.remove("active")),document.querySelectorAll(".nav-item").forEach(i=>i.classList.remove("active"));let n=document.getElementById(`screen-${e}`);n&&n.classList.add("active"),e==="home"?(document.querySelectorAll(".screen").forEach(i=>i.classList.remove("active")),document.getElementById("nav-home")?.classList.add("active")):t&&t.classList.add("active"),e==="profile"?p():e==="games"?g():e==="store"?n&&(n.innerHTML='<div style="text-align:center;padding:50px;color:#888;">\u0641\u0631\u0648\u0634\u06AF\u0627\u0647 \u0628\u0647 \u0632\u0648\u062F\u06CC...</div>'):e==="friends"?n&&(n.innerHTML='<div style="text-align:center;padding:50px;color:#888;">\u062F\u0648\u0633\u062A\u0627\u0646 \u0648 \u06AF\u0631\u0648\u0647\u200C\u0647\u0627 \u0628\u0647 \u0632\u0648\u062F\u06CC...</div>'):e==="rooms"&&n&&(n.innerHTML='<div style="text-align:center;padding:50px;color:#888;">\u0631\u0648\u0645\u200C\u0647\u0627\u06CC \u0639\u0645\u0648\u0645\u06CC \u0628\u0647 \u0632\u0648\u062F\u06CC...</div>')}};window.Telegram&&window.Telegram.WebApp&&(window.Telegram.WebApp.ready(),window.Telegram.WebApp.expand());window.Router=d;document.addEventListener("DOMContentLoaded",()=>{d.switch("home",document.getElementById("nav-home"))});})();
