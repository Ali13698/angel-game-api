(()=>{function h(e,i){console.log(`[INFO] ${e}`,i||"")}function c(e,i){console.error(`[ERROR] ${e}`,i||"")}function L(e,i){console.warn(`[WARN] ${e}`,i||"")}async function E(e){try{h("Fetching token for user",e);let t=await(await fetch("/api/auth/token?id="+e)).json();return t.ok?t.token:null}catch(i){return c("Fetch token failed",i),null}}async function I(e){let i=new URLSearchParams(location.search).get("id");if(!i)return{ok:!1,error:"No ID"};let t=await E(parseInt(i,10));if(!t)return{ok:!1,error:"Auth error"};h("Updating profile with patch:",e);try{let a=await fetch("/api/me",{method:"PATCH",headers:{Authorization:"Bearer "+t,"Content-Type":"application/json"},body:JSON.stringify(e)}),n=await a.json();return a.ok||c("Update profile failed on server",n),{ok:a.ok,error:n.error}}catch(a){return c("Update profile network error",a),{ok:!1,error:"Network error"}}}async function M(e){let i=new URLSearchParams(location.search).get("id");if(!i)return{ok:!1,error:"No ID"};let t=await E(parseInt(i,10));if(!t)return{ok:!1,error:"Auth error"};h("Uploading image, size:",e.size);let a=new FormData;a.append("image",e,"avatar.jpg");try{let n=await fetch("/api/me/upload",{method:"POST",headers:{Authorization:"Bearer "+t},body:a}),d=await n.json();return n.ok?h("Upload successful",d.url):c("Upload failed on server",d),{ok:n.ok,url:d.url,error:d.error}}catch(n){return c("Upload network error",n),{ok:!1,error:"Network error"}}}var O='<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"></path><circle cx="12" cy="12" r="3"></circle></svg>',H='<svg width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="19" y1="12" x2="5" y2="12"></line><polyline points="12 19 5 12 12 5"></polyline></svg>',S='<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>',F='<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#8b949e" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"></path><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path></svg>',P='<svg width="28" height="28" viewBox="0 0 24 24" fill="{fill}" stroke="{stroke}" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path></svg>',D={ludo:{name:"\u0645\u0646\u0686",icon:"\u{1F3AF}"},football:{name:"\u0641\u0648\u062A\u0628\u0627\u0644 \u062F\u0633\u062A\u06CC",icon:"\u26BD"},billiards:{name:"\u0628\u06CC\u0644\u06CC\u0627\u0631\u062F",icon:"\u{1F3B1}"}};function m(e){return String(e||"").replace(/[&<>"']/g,i=>({"&":"&amp;","<":"&lt;",">":"&gt;",'"':"&quot;","'":"&#39;"})[i])}function G(e,i,t){let a=e.display_name||e.game_username||"User"+e.id,n=e.avatar?`<img src="${m(e.avatar)}">`:e.name_emoji||"\u{1F3AE}",d=i?"#ef4444":"none",b=i?"#ef4444":"#ffffff",y=P.replace("{fill}",d).replace("{stroke}",b),r=t.length>0?t[0].url:"#15181e",o=e.banner_id?e.banner_id:r;return`
    <div class="profile-card">
      <div class="top-bar">
        <button class="icon-btn" id="btn-open-settings">${O}</button>
        <button class="icon-btn">${H}</button>
      </div>
      <div class="banner" style="background: ${o}; background-size: cover; background-position: center;">
        <button class="plus-badge banner-plus" id="btn-edit-banner">${S}</button>
      </div>
      <div class="avatar-container">
        <div class="avatar" id="btn-edit-avatar">${n}</div>
        <button class="plus-badge avatar-plus" id="btn-change-avatar">${S}</button>
      </div>
      <div class="profile-info">
        <div class="name" style="color:${e.profile_color||"#c084fc"}">${m(a)}</div>
        <div class="username">@${m(e.game_username||"unknown")}</div>
        <div class="bio-wrapper">
            <span class="bio">${m(e.bio||"---")}</span>
            <button class="pencil-btn" id="btn-edit-bio">${F}</button>
        </div>
      </div>
      <div class="main-stats">
        <div class="stat-box" id="like-box" style="cursor:pointer;">
          <div class="stat-value" id="heart-icon">${y}</div>
          <div class="stat-label" id="likes-count">${e.likes_count||0}</div>
        </div>
        <div class="stat-box"><div class="stat-value">${e.age_days||0}</div><div class="stat-label">\u0631\u0648\u0632</div></div>
        <div class="stat-box"><div class="stat-value">${e.wins||0}</div><div class="stat-label">\u0628\u0631\u062F</div></div>
      </div>
      <div class="action-buttons">
        <button class="btn-square-small green"><span class="icon-3d">\u{1F48E}</span> ${e.gems||0}</button>
        <button class="btn-square-small"><span class="icon-3d">\u{1FA99}</span> ${e.coins||0}</button>
      </div>
    </div>`}function q(e){if(!e||e.length===0)return'<div class="game-stats-card"><div style="text-align:center;padding:20px;color:#8b949e;">\u0647\u0646\u0648\u0632 \u0622\u0645\u0627\u0631\u06CC \u062B\u0628\u062A \u0646\u0634\u062F\u0647 \u0627\u0633\u062A.</div></div>';let i='<div class="game-stats-card" style="padding: 0; background: transparent;">';return e.forEach(t=>{let a=D[t.game_id]||{name:t.game_id,icon:"\u{1F3AE}"},n="Beginner";t.level>=80?n="Prime Gamer":t.level>=60?n="Master":t.level>=40?n="Skilled":t.level>=20?n="Captain":t.level>=10?n="Adept":t.level>=5&&(n="Gamer");let d=t.games_played>0?Math.round(t.wins/t.games_played*100):0;i+=`
      <div style="background: #1c2028; border-radius: 12px; padding: 15px; margin-bottom: 12px; display: flex; align-items: center; justify-content: space-between; box-shadow: 0 4px 10px rgba(0,0,0,0.2);">
        <div style="display: flex; align-items: center; gap: 15px; flex: 1;">
          <div style="position: relative; width: 60px; height: 60px; flex-shrink: 0;">
            <div style="width: 100%; height: 100%; background: linear-gradient(135deg, #2a2e38, #1c2028); border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 35px; border: 1px solid rgba(255,255,255,0.1);">
              ${a.icon}
            </div>
            <div style="position: absolute; bottom: -8px; right: -8px; width: 32px; height: 32px; background: #3b82f6; clip-path: polygon(50% 0%, 61% 35%, 98% 35%, 68% 57%, 79% 91%, 50% 70%, 21% 91%, 32% 57%, 2% 35%, 39% 35%); display: flex; align-items: center; justify-content: center; color: white; font-weight: bold; font-size: 13px; filter: drop-shadow(0 2px 3px rgba(0,0,0,0.5));">
              ${t.level}
            </div>
          </div>
          <div style="flex: 1;">
            <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 6px;">
              <span style="font-size: 16px; font-weight: 800; color: #fff;">${a.name}</span>
              <span style="font-size: 13px; color: #8b949e; font-weight: 600;">${n}</span>
            </div>
            <div style="height: 8px; background: #15181e; border-radius: 4px; overflow: hidden; margin-bottom: 8px; border: 1px solid rgba(255,255,255,0.05);">
              <div style="width: ${d}%; height: 100%; background: linear-gradient(90deg, #22c55e, #4ade80); border-radius: 4px; box-shadow: 0 0 8px rgba(74, 222, 128, 0.4);"></div>
            </div>
            <div style="display: flex; justify-content: space-between; font-size: 12px; color: #8b949e;">
              <span>${t.games_played.toLocaleString()} \u0628\u0627\u0632\u06CC \u06A9\u0631\u062F\u0647</span>
              <span>${t.wins.toLocaleString()} \u0628\u0631\u0646\u062F\u0647 \u0634\u062F\u0647</span>
            </div>
          </div>
        </div>
      </div>
    `}),i+"</div>"}function A(e,i,t,a,n){a.innerHTML=`<div class="container">${G(e,n,t)}<div class="section-title" style="margin-bottom: 10px; color:#8b949e;">\u0622\u0645\u0627\u0631 \u0645\u0646</div>${q(i)}</div>`}var $='<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>',V='<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect><circle cx="8.5" cy="8.5" r="1.5"></circle><polyline points="21 15 16 10 5 21"></polyline></svg>',W='<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect><text x="12" y="16" text-anchor="middle" font-size="10" fill="currentColor">GIF</text></svg>',K='<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 4h16v16H4z"></path><path d="M9 4v16"></path><path d="M15 4v16"></path><path d="M4 9h16"></path><path d="M4 15h16"></path></svg>',Z='<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><circle cx="12" cy="10" r="3"></circle><path d="M7 20.662V19a2 2 0 0 1 2-2h6a2 2 0 0 1 2 2v1.662"></path></svg>';function B(){return`
    <div id="avatar-options-modal" class="modal">
      <div class="modal-content" style="padding: 15px; background: #1c2028;">
        <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:15px;">
            <h3 style="margin:0; font-size:16px;">\u062A\u063A\u06CC\u06CC\u0631 \u062A\u0635\u0648\u06CC\u0631</h3>
            <button id="close-avatar-options" style="background:none; border:none; color:white; font-size:20px; cursor:pointer;">\u2716</button>
        </div>
        
        <div class="avatar-options-grid">
            <div class="option-btn locked" id="opt-gif">
                <div class="opt-icon">${W}</div>
                <div class="opt-label">\u0627\u0646\u062A\u062E\u0627\u0628 \u06AF\u06CC\u0641</div>
                <div class="lock-badge">${$} 2</div>
            </div>
            <div class="option-btn active" id="opt-gallery">
                <div class="opt-icon">${V}</div>
                <div class="opt-label">\u06AF\u0627\u0644\u0631\u06CC</div>
            </div>
            <div class="option-btn locked" id="opt-frame">
                <div class="opt-icon">${K}</div>
                <div class="opt-label">\u0642\u0627\u0628</div>
                <div class="lock-badge">${$}</div>
            </div>
            <div class="option-btn active" id="opt-dicebear">
                <div class="opt-icon">${Z}</div>
                <div class="opt-label">\u0622\u0648\u0627\u062A\u0627\u0631</div>
            </div>
        </div>
      </div>
    </div>`}function _(){return`
    <div id="cropper-modal" class="modal">
      <div class="modal-content" style="padding: 15px; background: #1c2028; width: 95%; max-width: 500px;">
        <h3 style="text-align:center; margin-bottom:15px;">\u062A\u0646\u0638\u06CC\u0645 \u0639\u06A9\u0633</h3>
        <div style="width: 100%; height: 300px; background: #000; border-radius: 8px; overflow: hidden;">
            <img id="image-to-crop" style="max-width: 100%; display: block;" />
        </div>
        <div style="display:flex; gap:10px; margin-top:15px;">
            <button class="modal-btn save" id="save-crop-btn">\u062A\u0627\u06CC\u06CC\u062F \u0648 \u0622\u067E\u0644\u0648\u062F</button>
            <button class="modal-btn close" id="close-crop-btn" style="width:40%;">\u0627\u0646\u0635\u0631\u0627\u0641</button>
        </div>
      </div>
    </div>`}function C(e,i){let t="";return i.forEach(a=>{t+=`<img src="${a.url}" class="avatar-option ${e.avatar===a.url?"selected":""}" data-url="${a.url}" />`}),`<div id="dicebear-modal" class="modal"><div class="modal-content"><h3 style="text-align:center; margin-bottom:15px;">\u0627\u0646\u062A\u062E\u0627\u0628 \u0622\u0648\u0627\u062A\u0627\u0631</h3><div class="avatar-grid" id="dicebear-grid">${t}</div><button class="modal-btn save" id="save-dicebear-btn">\u062A\u0627\u06CC\u06CC\u062F</button><button class="modal-btn close" id="close-dicebear-btn">\u0627\u0646\u0635\u0631\u0627\u0641</button></div></div>`}function T(e){return`
    <div id="edit-modal" class="modal">
      <div class="modal-content">
        <h3 style="text-align:center; margin-bottom:15px;">\u0648\u06CC\u0631\u0627\u06CC\u0634 \u0627\u0637\u0644\u0627\u0639\u0627\u062A</h3>
        <label style="font-size:12px; color:#8b949e; margin-bottom:5px; display:block;">\u0646\u0627\u0645 \u06A9\u0627\u0631\u0628\u0631\u06CC (\u0641\u0642\u0637 \u062D\u0631\u0648\u0641 \u0648 \u0627\u0639\u062F\u0627\u062F - \u062D\u062F\u0627\u06A9\u062B\u0631 \u06F1\u06F0 \u06A9\u0627\u0631\u0627\u06A9\u062A\u0631)</label>
        <input type="text" id="username-input" maxlength="10" pattern="[a-zA-Z0-9]+" value="${m(e.game_username||"")}" style="width:100%; background:var(--inner-card-bg); border:1px solid #333; border-radius:8px; color:white; padding:10px; margin-bottom:15px; font-family:inherit;">
        <label style="font-size:12px; color:#8b949e; margin-bottom:5px; display:block;">\u0628\u06CC\u0648\u06AF\u0631\u0627\u0641\u06CC (\u062D\u062F\u0627\u06A9\u062B\u0631 \u06F1\u06F0\u06F0 \u06A9\u0627\u0631\u0627\u06A9\u062A\u0631)</label>
        <textarea id="bio-input" maxlength="100" placeholder="\u0628\u06CC\u0648\u06AF\u0631\u0627\u0641\u06CC \u062C\u062F\u06CC\u062F...">${m(e.bio||"")}</textarea>
        <button class="modal-btn save" id="save-profile-btn">\u0630\u062E\u06CC\u0631\u0647</button>
        <button class="modal-btn close" id="close-edit-btn">\u0627\u0646\u0635\u0631\u0627\u0641</button>
      </div>
    </div>`}function U(e){let i="";return e.forEach(t=>{t.owned&&(i+=`<div class="banner-option" style="background: ${t.url}; background-size: cover; background-position: center;" data-bg="${t.url}"></div>`)}),`<div id="banner-modal" class="modal"><div class="modal-content"><h3 style="text-align:center; margin-bottom:15px;">\u0627\u0646\u062A\u062E\u0627\u0628 \u0628\u0646\u0631</h3><div class="banner-grid" id="banner-grid">${i}</div><button class="modal-btn close" id="close-banner-btn">\u0627\u0646\u0635\u0631\u0627\u0641</button></div></div>`}function j(){return`
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
    </div>`}function v(e,i="error"){L("Showing Toast:",e);let t=document.getElementById("app-toast");t&&t.remove();let a=document.createElement("div");a.id="app-toast",a.className="toast "+i,a.innerText=e,document.body.appendChild(a),setTimeout(()=>a.remove(),3e3)}var k=[{id:"b1",type:"banner",url:"linear-gradient(135deg, #7c3aed 0%, #ec4899 100%)",price:0,owned:!0},{id:"b2",type:"banner",url:"linear-gradient(135deg, #06b6d4 0%, #3b82f6 100%)",price:0,owned:!0},{id:"b3",type:"banner",url:"linear-gradient(135deg, #f59e0b 0%, #ef4444 100%)",price:0,owned:!0},{id:"b4",type:"banner",url:"linear-gradient(135deg, #10b981 0%, #3b82f6 100%)",price:100,owned:!1}],N=[{id:"a1",type:"avatar",url:"https://api.dicebear.com/7.x/adventurer/svg?seed=1",price:0,owned:!0},{id:"a2",type:"avatar",url:"https://api.dicebear.com/7.x/avataaars/svg?seed=2",price:0,owned:!0},{id:"a3",type:"avatar",url:"https://api.dicebear.com/7.x/bottts/svg?seed=3",price:0,owned:!0}];var g=document.getElementById("app"),s=null,z=[],u=null;async function J(){let e=new URLSearchParams(location.search).get("id");if(!e){g&&(g.innerHTML='<div style="text-align:center;padding:20px;">No ID provided</div>');return}try{let t=await(await fetch("/api/profile/"+e)).json();if(!t.ok)throw new Error(t.error||"User not found");s=t.user,z=t.stats||[],g&&s&&R()}catch(i){c("Load profile failed",i),g&&(g.innerHTML='<div style="text-align:center;padding:20px;color:red;direction:rtl;">Error: '+i.message+"</div>")}}function R(){if(!(!s||!g)){try{let e=localStorage.getItem("liked_"+s.id)==="true";A(s,z,k,g,e)}catch(e){c("Main UI render failed",e),g.innerHTML='<div style="color:red;padding:20px;">\u062E\u0637\u0627 \u062F\u0631 \u0646\u0645\u0627\u06CC\u0634 \u067E\u0631\u0648\u0641\u0627\u06CC\u0644</div>';return}try{document.querySelectorAll(".modal").forEach(i=>i.remove());let e=document.createElement("div");e.innerHTML=B()+C(s,N)+T(s)+U(k)+j()+_(),document.body.appendChild(e)}catch(e){c("Modals render failed",e)}try{Y()}catch(e){c("Event attachment failed",e)}}}function Y(){if(!s)return;let e=(r,o,l)=>{let p=document.getElementById(r);p?p.addEventListener(o,l):console.warn("Element not found:",r)},i=document.getElementById("edit-modal");e("btn-edit-bio","click",()=>i?.classList.add("active")),e("close-edit-btn","click",()=>i?.classList.remove("active")),e("save-profile-btn","click",()=>{let r=document.getElementById("username-input")?.value||"",o=document.getElementById("bio-input")?.value||"";if(!/^[a-zA-Z0-9]{1,10}$/.test(r)){v("\u0646\u0627\u0645 \u06A9\u0627\u0631\u0628\u0631\u06CC \u0641\u0642\u0637 \u062D\u0631\u0648\u0641 \u0648 \u0627\u0639\u062F\u0627\u062F (\u062D\u062F\u0627\u06A9\u062B\u0631 \u06F1\u06F0 \u06A9\u0627\u0631\u0627\u06A9\u062A\u0631) \u0628\u0627\u0634\u062F");return}x({game_username:r,bio:o}),i?.classList.remove("active")});let t=document.getElementById("avatar-options-modal"),a=document.getElementById("dicebear-modal"),n=document.getElementById("cropper-modal");e("btn-change-avatar","click",()=>t?.classList.add("active")),e("close-avatar-options","click",()=>t?.classList.remove("active")),e("opt-gif","click",()=>v("\u062E\u0631\u06CC\u062F \u06AF\u06CC\u0641 \u0628\u0647 \u0632\u0648\u062F\u06CC!","error")),e("opt-frame","click",()=>v("\u062E\u0631\u06CC\u062F \u0642\u0627\u0628 \u0628\u0647 \u0632\u0648\u062F\u06CC!","error")),e("opt-dicebear","click",()=>{t?.classList.remove("active"),a?.classList.add("active")});let d=s.avatar||"";e("close-dicebear-btn","click",()=>a?.classList.remove("active")),document.querySelectorAll("#dicebear-grid .avatar-option").forEach(r=>r.addEventListener("click",o=>{document.querySelectorAll("#dicebear-grid .avatar-option").forEach(l=>l.classList.remove("selected")),o.target.classList.add("selected"),d=o.target.getAttribute("data-url")||""})),e("save-dicebear-btn","click",()=>{x({avatar:d}),a?.classList.remove("active")}),e("opt-gallery","click",()=>{t?.classList.remove("active");let r=document.createElement("input");r.type="file",r.accept="image/*",r.onchange=o=>{let l=o.target.files?.[0];if(!l)return;let p=new FileReader;p.onload=f=>{let w=document.getElementById("image-to-crop");w.src=f.target?.result,n?.classList.add("active"),u&&u.destroy(),u=new window.Cropper(w,{aspectRatio:1,viewMode:1,dragMode:"move",autoCropArea:1,background:!1,cropBoxResizable:!1,cropBoxMovable:!1,toggleDragModeOnDblclick:!1})},p.readAsDataURL(l)},r.click()}),e("close-crop-btn","click",()=>{n?.classList.remove("active"),u&&u.destroy()}),e("save-crop-btn","click",()=>{u&&u.getCroppedCanvas({width:256,height:256}).toBlob(async r=>{if(!r)return;n?.classList.remove("active"),v("\u062F\u0631 \u062D\u0627\u0644 \u0622\u067E\u0644\u0648\u062F...","success");let o=await M(r);o.ok&&o.url?(await x({avatar:o.url}),v("\u0622\u0648\u0627\u062A\u0627\u0631 \u0628\u0627 \u0645\u0648\u0641\u0642\u06CC\u062A \u062A\u063A\u06CC\u06CC\u0631 \u06A9\u0631\u062F","success")):v("\u062E\u0637\u0627 \u062F\u0631 \u0622\u067E\u0644\u0648\u062F: "+(o.error||"\u0646\u0627\u0634\u0646\u0627\u062E\u062A\u0647"))},"image/jpeg",.9)});let b=document.getElementById("banner-modal");e("btn-edit-banner","click",()=>b?.classList.add("active")),e("close-banner-btn","click",()=>b?.classList.remove("active")),document.querySelectorAll(".banner-option").forEach(r=>r.addEventListener("click",o=>{let l=o.target.getAttribute("data-bg");l&&(x({banner_id:l}),b?.classList.remove("active"))})),e("like-box","click",()=>{if(!s)return;let o=!(localStorage.getItem("liked_"+s.id)==="true");localStorage.setItem("liked_"+s.id,String(o)),o?s.likes_count=(s.likes_count||0)+1:s.likes_count=Math.max(0,(s.likes_count||0)-1);let l=document.getElementById("heart-icon"),p=document.getElementById("likes-count");if(l){let f=l.querySelector("path");f&&(f.setAttribute("fill",o?"#ef4444":"none"),f.setAttribute("stroke",o?"#ef4444":"#ffffff"))}p&&(p.innerText=s.likes_count.toString())});let y=document.getElementById("settings-modal");e("btn-open-settings","click",()=>y?.classList.add("active")),e("close-settings-btn","click",()=>y?.classList.remove("active"))}async function x(e){if(!s)return;let i=await I(e);i.ok?(Object.assign(s,e),R()):v(i.error||"\u062E\u0637\u0627 \u062F\u0631 \u0630\u062E\u06CC\u0631\u0647 \u0627\u0637\u0644\u0627\u0639\u0627\u062A")}J();})();
