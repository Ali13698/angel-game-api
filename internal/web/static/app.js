(()=>{async function w(t){try{let e=await(await fetch("/api/auth/token?id="+t)).json();return e.ok?e.token:null}catch{return null}}async function L(t){let a=new URLSearchParams(location.search).get("id");if(!a)return!1;let e=await w(parseInt(a,10));if(!e)return!1;try{return(await fetch("/api/me",{method:"PATCH",headers:{Authorization:"Bearer "+e,"Content-Type":"application/json"},body:JSON.stringify(t)})).ok}catch{return!1}}var I='<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"></path><circle cx="12" cy="12" r="3"></circle></svg>',B='<svg width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="19" y1="12" x2="5" y2="12"></line><polyline points="12 19 5 12 12 5"></polyline></svg>',$='<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>',S='<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#8b949e" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"></path><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path></svg>',_='<svg width="28" height="28" viewBox="0 0 24 24" fill="{fill}" stroke="{stroke}" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path></svg>',M={ludo:{name:"\u0645\u0646\u0686",icon:"\u{1F3AF}"},backgammon:{name:"\u0628\u06A9\u200C\u06AF\u0645\u0648\u0646",icon:"\u{1F3B2}"},chess:{name:"\u0634\u0637\u0631\u0646\u062C",icon:"\u265F\uFE0F"},football:{name:"\u0641\u0648\u062A\u0628\u0627\u0644 \u062F\u0633\u062A\u06CC",icon:"\u26BD"},billiards:{name:"\u0628\u06CC\u0644\u06CC\u0627\u0631\u062F",icon:"\u{1F3B1}"}},A=["adventurer","avataaars","bottts","fun-emoji","lorelei","micah","miniavs","notionists","open-peeps","personas"],T=Array.from({length:10},(t,a)=>`https://api.dicebear.com/7.x/${A[a%10]}/svg?seed=Angel${a+1}`),k=["linear-gradient(135deg, #7c3aed 0%, #ec4899 100%)","linear-gradient(135deg, #06b6d4 0%, #3b82f6 100%)","linear-gradient(135deg, #f59e0b 0%, #ef4444 100%)","linear-gradient(135deg, #10b981 0%, #3b82f6 100%)",'url("https://images.unsplash.com/photo-1579546929518-9e396f3cc809?w=600&q=80")','url("https://images.unsplash.com/photo-1557682250-33bd709cbe85?w=600&q=80")','url("https://images.unsplash.com/photo-1558591710-4b4a1ae0f04d?w=600&q=80")','url("https://images.unsplash.com/photo-1618005182384-a83a8bd57fbe?w=600&q=80")','url("https://images.unsplash.com/photo-1508615039623-a25605d2b022?w=600&q=80")','url("https://images.unsplash.com/photo-1534447677768-be436bb09401?w=600&q=80")'];function c(t){return String(t||"").replace(/[&<>"']/g,a=>({"&":"&amp;","<":"&lt;",">":"&gt;",'"':"&quot;","'":"&#39;"})[a])}function H(t,a,e){let i=t.display_name||t.game_username||"User"+t.id,l=t.avatar?`<img src="${c(t.avatar)}">`:t.name_emoji||"\u{1F3AE}",h=a?"#ef4444":"none",m=a?"#ef4444":"#ffffff",o=_.replace("{fill}",h).replace("{stroke}",m),d=t.banner_id?`background: ${t.banner_id};`:`background: ${k[0]};`;return`
    <div class="profile-card">
      <div class="top-bar">
        <button class="icon-btn" id="btn-open-settings">${I}</button>
        <button class="icon-btn">${B}</button>
      </div>
      
      <div class="banner" style="${d}">
        <button class="plus-badge banner-plus" id="btn-edit-banner">${$}</button>
      </div>
      
      <div class="avatar-container">
        <div class="avatar" id="btn-edit-avatar">${l}</div>
        <button class="plus-badge avatar-plus" id="btn-change-avatar">${$}</button>
      </div>
      
      <div class="profile-info">
        <div class="name" style="color:${t.profile_color||"#c084fc"}">${c(i)}</div>
        <div class="username">@${c(t.game_username||"unknown")}</div>
        <div class="bio-wrapper">
            <span class="bio">${c(t.bio||"---")}</span>
            <button class="pencil-btn" id="btn-edit-bio">${S}</button>
        </div>
      </div>
      
      <div class="main-stats">
        <div class="stat-box" id="like-box" style="cursor:pointer;">
          <div class="stat-value" id="heart-icon">${o}</div>
          <div class="stat-label" id="likes-count">${t.likes_count||0}</div>
        </div>
        <div class="stat-box"><div class="stat-value">${t.age_days||0}</div><div class="stat-label">\u0631\u0648\u0632</div></div>
        <div class="stat-box"><div class="stat-value">${t.wins||0}</div><div class="stat-label">\u0628\u0631\u062F</div></div>
      </div>

      <div class="action-buttons">
        <button class="btn-square green">\u{1F48E} ${t.gems||0}</button>
        <button class="btn-square">\u{1FA99} ${t.coins||0}</button>
      </div>
      
      <div style="padding: 0 20px 20px;"><button class="btn-hype">\u0647\u0627\u06CC\u067E \u{1F680}</button></div>
    </div>`}function C(t){if(!t||t.length===0)return'<div class="game-stats-card"><div style="text-align:center;padding:20px;color:#8b949e;">\u0647\u0646\u0648\u0632 \u0622\u0645\u0627\u0631\u06CC \u062B\u0628\u062A \u0646\u0634\u062F\u0647 \u0627\u0633\u062A.</div></div>';let a='<div class="game-stats-card">';return t.forEach(e=>{let i=M[e.game_id]||{name:e.game_id,icon:"\u{1F3AE}"},l=e.games_played>0?Math.round(e.wins/e.games_played*100):0;a+=`<div class="game-item"><div class="game-icon">${i.icon}</div><div class="game-details"><div class="game-header"><div><span class="game-name">${i.name}</span></div><div class="game-level">\u0645\u0631\u062D\u0644\u0647 ${e.level}</div></div><div class="progress-bar"><div class="progress-fill" style="width: ${l}%"></div></div><div class="game-meta"><span>${e.games_played} \u0628\u0627\u0632\u06CC \u06A9\u0631\u062F\u0647</span><span>${e.wins} \u0628\u0631\u0646\u062F\u0647 \u0634\u062F\u0647</span></div></div><div class="game-star">${e.level}</div></div>`}),a+"</div>"}function N(t){let a="";return T.forEach(e=>{a+=`<img src="${e}" class="avatar-option ${t.avatar===e?"selected":""}" data-url="${e}" />`}),`<div id="avatar-modal" class="modal"><div class="modal-content"><h3 style="text-align:center; margin-bottom:15px;">\u0627\u0646\u062A\u062E\u0627\u0628 \u0622\u0648\u0627\u062A\u0627\u0631</h3><div class="avatar-grid" id="avatar-grid">${a}</div><button class="modal-btn save" id="save-avatar-btn">\u062A\u0627\u06CC\u06CC\u062F</button><button class="modal-btn close" id="close-avatar-btn">\u0627\u0646\u0635\u0631\u0627\u0641</button></div></div>`}function j(t){return`<div id="bio-modal" class="modal"><div class="modal-content"><h3 style="text-align:center; margin-bottom:15px;">\u0648\u06CC\u0631\u0627\u06CC\u0634 \u0628\u06CC\u0648\u06AF\u0631\u0627\u0641\u06CC</h3><textarea id="bio-input" placeholder="\u0628\u06CC\u0648\u06AF\u0631\u0627\u0641\u06CC \u062C\u062F\u06CC\u062F...">${c(t.bio||"")}</textarea><button class="modal-btn save" id="save-bio-btn">\u0630\u062E\u06CC\u0631\u0647</button><button class="modal-btn close" id="close-bio-btn">\u0627\u0646\u0635\u0631\u0627\u0641</button></div></div>`}function q(){let t="";return k.forEach((a,e)=>{t+=`<div class="banner-option" style="background: ${a};" data-bg="${a}" data-idx="${e}"></div>`}),`<div id="banner-modal" class="modal"><div class="modal-content"><h3 style="text-align:center; margin-bottom:15px;">\u0627\u0646\u062A\u062E\u0627\u0628 \u0628\u0646\u0631</h3><div class="banner-grid" id="banner-grid">${t}</div><button class="modal-btn close" id="close-banner-btn">\u0627\u0646\u0635\u0631\u0627\u0641</button></div></div>`}function P(){return`
    <div id="settings-modal" class="modal">
      <div class="modal-content" style="padding: 15px;">
        <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:15px;"><h3 style="margin:0;">\u062A\u0646\u0638\u06CC\u0645\u0627\u062A</h3><button id="close-settings-btn" style="background:none; border:none; color:white; font-size:20px; cursor:pointer;">\u2716</button></div>
        <ul class="settings-list">
            <li class="settings-item"><span>\u062D\u0633\u0627\u0628</span> <span class="settings-icon">\u{1F464}</span></li>
            <li class="settings-item"><span>\u062F\u0633\u062A\u06AF\u0627\u0647\u200C\u0647\u0627</span> <span class="settings-icon">\u{1F4F1}</span></li>
            <li class="settings-item"><span>\u0632\u0628\u0627\u0646</span> <span class="settings-icon">\u{1F30D}</span></li>
            <li class="settings-item"><span>\u0628\u0627\u0632\u06CC\u200C\u0647\u0627</span> <span class="settings-icon">\u{1F3AE}</span></li>
            <li class="settings-item"><span>\u0627\u0639\u0644\u0627\u0646\u200C\u0647\u0627</span> <span class="settings-icon">\u{1F514}</span></li>
            <li class="settings-item"><span>\u062D\u0631\u06CC\u0645 \u062E\u0635\u0648\u0635\u06CC</span> <span class="settings-icon">\u{1F512}</span></li>
            <li class="settings-item"><span>\u062A\u0645\u200C\u0647\u0627</span> <span class="settings-icon">\u{1F3A8}</span></li>
            <li class="settings-item"><span>\u062F\u0631\u0628\u0627\u0631\u0647 \u0622\u0646\u062C\u0644</span> <span class="settings-icon">\u24D8</span></li>
            <li class="settings-item"><span>\u062F\u0631\u06CC\u0627\u0641\u062A \u067E\u0634\u062A\u06CC\u0628\u0627\u0646\u06CC</span> <span class="settings-icon">\u2753</span></li>
        </ul>
      </div>
    </div>`}function r(t,a,e,i,l){let h=localStorage.getItem("liked_"+t.id)==="true";e.innerHTML=`
    <div class="container">
      ${H(t,h,()=>{})}
      <div class="section-title" style="margin-bottom: 10px; color:#8b949e;">\u0622\u0645\u0627\u0631 \u0645\u0646</div>
      ${C(a)}
    </div>
    ${N(t)}
    ${j(t)}
    ${q()}
    ${P()}
  `;let m=t.avatar||"",o=document.getElementById("avatar-modal");document.getElementById("btn-change-avatar")?.addEventListener("click",()=>o?.classList.add("active")),document.getElementById("close-avatar-btn")?.addEventListener("click",()=>o?.classList.remove("active")),document.querySelectorAll(".avatar-option").forEach(y=>y.addEventListener("click",p=>{document.querySelectorAll(".avatar-option").forEach(g=>g.classList.remove("selected")),p.target.classList.add("selected"),m=p.target.getAttribute("data-url")||""})),document.getElementById("save-avatar-btn")?.addEventListener("click",()=>{i({avatar:m}),o?.classList.remove("active")});let d=document.getElementById("bio-modal");document.getElementById("btn-edit-bio")?.addEventListener("click",()=>d?.classList.add("active")),document.getElementById("close-bio-btn")?.addEventListener("click",()=>d?.classList.remove("active")),document.getElementById("save-bio-btn")?.addEventListener("click",()=>{i({bio:document.getElementById("bio-input").value}),d?.classList.remove("active")});let E=t.banner_id||k[0],f=document.getElementById("banner-modal");document.getElementById("btn-edit-banner")?.addEventListener("click",()=>f?.classList.add("active")),document.getElementById("close-banner-btn")?.addEventListener("click",()=>f?.classList.remove("active")),document.querySelectorAll(".banner-option").forEach(y=>y.addEventListener("click",p=>{let g=p.target.getAttribute("data-bg");g&&(E=g,i({banner_id:E}),f?.classList.remove("active"))})),document.getElementById("like-box")?.addEventListener("click",l);let x=document.getElementById("settings-modal");document.getElementById("btn-open-settings")?.addEventListener("click",()=>x?.classList.add("active")),document.getElementById("close-settings-btn")?.addEventListener("click",()=>x?.classList.remove("active"))}var n=document.getElementById("app"),s=null,v=[];async function U(){let t=new URLSearchParams(location.search).get("id");if(!t){n&&(n.innerHTML='<div style="text-align:center;padding:20px;">No ID</div>');return}try{let e=await(await fetch("/api/profile/"+t)).json();if(!e.ok){n&&(n.innerHTML='<div style="text-align:center;padding:20px;">Not Found</div>');return}s=e.user,v=e.stats||[],n&&s&&r(s,v,n,b,u)}catch{n&&(n.innerHTML='<div style="text-align:center;padding:20px;">Error</div>')}}async function b(t){if(!s)return;await L(t)&&n?(Object.assign(s,t),r(s,v,n,b,u)):alert("\u062E\u0637\u0627 \u062F\u0631 \u0630\u062E\u06CC\u0631\u0647 \u0627\u0637\u0644\u0627\u0639\u0627\u062A")}function u(){if(!s)return;localStorage.getItem("liked_"+s.id)==="true"?(localStorage.setItem("liked_"+s.id,"false"),r(s,v,n,b,u)):(localStorage.setItem("liked_"+s.id,"true"),r(s,v,n,b,u))}U();})();
