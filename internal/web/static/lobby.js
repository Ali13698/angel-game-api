(()=>{async function o(){let s=document.getElementById("screen-profile");if(s){s.innerHTML='<div style="text-align:center;padding:50px;color:#888;">\u062F\u0631 \u062D\u0627\u0644 \u062F\u0631\u06CC\u0627\u0641\u062A \u0627\u0637\u0644\u0627\u0639\u0627\u062A...</div>';try{let i=await(await fetch("/api/me")).json();if(!i.ok)throw new Error(i.error);let e=i.user,r=e.wins||0,c=e.coins||0,d=e.gems||0,l=e.display_name||e.game_username||"\u06A9\u0627\u0631\u0628\u0631 \u0645\u0647\u0645\u0627\u0646";s.innerHTML=`
            <div class="profile-card">
                <div class="banner"></div>
                <div class="profile-avatar">
                    <svg viewBox="0 0 160 160"><rect width="160" height="160" fill="#475346"/><path d="M80 13 C27 10 13 43 25 84 C32 108 54 137 80 147 C104 135 127 109 136 80 C149 37 126 10 80 13Z" fill="#8b9472"/><path d="M34 62 Q62 54 71 88 Q40 102 34 62Z" fill="#07100b"/><path d="M127 62 Q99 54 89 88 Q120 102 127 62Z" fill="#07100b"/></svg>
                </div>
                <div class="profile-info"><h2 class="username">${l}</h2></div>
                <div class="metrics">
                    <div class="metric"><span>\u0628\u0631\u062F\u0647\u0627</span><strong>${r}</strong></div>
                    <div class="metric"><span>\u0633\u06A9\u0647</span><strong>${c}</strong></div>
                    <div class="metric"><span>\u0627\u0644\u0645\u0627\u0633</span><strong>${d}</strong></div>
                </div>
            </div>
        `}catch{s.innerHTML='<div style="text-align:center;padding:50px;color:#ef4444;">\u062E\u0637\u0627 \u062F\u0631 \u062F\u0631\u06CC\u0627\u0641\u062A \u0627\u0637\u0644\u0627\u0639\u0627\u062A \u067E\u0631\u0648\u0641\u0627\u06CC\u0644</div>'}}}function a(){let s=document.getElementById("screen-games");if(!s)return;let n=[{id:"ludo",name:"\u0645\u0646\u0686 \u0641\u0636\u0627\u06CC\u06CC",icon:"\u{1F3B2}",route:"/games/ludo/lobby.html"}];s.innerHTML=`
        <div class="games-grid">
            ${n.map(i=>`
                <div class="game-card" onclick="window.location.href='${i.route}'">
                    <div class="game-icon">${i.icon}</div>
                    <div class="game-name">${i.name}</div>
                </div>
            `).join("")}
        </div>
    `}var t=class{static switch(n,i){document.querySelectorAll(".screen").forEach(r=>r.classList.remove("active")),document.querySelectorAll(".nav-item").forEach(r=>r.classList.remove("active"));let e=document.getElementById(`screen-${n}`);e&&e.classList.add("active"),n==="home"?(document.querySelectorAll(".screen").forEach(r=>r.classList.remove("active")),document.getElementById("nav-home")?.classList.add("active")):i&&i.classList.add("active"),n==="profile"?o():n==="games"?a():n==="store"?e&&(e.innerHTML='<div style="text-align:center;padding:50px;color:#888;">\u0641\u0631\u0648\u0634\u06AF\u0627\u0647 \u0628\u0647 \u0632\u0648\u062F\u06CC...</div>'):n==="friends"?e&&(e.innerHTML='<div style="text-align:center;padding:50px;color:#888;">\u062F\u0648\u0633\u062A\u0627\u0646 \u0648 \u06AF\u0631\u0648\u0647\u200C\u0647\u0627 \u0628\u0647 \u0632\u0648\u062F\u06CC...</div>'):n==="rooms"&&e&&(e.innerHTML='<div style="text-align:center;padding:50px;color:#888;">\u0631\u0648\u0645\u200C\u0647\u0627\u06CC \u0639\u0645\u0648\u0645\u06CC \u0628\u0647 \u0632\u0648\u062F\u06CC...</div>')}};window.Telegram&&window.Telegram.WebApp&&(window.Telegram.WebApp.ready(),window.Telegram.WebApp.expand());window.Router=t;document.addEventListener("DOMContentLoaded",()=>{t.switch("home",document.getElementById("nav-home"))});})();
