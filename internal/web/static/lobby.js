(()=>{function h(){let t=new URLSearchParams(window.location.search).get("id");if(t)return t;let e=window.Telegram?.WebApp?.initDataUnsafe?.user?.id;return e?String(e):"1001"}async function g(){try{let t=h(),e=await fetch(`/api/profile/${t}`),a=await e.json();if(!e.ok||!a.ok)throw new Error(a.error||"User not found");return a.user}catch(t){return console.error("API Error in fetchProfileData:",t),null}}function k(){let t=new URLSearchParams(window.location.search).get("id");if(t)return parseInt(t,10);let e=window.Telegram?.WebApp?.initDataUnsafe?.user?.id;return Number(e)||1001}async function b(){try{let t=k();return(await(await fetch(`/api/auth/token?id=${t}`)).json()).token||""}catch{return""}}async function l(t){try{let e=await b();if(!e)throw new Error("No token");return await(await fetch("/api/me",{method:"PATCH",headers:{"Content-Type":"application/json",Authorization:"Bearer "+e},body:JSON.stringify(t)})).json()}catch{return{ok:!1,error:"\u062E\u0637\u0627 \u062F\u0631 \u0627\u062D\u0631\u0627\u0632 \u0647\u0648\u06CC\u062A"}}}async function v(t,e){try{let a=await b();return a?(await(await fetch("/api/profile/like",{method:"POST",headers:{"Content-Type":"application/json",Authorization:"Bearer "+a},body:JSON.stringify({target_id:t,is_liked:e})})).json()).ok:!1}catch{return!1}}function u(t){let e=document.getElementById("modal-bio"),a=document.getElementById("input-bio");document.getElementById("btn-edit-bio")?.addEventListener("click",()=>{a&&(a.value=t?.bio||""),e?.classList.add("active")}),document.getElementById("btn-close-bio")?.addEventListener("click",()=>e?.classList.remove("active")),document.getElementById("btn-save-bio")?.addEventListener("click",async()=>{if(a){let n=await l({bio:a.value});n.ok?window.location.reload():alert("\u062E\u0637\u0627: "+(n.error||"\u0627\u0631\u062A\u0628\u0627\u0637 \u0628\u0627 \u0633\u0631\u0648\u0631 \u0628\u0631\u0642\u0631\u0627\u0631 \u0646\u0634\u062F"))}}),document.getElementById("btn-change-avatar")?.addEventListener("click",()=>document.getElementById("modal-avatar")?.classList.add("active")),document.getElementById("btn-close-avatar")?.addEventListener("click",()=>document.getElementById("modal-avatar")?.classList.remove("active")),document.getElementById("btn-edit-banner")?.addEventListener("click",()=>document.getElementById("modal-banner")?.classList.add("active")),document.getElementById("btn-close-banner")?.addEventListener("click",()=>document.getElementById("modal-banner")?.classList.remove("active")),window.selectAvatar=async n=>{let i=await l({avatar:n});i.ok?window.location.reload():alert("\u062E\u0637\u0627: "+(i.error||""))},window.selectBanner=async n=>{let i=await l({banner_id:n});i.ok?window.location.reload():alert("\u062E\u0637\u0627: "+(i.error||""))}}function m(){let t=["https://api.dicebear.com/7.x/bottts/svg?seed=1","https://api.dicebear.com/7.x/bottts/svg?seed=2","https://api.dicebear.com/7.x/bottts/svg?seed=3","https://api.dicebear.com/7.x/bottts/svg?seed=4","https://api.dicebear.com/7.x/bottts/svg?seed=5"],e=["linear-gradient(135deg, #7c3aed 0%, #ec4899 100%)","linear-gradient(135deg, #3b82f6 0%, #2dd4bf 100%)","linear-gradient(135deg, #f59e0b 0%, #ef4444 100%)"];return`
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
                ${t.map(a=>`<img src="${a}" class="avatar-option" data-url="${a}" onclick="window.selectAvatar('${a}')">`).join("")}
            </div>
            <button class="modal-btn close" id="btn-close-avatar">\u0628\u0633\u062A\u0646</button>
        </div>
    </div>
    <div id="modal-banner" class="modal">
        <div class="modal-content">
            <h3>\u062A\u063A\u06CC\u06CC\u0631 \u0628\u0646\u0631</h3>
            <div class="banner-grid">
                ${e.map(a=>`<div class="banner-option" style="background:${a}" onclick="window.selectBanner('${a}')"></div>`).join("")}
            </div>
            <button class="modal-btn close" id="btn-close-banner">\u0628\u0633\u062A\u0646</button>
        </div>
    </div>`}async function f(){let t=document.getElementById("screen-profile");if(!t)return;t.innerHTML='<div style="text-align:center;padding:50px;color:#888;">\u062F\u0631 \u062D\u0627\u0644 \u062F\u0631\u06CC\u0627\u0641\u062A...</div>';let e=await g();if(!e){t.innerHTML='<div style="text-align:center;padding:50px;color:#ef4444;">\u062E\u0637\u0627 \u062F\u0631 \u062F\u0631\u06CC\u0627\u0641\u062A \u0627\u0637\u0644\u0627\u0639\u0627\u062A</div>';return}let a=e.avatar||"https://api.dicebear.com/7.x/bottts/svg?seed="+e.id,n=document.querySelector("#header-avatar svg");n&&(n.outerHTML=`<img src="${a}" style="width:100%; height:100%; object-fit:cover;">`),t.innerHTML=`
        <div class="profile-card">
            <div class="banner" style="background: ${e.banner_id||"linear-gradient(135deg, #7c3aed 0%, #ec4899 100%)"}; background-size: cover;">
                <button class="plus-badge banner-plus" id="btn-edit-banner">+</button>
            </div>
            <div class="avatar-container">
                <div class="avatar"><img id="profile-avatar-img" src="${a}" /></div>
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
        ${m()}
    `,u(e);let i=!1;document.getElementById("like-box")?.addEventListener("click",async()=>{let o=document.getElementById("heart-icon")?.querySelector("svg"),r=document.getElementById("likes-count"),s=!i;if(o&&(o.setAttribute("fill",s?"#ef4444":"none"),o.setAttribute("stroke",s?"#ef4444":"#ffffff")),r){let c=parseInt(r.innerText);r.innerText=(s?c+1:Math.max(0,c-1)).toString()}await v(e.id,s)?i=s:(o&&(o.setAttribute("fill",i?"#ef4444":"none"),o.setAttribute("stroke",i?"#ef4444":"#ffffff")),r&&(r.innerText=(e.likes_count||0).toString()),alert("\u062E\u0637\u0627 \u062F\u0631 \u062B\u0628\u062A \u0644\u0627\u06CC\u06A9"))})}function p(){let t=document.getElementById("screen-games");t&&(t.innerHTML=`
        <style>
            .angel-games-wrapper {
                display: flex; 
                flex-direction: column; 
                align-items: center; 
                gap: 20px; 
                padding: 10px;
                margin-top: 40px; /* \u0641\u0627\u0635\u0644\u0647 \u0627\u0632 \u0628\u0627\u0644\u0627 */
            }
            .angel-game-card {
                border-radius: 16px; 
                overflow: hidden; 
                cursor: pointer; 
                transition: transform 0.2s; 
                box-sizing: border-box; 
                width: 100%;
                border: 3px solid #8b5cf6;
            }
            .angel-game-card:active { 
                transform: scale(0.98); 
            }
            
            /* --- Ludo Card --- */
            .angel-ludo-card {
                max-width: 600px; 
                aspect-ratio: 3 / 1; 
                display: flex; 
            }
            .angel-ludo-col {
                flex: 1; display: flex; justify-content: center; align-items: center; position: relative;
            }
            .angel-bg-blue   { background-color: #2188d9; }
            .angel-bg-yellow { background-color: #f5be18; }
            .angel-bg-green  { background-color: #26b32b; }
            .angel-bg-red    { background-color: #e52e2e; }
            
            .angel-home-base {
                width: 75%; height: 75%; display: grid; grid-template-columns: 1fr 1fr; grid-template-rows: 1fr 1fr; gap: 15%; padding: 15%; border-radius: 12px; box-sizing: border-box;
            }
            .angel-bg-blue .angel-home-base   { background-color: #1862a8; }
            .angel-bg-yellow .angel-home-base { background-color: #d19f0b; }
            .angel-bg-green .angel-home-base  { background-color: #1a861e; }
            .angel-bg-red .angel-home-base    { background-color: #b82020; }
            
            .angel-token { border-radius: 50%; aspect-ratio: 1 / 1; width: 100%; box-sizing: border-box; }
            .angel-token.ring { border: 6px solid #ffffff; }
            
            .angel-bg-blue .angel-token.solid   { background-color: #2188d9; }
            .angel-bg-yellow .angel-token.solid { background-color: #f5be18; }
            .angel-bg-green .angel-token.solid  { background-color: #26b32b; }
            .angel-bg-red .angel-token.solid    { background-color: #e52e2e; }
            
            /* --- Tic Tac Toe Card (\u0645\u0633\u062A\u0637\u06CC\u0644\u06CC) --- */
            .angel-ttt-card {
                max-width: 600px; 
                aspect-ratio: 3 / 1; 
                background: #1c1e26; 
                display: flex; 
                justify-content: center; 
                align-items: center; 
            }
            .angel-ttt-icon {
                width: 20%; height: auto; aspect-ratio: 1/1; border-radius: 50%; border: 12px solid #e52e2e; box-sizing: border-box;
            }

            /* --- Football Card --- */
            .angel-football-card {
                max-width: 600px; 
                aspect-ratio: 3 / 1; 
                background: transparent; 
                display: flex; 
                justify-content: center; 
                align-items: center;
                overflow: hidden;
            }
            .angel-football-card img {
                width: 100%;
                height: 100%;
                object-fit: cover;
                display: block;
                /* \u0627\u06AF\u0631 \u0639\u06A9\u0633 \u067E\u0633\u200C\u0632\u0645\u06CC\u0646\u0647 \u062F\u0627\u0631\u062F\u060C \u0627\u06CC\u0646 \u062E\u0637 \u0628\u0627\u0639\u062B \u0645\u06CC\u0634\u0647 \u0628\u0627 \u062A\u0645 \u062A\u0627\u0631\u06CC\u06A9 \u0633\u0627\u06CC\u062A \u062A\u0631\u06A9\u06CC\u0628 \u0628\u0634\u0647 */
                mix-blend-mode: screen; 
            }
            
            @media (max-width: 400px) {
                .angel-token.ring { border-width: 4px; }
                .angel-ttt-icon { border-width: 8px; }
            }
        </style>
        
        <div class="angel-games-wrapper">
            <!-- Ludo -->
            <div class="angel-game-card angel-ludo-card" onclick="window.location.href='/games/ludo/lobby.html'">
                <div class="angel-ludo-col angel-bg-blue">
                    <div class="angel-home-base">
                        <div class="angel-token ring"></div><div class="angel-token solid"></div>
                        <div class="angel-token ring"></div><div class="angel-token ring"></div>
                    </div>
                </div>
                <div class="angel-ludo-col angel-bg-yellow">
                    <div class="angel-home-base">
                        <div class="angel-token solid"></div><div class="angel-token solid"></div>
                        <div class="angel-token ring"></div><div class="angel-token ring"></div>
                    </div>
                </div>
                <div class="angel-ludo-col angel-bg-green">
                    <div class="angel-home-base">
                        <div class="angel-token solid"></div><div class="angel-token ring"></div>
                        <div class="angel-token ring"></div><div class="angel-token ring"></div>
                    </div>
                </div>
                <div class="angel-ludo-col angel-bg-red">
                    <div class="angel-home-base">
                        <div class="angel-token ring"></div><div class="angel-token ring"></div>
                        <div class="angel-token ring"></div><div class="angel-token ring"></div>
                    </div>
                </div>
            </div>

            <!-- Tic Tac Toe -->
            <div class="angel-game-card angel-ttt-card" onclick="window.location.href='/games/tic-tac-toe/index.html'">
                <div class="angel-ttt-icon"></div>
            </div>

            <!-- Football -->
            <div class="angel-game-card angel-football-card" onclick="alert('\u0628\u0627\u0632\u06CC \u0641\u0648\u062A\u0628\u0627\u0644 \u0628\u0647 \u0632\u0648\u062F\u06CC \u0627\u0636\u0627\u0641\u0647 \u0645\u06CC\u200C\u0634\u0648\u062F!')">
                <img alt="\u0641\u0648\u062A\u0628\u0627\u0644" src="data:image/webp;base64,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" />
            </div>
        </div>
    `)}var d=class{static switch(e,a){document.querySelectorAll(".screen").forEach(i=>i.classList.remove("active")),document.querySelectorAll(".nav-item").forEach(i=>i.classList.remove("active"));let n=document.getElementById(`screen-${e}`);n&&n.classList.add("active"),e==="home"?(document.querySelectorAll(".screen").forEach(i=>i.classList.remove("active")),document.getElementById("nav-home")?.classList.add("active")):a&&a.classList.add("active"),e==="profile"?f():e==="games"?p():e==="store"?n&&(n.innerHTML='<div style="text-align:center;padding:50px;color:#888;">\u0641\u0631\u0648\u0634\u06AF\u0627\u0647 \u0628\u0647 \u0632\u0648\u062F\u06CC...</div>'):e==="friends"?n&&(n.innerHTML='<div style="text-align:center;padding:50px;color:#888;">\u062F\u0648\u0633\u062A\u0627\u0646 \u0648 \u06AF\u0631\u0648\u0647\u200C\u0647\u0627 \u0628\u0647 \u0632\u0648\u062F\u06CC...</div>'):e==="rooms"&&n&&(n.innerHTML='<div style="text-align:center;padding:50px;color:#888;">\u0631\u0648\u0645\u200C\u0647\u0627\u06CC \u0639\u0645\u0648\u0645\u06CC \u0628\u0647 \u0632\u0648\u062F\u06CC...</div>')}};window.Telegram&&window.Telegram.WebApp&&(window.Telegram.WebApp.ready(),window.Telegram.WebApp.expand());window.Router=d;document.addEventListener("DOMContentLoaded",()=>{d.switch("profile",document.getElementById("nav-profile"))});})();
