(()=>{var l="",u=["\u{1F409}","\u{1F98A}","\u{1F31F}","\u{1F43C}","\u{1F981}","\u{1F438}","\u{1F43A}","\u{1F984}","\u{1F42F}","\u{1F428}","\u{1F47E}","\u{1F916}","\u{1F9D1}","\u{1F469}","\u{1F9D4}","\u{1F451}"],s=null,o=null;function m(){let t=new URLSearchParams(location.search).get("id");if(!t)return null;let i=parseInt(t,10);return isNaN(i)?null:i}async function b(t){try{let n=await(await fetch(`${l}/api/auth/token?id=${t}`)).json();return n.ok&&n.token?n.token:null}catch{return null}}async function f(t){return(await fetch(`${l}/api/profile/${t}`)).json()}async function v(t){return o?(await fetch(`${l}/api/me`,{method:"PATCH",headers:{Authorization:`Bearer ${o}`,"Content-Type":"application/json"},body:JSON.stringify(t)})).json():{ok:!1,error:"no token"}}function d(t){let i=document.getElementById("app");if(!i)return;let n=t.display_name||t.first_name||t.game_username||`User${t.id}`,e=t.avatar?`<img src="${a(t.avatar)}" alt="">`:t.name_emoji||"\u{1F3AE}",c=t.profile_color||"#8b5cf6";i.innerHTML=`
    <div class="banner" style="background:${a(t.banner_id?"linear-gradient(135deg,#7c3aed,#ec4899)":c)}"></div>
    <div class="avatar-wrap">
      <div class="avatar">${e}</div>
    </div>
    <div class="name" style="color:${c}">${a(n)}</div>
    <div class="username">@${a(t.game_username||"unknown")}</div>
    <div class="bio">${a(t.bio||"\u0628\u062F\u0648\u0646 \u0628\u06CC\u0648")}</div>

    <div class="stats">
      <div><div class="stat-value">${t.wins}</div><div class="stat-label">\u0628\u0631\u062F</div></div>
      <div><div class="stat-value">${t.likes_count}</div><div class="stat-label">\u0647\u0627\u0644\u0647</div></div>
      <div><div class="stat-value">${t.age_days}</div><div class="stat-label">\u0631\u0648\u0632</div></div>
    </div>

    <div class="stats">
      <div><div class="stat-value">${t.coins}</div><div class="stat-label">\u0633\u06A9\u0647</div></div>
      <div><div class="stat-value">${t.gems}</div><div class="stat-label">\u0627\u0644\u0645\u0627\u0633</div></div>
      <div><div class="stat-value">${t.level}</div><div class="stat-label">\u0633\u0637\u062D</div></div>
    </div>

    <div class="actions">
      <button class="btn" id="btn-edit-bio">\u062A\u063A\u06CC\u06CC\u0631 \u0628\u06CC\u0648</button>
      <button class="btn secondary" id="btn-edit-emoji">\u062A\u063A\u06CC\u06CC\u0631 \u0627\u06CC\u0645\u0648\u062C\u06CC</button>
    </div>

    <div class="section">
      <div class="section-title">\u0627\u06CC\u0645\u0648\u062C\u06CC \u0646\u0627\u0645</div>
      <div class="emoji-row" id="emoji-row">
        ${u.map(r=>`<button class="emoji-btn${t.name_emoji===r?" active":""}" data-emoji="${r}">${r}</button>`).join("")}
      </div>
    </div>
  `,g()}function g(){document.getElementById("btn-edit-bio")?.addEventListener("click",async()=>{let i=prompt("\u0628\u06CC\u0648 \u062C\u062F\u06CC\u062F:",s?.bio||"");if(i===null)return;let n=await v({bio:i});n.ok&&n.user?(s=n.user,d(n.user)):alert("\u062E\u0637\u0627: "+(n.error||"unknown"))}),document.querySelectorAll(".emoji-btn").forEach(i=>{i.addEventListener("click",async()=>{let n=i.dataset.emoji;if(!n)return;let e=await v({name_emoji:n});e.ok&&e.user&&(s=e.user,d(e.user))})})}function a(t){return String(t).replace(/[&<>"']/g,i=>({"&":"&amp;","<":"&lt;",">":"&gt;",'"':"&quot;","'":"&#39;"})[i])}async function p(){let t=m(),i=document.getElementById("app");if(!t){i&&(i.innerHTML='<div class="status error">?id= \u062F\u0631 URL \u0644\u0627\u0632\u0645\u0647 \u2014 \u0645\u062B\u0644\u0627\u064B /profile/?id=1</div>');return}o=await b(t);let n=await f(t);if(!n.ok||!n.user){i&&(i.innerHTML=`<div class="status error">${a(n.error||"not found")}</div>`);return}s=n.user,d(n.user)}p();})();
