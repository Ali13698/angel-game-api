(()=>{var f="";async function h(t){try{let a=await(await fetch(f+"/api/auth/token?id="+t)).json();return a.ok?a.token:null}catch{return null}}async function c(t){let i=new URLSearchParams(location.search).get("id");if(!i)return!1;let a=await h(parseInt(i,10));if(!a)return!1;try{return(await fetch(f+"/api/me",{method:"PATCH",headers:{Authorization:"Bearer "+a,"Content-Type":"application/json"},body:JSON.stringify(t)})).ok}catch{return!1}}var g=["\u{1F409}","\u{1F98A}","\u{1F31F}","\u{1F43C}","\u{1F981}","\u{1F438}","\u{1F43A}","\u{1F984}","\u{1F42F}","\u{1F428}","\u{1F47E}","\u{1F916}","\u{1F9D1}","\u{1F469}","\u{1F9D4}","\u{1F451}"];function s(t){return String(t||"").replace(/[&<>"']/g,i=>({"&":"&amp;","<":"&lt;",">":"&gt;",'"':"&quot;","'":"&#39;"})[i])}function r(t,i,a,v){let u=t.display_name||t.game_username||"User"+t.id,p=t.profile_color||"#8b5cf6",b=t.avatar?`<img src="${s(t.avatar)}">`:t.name_emoji||"\u{1F3AE}";i.innerHTML=`
    <div class="banner"></div>
    <div class="avatar-wrap"><div class="avatar">${b}</div></div>
    <div class="name" style="color:${p}">${s(u)}</div>
    <div class="username">@${s(t.game_username||"unknown")}</div>
    <div class="bio">${s(t.bio||"---")}</div>
    <div class="stats">
      <div><div class="stat-value">${t.coins}</div><div class="stat-label">\u0633\u06A9\u0647</div></div>
      <div><div class="stat-value">${t.gems}</div><div class="stat-label">\u0627\u0644\u0645\u0627\u0633</div></div>
      <div><div class="stat-value">${t.level}</div><div class="stat-label">\u0633\u0637\u062D</div></div>
    </div>
    <button class="btn" id="btn-edit-bio">\u062A\u063A\u06CC\u06CC\u0631 \u0628\u06CC\u0648</button>
    <div class="section">
      <div class="section-title">\u0627\u06CC\u0645\u0648\u062C\u06CC \u0646\u0627\u0645</div>
      <div class="emoji-row">
        ${g.map(o=>`<button class="emoji-btn ${t.name_emoji===o?"active":""}" data-emoji="${o}">${o}</button>`).join("")}
      </div>
    </div>
  `,document.getElementById("btn-edit-bio")?.addEventListener("click",a),document.querySelectorAll(".emoji-btn").forEach(o=>{o.addEventListener("click",()=>{let m=o.dataset.emoji;m&&v(m)})})}var e=document.getElementById("app"),n=null;async function y(){let t=new URLSearchParams(location.search).get("id");if(!t){e&&(e.innerHTML='<div class="status">No ID</div>');return}try{let a=await(await fetch("/api/profile/"+t)).json();if(!a.ok){e&&(e.innerHTML='<div class="status">Not Found</div>');return}n=a.user,e&&r(n,e,l,d)}catch{e&&(e.innerHTML='<div class="status">Error</div>')}}function l(){let t=prompt("\u0628\u06CC\u0648 \u062C\u062F\u06CC\u062F:",n?.bio||"");t!==null&&c({bio:t}).then(i=>{i&&n&&e&&(n.bio=t,r(n,e,l,d))})}function d(t){c({name_emoji:t}).then(i=>{i&&n&&e&&(n.name_emoji=t,r(n,e,l,d))})}y();})();
