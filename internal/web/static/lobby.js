(()=>{function l(){let n=new URLSearchParams(window.location.search).get("id");if(n)return n;let e=window.Telegram?.WebApp?.initDataUnsafe?.user?.id;return e?String(e):"1001"}async function o(){try{let n=l(),e=await fetch(`/api/profile/${n}`),i=await e.json();if(!e.ok||!i.ok)throw new Error(i.error||"User not found");return i.user}catch(n){return console.error("API Error in fetchProfileData:",n),null}}async function a(){let n=document.getElementById("screen-profile");if(!n)return;n.innerHTML='<div style="text-align:center;padding:50px;color:#888;">\u062F\u0631 \u062D\u0627\u0644 \u062F\u0631\u06CC\u0627\u0641\u062A \u0627\u0637\u0644\u0627\u0639\u0627\u062A...</div>';let e=await o();if(!e){n.innerHTML='<div style="text-align:center;padding:50px;color:#ef4444;">\u062E\u0637\u0627 \u062F\u0631 \u062F\u0631\u06CC\u0627\u0641\u062A \u0627\u0637\u0644\u0627\u0639\u0627\u062A \u067E\u0631\u0648\u0641\u0627\u06CC\u0644</div>';return}let i=e.wins||0,r=e.coins||0,t=e.gems||0,d=e.display_name||e.game_username||"\u06A9\u0627\u0631\u0628\u0631 \u0645\u0647\u0645\u0627\u0646";n.innerHTML=`
        <div class="profile-card">
            <div class="banner"></div>
            <div class="profile-avatar">
                <svg viewBox="0 0 160 160"><rect width="160" height="160" fill="#475346"/><path d="M80 13 C27 10 13 43 25 84 C32 108 54 137 80 147 C104 135 127 109 136 80 C149 37 126 10 80 13Z" fill="#8b9472"/><path d="M34 62 Q62 54 71 88 Q40 102 34 62Z" fill="#07100b"/><path d="M127 62 Q99 54 89 88 Q120 102 127 62Z" fill="#07100b"/></svg>
            </div>
            <div class="profile-info"><h2 class="username">${d}</h2></div>
            <div class="metrics">
                <div class="metric"><span>\u0628\u0631\u062F\u0647\u0627</span><strong>${i}</strong></div>
                <div class="metric"><span>\u0633\u06A9\u0647</span><strong>${r}</strong></div>
                <div class="metric"><span>\u0627\u0644\u0645\u0627\u0633</span><strong>${t}</strong></div>
            </div>
        </div>
    `}function c(){let n=document.getElementById("screen-games");if(!n)return;let e=[{id:"ludo",name:"\u0645\u0646\u0686 \u0641\u0636\u0627\u06CC\u06CC",icon:"\u{1F3B2}",route:"/games/ludo/lobby.html"}];n.innerHTML=`
        <div class="games-grid">
            ${e.map(i=>`
                <div class="game-card" onclick="window.location.href='${i.route}'">
                    <div class="game-icon">${i.icon}</div>
                    <div class="game-name">${i.name}</div>
                </div>
            `).join("")}
        </div>
    `}var s=class{static switch(e,i){document.querySelectorAll(".screen").forEach(t=>t.classList.remove("active")),document.querySelectorAll(".nav-item").forEach(t=>t.classList.remove("active"));let r=document.getElementById(`screen-${e}`);r&&r.classList.add("active"),e==="home"?(document.querySelectorAll(".screen").forEach(t=>t.classList.remove("active")),document.getElementById("nav-home")?.classList.add("active")):i&&i.classList.add("active"),e==="profile"?a():e==="games"?c():e==="store"?r&&(r.innerHTML='<div style="text-align:center;padding:50px;color:#888;">\u0641\u0631\u0648\u0634\u06AF\u0627\u0647 \u0628\u0647 \u0632\u0648\u062F\u06CC...</div>'):e==="friends"?r&&(r.innerHTML='<div style="text-align:center;padding:50px;color:#888;">\u062F\u0648\u0633\u062A\u0627\u0646 \u0648 \u06AF\u0631\u0648\u0647\u200C\u0647\u0627 \u0628\u0647 \u0632\u0648\u062F\u06CC...</div>'):e==="rooms"&&r&&(r.innerHTML='<div style="text-align:center;padding:50px;color:#888;">\u0631\u0648\u0645\u200C\u0647\u0627\u06CC \u0639\u0645\u0648\u0645\u06CC \u0628\u0647 \u0632\u0648\u062F\u06CC...</div>')}};window.Telegram&&window.Telegram.WebApp&&(window.Telegram.WebApp.ready(),window.Telegram.WebApp.expand());window.Router=s;document.addEventListener("DOMContentLoaded",()=>{s.switch("home",document.getElementById("nav-home"))});})();
