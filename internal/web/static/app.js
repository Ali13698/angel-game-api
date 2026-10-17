(()=>{async function _(t){try{let e=await(await fetch("/api/auth/token?id="+t)).json();return e.ok?e.token:null}catch{return null}}async function S(t){let n=new URLSearchParams(location.search).get("id");if(!n)return!1;let e=await _(parseInt(n,10));if(!e)return!1;try{return(await fetch("/api/me",{method:"PATCH",headers:{Authorization:"Bearer "+e,"Content-Type":"application/json"},body:JSON.stringify(t)})).ok}catch{return!1}}var T='<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"></path><circle cx="12" cy="12" r="3"></circle></svg>',H='<svg width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="19" y1="12" x2="5" y2="12"></line><polyline points="12 19 5 12 12 5"></polyline></svg>',M='<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>',A='<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#8b949e" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"></path><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path></svg>',C='<svg width="28" height="28" viewBox="0 0 24 24" fill="{fill}" stroke="{stroke}" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path></svg>',j={ludo:{name:"\u0645\u0646\u0686",icon:"\u{1F3AF}"},backgammon:{name:"\u0628\u06A9\u200C\u06AF\u0645\u0648\u0646",icon:"\u{1F3B2}"},chess:{name:"\u0634\u0637\u0631\u0646\u062C",icon:"\u265F\uFE0F"},football:{name:"\u0641\u0648\u062A\u0628\u0627\u0644 \u062F\u0633\u062A\u06CC",icon:"\u26BD"},billiards:{name:"\u0628\u06CC\u0644\u06CC\u0627\u0631\u062F",icon:"\u{1F3B1}"}};function r(t){return String(t||"").replace(/[&<>"']/g,n=>({"&":"&amp;","<":"&lt;",">":"&gt;",'"':"&quot;","'":"&#39;"})[n])}function N(t,n,e){let l=t.display_name||t.game_username||"User"+t.id,c=t.avatar?`<img src="${r(t.avatar)}">`:t.name_emoji||"\u{1F3AE}",d=n?"#ef4444":"none",p=n?"#ef4444":"#ffffff",k=C.replace("{fill}",d).replace("{stroke}",p),E=e.length>0?e[0].url:"#15181e",g=t.banner_id?t.banner_id:E;return`
    <div class="profile-card">
      <div class="top-bar">
        <button class="icon-btn" id="btn-open-settings">${T}</button>
        <button class="icon-btn">${H}</button>
      </div>
      
      <div class="banner" style="background: ${g}; background-size: cover; background-position: center;">
        <button class="plus-badge banner-plus" id="btn-edit-banner">${M}</button>
        <input type="file" id="banner-file-input" accept="image/*" style="display:none;">
      </div>
      
      <div class="avatar-container">
        <div class="avatar" id="btn-edit-avatar">${c}</div>
        <button class="plus-badge avatar-plus" id="btn-change-avatar">${M}</button>
        <input type="file" id="avatar-file-input" accept="image/*" style="display:none;">
      </div>
      
      <div class="profile-info">
        <div class="name" style="color:${t.profile_color||"#c084fc"}">${r(l)}</div>
        <div class="username">@${r(t.game_username||"unknown")}</div>
        <div class="bio-wrapper">
            <span class="bio">${r(t.bio||"---")}</span>
            <button class="pencil-btn" id="btn-edit-bio">${A}</button>
        </div>
      </div>
      
      <div class="main-stats">
        <div class="stat-box" id="like-box" style="cursor:pointer;">
          <div class="stat-value" id="heart-icon">${k}</div>
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
    </div>`}function P(t){if(!t||t.length===0)return'<div class="game-stats-card"><div style="text-align:center;padding:20px;color:#8b949e;">\u0647\u0646\u0648\u0632 \u0622\u0645\u0627\u0631\u06CC \u062B\u0628\u062A \u0646\u0634\u062F\u0647 \u0627\u0633\u062A.</div></div>';let n='<div class="game-stats-card">';return t.forEach(e=>{let l=j[e.game_id]||{name:e.game_id,icon:"\u{1F3AE}"},c=e.games_played>0?Math.round(e.wins/e.games_played*100):0;n+=`<div class="game-item"><div class="game-icon">${l.icon}</div><div class="game-details"><div class="game-header"><div><span class="game-name">${l.name}</span></div><div class="game-level">\u0645\u0631\u062D\u0644\u0647 ${e.level}</div></div><div class="progress-bar"><div class="progress-fill" style="width: ${c}%"></div></div><div class="game-meta"><span>${e.games_played} \u0628\u0627\u0632\u06CC \u06A9\u0631\u062F\u0647</span><span>${e.wins} \u0628\u0631\u0646\u062F\u0647 \u0634\u062F\u0647</span></div></div><div class="game-star">${e.level}</div></div>`}),n+"</div>"}function G(t,n){let e="";return n.forEach(l=>{e+=`<img src="${l.url}" class="avatar-option ${t.avatar===l.url?"selected":""}" data-url="${l.url}" />`}),`<div id="avatar-modal" class="modal"><div class="modal-content"><h3 style="text-align:center; margin-bottom:15px;">\u0627\u0646\u062A\u062E\u0627\u0628 \u0622\u0648\u0627\u062A\u0627\u0631</h3><div class="avatar-grid" id="avatar-grid">${e}</div><button class="modal-btn save" id="save-avatar-btn">\u062A\u0627\u06CC\u06CC\u062F</button><button class="modal-btn close" id="close-avatar-btn">\u0627\u0646\u0635\u0631\u0627\u0641</button></div></div>`}function U(t){return`<div id="bio-modal" class="modal"><div class="modal-content"><h3 style="text-align:center; margin-bottom:15px;">\u0648\u06CC\u0631\u0627\u06CC\u0634 \u0628\u06CC\u0648\u06AF\u0631\u0627\u0641\u06CC</h3><textarea id="bio-input" placeholder="\u0628\u06CC\u0648\u06AF\u0631\u0627\u0641\u06CC \u062C\u062F\u06CC\u062F...">${r(t.bio||"")}</textarea><button class="modal-btn save" id="save-bio-btn">\u0630\u062E\u06CC\u0631\u0647</button><button class="modal-btn close" id="close-bio-btn">\u0627\u0646\u0635\u0631\u0627\u0641</button></div></div>`}function z(t){let n="";return t.forEach(e=>{n+=`<div class="banner-option" style="background: ${e.url}; background-size: cover; background-position: center;" data-bg="${e.url}"></div>`}),`<div id="banner-modal" class="modal"><div class="modal-content"><h3 style="text-align:center; margin-bottom:15px;">\u0627\u0646\u062A\u062E\u0627\u0628 \u0628\u0646\u0631</h3><div class="banner-grid" id="banner-grid">${n}</div><button class="modal-btn close" id="close-banner-btn">\u0627\u0646\u0635\u0631\u0627\u0641</button></div></div>`}function O(){return`
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
    </div>`}function v(t,n,e,l,c,d,p,k){let E=localStorage.getItem("liked_"+t.id)==="true";c.innerHTML=`
    <div class="container">
      ${N(t,E,e)}
      <div class="section-title" style="margin-bottom: 10px; color:#8b949e;">\u0622\u0645\u0627\u0631 \u0645\u0646</div>
      ${P(n)}
    </div>
    ${G(t,l)}
    ${U(t)}
    ${z(e)}
    ${O()}
  `;let g=t.avatar||"",u=document.getElementById("avatar-modal");document.getElementById("btn-change-avatar")?.addEventListener("click",()=>u?.classList.add("active")),document.getElementById("close-avatar-btn")?.addEventListener("click",()=>u?.classList.remove("active")),document.querySelectorAll(".avatar-option").forEach(a=>a.addEventListener("click",o=>{document.querySelectorAll(".avatar-option").forEach(f=>f.classList.remove("selected")),o.target.classList.add("selected"),g=o.target.getAttribute("data-url")||""})),document.getElementById("save-avatar-btn")?.addEventListener("click",()=>{d({avatar:g}),u?.classList.remove("active")});let L=document.getElementById("bio-modal");document.getElementById("btn-edit-bio")?.addEventListener("click",()=>L?.classList.add("active")),document.getElementById("close-bio-btn")?.addEventListener("click",()=>L?.classList.remove("active")),document.getElementById("save-bio-btn")?.addEventListener("click",()=>{d({bio:document.getElementById("bio-input").value}),L?.classList.remove("active")});let b=document.getElementById("banner-modal");document.getElementById("btn-edit-banner")?.addEventListener("click",()=>b?.classList.add("active")),document.getElementById("close-banner-btn")?.addEventListener("click",()=>b?.classList.remove("active")),document.querySelectorAll(".banner-option").forEach(a=>a.addEventListener("click",o=>{let f=o.target.getAttribute("data-bg");f&&(d({banner_id:f}),b?.classList.remove("active"))}));let x=document.getElementById("avatar-file-input"),I=document.getElementById("banner-file-input");document.getElementById("btn-change-avatar")?.addEventListener("contextmenu",a=>{a.preventDefault(),x.click()});let $=document.getElementById("avatar-grid");if($){let a=document.createElement("div");a.className="avatar-option",a.style.cssText="display:flex;align-items:center;justify-content:center;font-size:24px;background:#2a2e38;color:white;cursor:pointer;",a.innerHTML="+",a.onclick=()=>x.click(),$.appendChild(a)}let w=document.getElementById("banner-grid");if(w){let a=document.createElement("div");a.className="banner-option",a.style.cssText="display:flex;align-items:center;justify-content:center;font-size:24px;background:#2a2e38;color:white;cursor:pointer;",a.innerHTML="+",a.onclick=()=>I.click(),w.appendChild(a)}x?.addEventListener("change",a=>{let o=a.target.files?.[0];o&&(p(o,"avatar"),u?.classList.remove("active"))}),I?.addEventListener("change",a=>{let o=a.target.files?.[0];o&&(p(o,"banner"),b?.classList.remove("active"))}),document.getElementById("like-box")?.addEventListener("click",k);let B=document.getElementById("settings-modal");document.getElementById("btn-open-settings")?.addEventListener("click",()=>B?.classList.add("active")),document.getElementById("close-settings-btn")?.addEventListener("click",()=>B?.classList.remove("active"))}var i=document.getElementById("app"),s=null,m=[];async function q(){let t=new URLSearchParams(location.search).get("id");if(!t){i&&(i.innerHTML='<div style="text-align:center;padding:20px;">No ID</div>');return}try{let e=await(await fetch("/api/profile/"+t)).json();if(!e.ok){i&&(i.innerHTML='<div style="text-align:center;padding:20px;">Not Found</div>');return}s=e.user,m=e.stats||[],i&&s&&v(s,m,i,y,h)}catch{i&&(i.innerHTML='<div style="text-align:center;padding:20px;">Error</div>')}}async function y(t){if(!s)return;await S(t)&&i?(Object.assign(s,t),v(s,m,i,y,h)):alert("\u062E\u0637\u0627 \u062F\u0631 \u0630\u062E\u06CC\u0631\u0647 \u0627\u0637\u0644\u0627\u0639\u0627\u062A")}function h(){if(!s)return;localStorage.getItem("liked_"+s.id)==="true"?(localStorage.setItem("liked_"+s.id,"false"),v(s,m,i,y,h)):(localStorage.setItem("liked_"+s.id,"true"),v(s,m,i,y,h))}q();})();
// Force rebuild 1792259401
