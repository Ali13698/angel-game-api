(()=>{function h(){let t=new URLSearchParams(window.location.search).get("id");if(t)return t;let e=window.Telegram?.WebApp?.initDataUnsafe?.user?.id;return e?String(e):"1001"}async function g(){try{let t=h(),e=await fetch(`/api/profile/${t}`),a=await e.json();if(!e.ok||!a.ok)throw new Error(a.error||"User not found");return a.user}catch(t){return console.error("API Error in fetchProfileData:",t),null}}function w(){let t=new URLSearchParams(window.location.search).get("id");if(t)return parseInt(t,10);let e=window.Telegram?.WebApp?.initDataUnsafe?.user?.id;return Number(e)||1001}async function v(){try{let t=w();return(await(await fetch(`/api/auth/token?id=${t}`)).json()).token||""}catch{return""}}async function l(t){try{let e=await v();if(!e)throw new Error("No token");return await(await fetch("/api/me",{method:"PATCH",headers:{"Content-Type":"application/json",Authorization:"Bearer "+e},body:JSON.stringify(t)})).json()}catch{return{ok:!1,error:"\u062E\u0637\u0627 \u062F\u0631 \u0627\u062D\u0631\u0627\u0632 \u0647\u0648\u06CC\u062A"}}}async function b(t,e){try{let a=await v();return a?(await(await fetch("/api/profile/like",{method:"POST",headers:{"Content-Type":"application/json",Authorization:"Bearer "+a},body:JSON.stringify({target_id:t,is_liked:e})})).json()).ok:!1}catch{return!1}}function u(t){let e=document.getElementById("modal-bio"),a=document.getElementById("input-bio");document.getElementById("btn-edit-bio")?.addEventListener("click",()=>{a&&(a.value=t?.bio||""),e?.classList.add("active")}),document.getElementById("btn-close-bio")?.addEventListener("click",()=>e?.classList.remove("active")),document.getElementById("btn-save-bio")?.addEventListener("click",async()=>{if(a){let n=await l({bio:a.value});n.ok?window.location.reload():alert("\u062E\u0637\u0627: "+(n.error||"\u0627\u0631\u062A\u0628\u0627\u0637 \u0628\u0627 \u0633\u0631\u0648\u0631 \u0628\u0631\u0642\u0631\u0627\u0631 \u0646\u0634\u062F"))}}),document.getElementById("btn-change-avatar")?.addEventListener("click",()=>document.getElementById("modal-avatar")?.classList.add("active")),document.getElementById("btn-close-avatar")?.addEventListener("click",()=>document.getElementById("modal-avatar")?.classList.remove("active")),document.getElementById("btn-edit-banner")?.addEventListener("click",()=>document.getElementById("modal-banner")?.classList.add("active")),document.getElementById("btn-close-banner")?.addEventListener("click",()=>document.getElementById("modal-banner")?.classList.remove("active")),window.selectAvatar=async n=>{let i=await l({avatar:n});i.ok?window.location.reload():alert("\u062E\u0637\u0627: "+(i.error||""))},window.selectBanner=async n=>{let i=await l({banner_id:n});i.ok?window.location.reload():alert("\u062E\u0637\u0627: "+(i.error||""))}}function m(){let t=["https://api.dicebear.com/7.x/bottts/svg?seed=1","https://api.dicebear.com/7.x/bottts/svg?seed=2","https://api.dicebear.com/7.x/bottts/svg?seed=3","https://api.dicebear.com/7.x/bottts/svg?seed=4","https://api.dicebear.com/7.x/bottts/svg?seed=5"],e=["linear-gradient(135deg, #7c3aed 0%, #ec4899 100%)","linear-gradient(135deg, #3b82f6 0%, #2dd4bf 100%)","linear-gradient(135deg, #f59e0b 0%, #ef4444 100%)"];return`
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
    `,u(e);let i=!1;document.getElementById("like-box")?.addEventListener("click",async()=>{let o=document.getElementById("heart-icon")?.querySelector("svg"),r=document.getElementById("likes-count"),s=!i;if(o&&(o.setAttribute("fill",s?"#ef4444":"none"),o.setAttribute("stroke",s?"#ef4444":"#ffffff")),r){let c=parseInt(r.innerText);r.innerText=(s?c+1:Math.max(0,c-1)).toString()}await b(e.id,s)?i=s:(o&&(o.setAttribute("fill",i?"#ef4444":"none"),o.setAttribute("stroke",i?"#ef4444":"#ffffff")),r&&(r.innerText=(e.likes_count||0).toString()),alert("\u062E\u0637\u0627 \u062F\u0631 \u062B\u0628\u062A \u0644\u0627\u06CC\u06A9"))})}function p(){let t=document.getElementById("screen-games");t&&(t.innerHTML=`
        <style>
            .angel-games-wrapper {
                display: flex; flex-direction: column; align-items: center; gap: 20px; padding: 10px;
            }
            .angel-game-card {
                border-radius: 16px; overflow: hidden; cursor: pointer; transition: transform 0.2s; box-sizing: border-box; width: 100%;
            }
            .angel-game-card:active { transform: scale(0.98); }
            
            /* Ludo Card */
            .angel-ludo-card {
                max-width: 600px; aspect-ratio: 3 / 1; display: flex; border: 3px solid #8b5cf6;
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
            
            /* TTT Card */
            .angel-ttt-card {
                max-width: 250px; aspect-ratio: 1 / 1; background: #1c1e26; border: 3px solid #8b5cf6;
                display: flex; justify-content: center; align-items: center; align-self: center; width: 60%;
            }
            .angel-ttt-icon {
                width: 40%; height: 40%; border-radius: 50%; border: 12px solid #e52e2e; box-sizing: border-box;
            }
            
            @media (max-width: 400px) {
                .angel-ludo-card { aspect-ratio: 2.5 / 1; }
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
        </div>
    `)}var d=class{static switch(e,a){document.querySelectorAll(".screen").forEach(i=>i.classList.remove("active")),document.querySelectorAll(".nav-item").forEach(i=>i.classList.remove("active"));let n=document.getElementById(`screen-${e}`);n&&n.classList.add("active"),e==="home"?(document.querySelectorAll(".screen").forEach(i=>i.classList.remove("active")),document.getElementById("nav-home")?.classList.add("active")):a&&a.classList.add("active"),e==="profile"?f():e==="games"?p():e==="store"?n&&(n.innerHTML='<div style="text-align:center;padding:50px;color:#888;">\u0641\u0631\u0648\u0634\u06AF\u0627\u0647 \u0628\u0647 \u0632\u0648\u062F\u06CC...</div>'):e==="friends"?n&&(n.innerHTML='<div style="text-align:center;padding:50px;color:#888;">\u062F\u0648\u0633\u062A\u0627\u0646 \u0648 \u06AF\u0631\u0648\u0647\u200C\u0647\u0627 \u0628\u0647 \u0632\u0648\u062F\u06CC...</div>'):e==="rooms"&&n&&(n.innerHTML='<div style="text-align:center;padding:50px;color:#888;">\u0631\u0648\u0645\u200C\u0647\u0627\u06CC \u0639\u0645\u0648\u0645\u06CC \u0628\u0647 \u0632\u0648\u062F\u06CC...</div>')}};window.Telegram&&window.Telegram.WebApp&&(window.Telegram.WebApp.ready(),window.Telegram.WebApp.expand());window.Router=d;document.addEventListener("DOMContentLoaded",()=>{d.switch("profile",document.getElementById("nav-profile"))});})();
