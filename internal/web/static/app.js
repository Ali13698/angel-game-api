(()=>{var d="";async function g(a){try{let t=await(await fetch(d+"/api/auth/token?id="+a)).json();return t.ok?t.token:null}catch{return null}}async function c(a){let e=new URLSearchParams(location.search).get("id");if(!e)return!1;let t=await g(parseInt(e,10));if(!t)return!1;try{return(await fetch(d+"/api/me",{method:"PATCH",headers:{Authorization:"Bearer "+t,"Content-Type":"application/json"},body:JSON.stringify(a)})).ok}catch{return!1}}var b=[{id:"tank",name:"\u062A\u0627\u0646\u06A9",icon:"\u{1F680}",level:25,rank:"Captain",played:384,won:193,progress:80},{id:"chess",name:"\u0634\u0637\u0631\u0646\u062C",icon:"\u265F\uFE0F",level:16,rank:"Skilled",played:196,won:88,progress:45},{id:"backgammon",name:"\u0628\u06A9\u200C\u06AF\u0645\u0648\u0646",icon:"\u{1F3B2}",level:15,rank:"Skilled",played:142,won:60,progress:90},{id:"ludo",name:"\u0645\u0646\u0686",icon:"\u{1F3AF}",level:7,rank:"Gamer",played:47,won:10,progress:30},{id:"basketball",name:"\u0628\u0633\u06A9\u062A\u0628\u0627\u0644",icon:"\u{1F3C0}",level:6,rank:"Beginner",played:16,won:9,progress:50}];function o(a){return String(a||"").replace(/[&<>"']/g,e=>({"&":"&amp;","<":"&lt;",">":"&gt;",'"':"&quot;","'":"&#39;"})[e])}function r(a,e,t,p){let m=a.display_name||a.game_username||"User"+a.id,u=a.profile_color||"#c084fc",l=`
  <div class="container">
    <div class="profile-card">
      <div class="banner"></div>
      <div class="avatar-container">
        <div class="avatar">${a.avatar?`<img src="${o(a.avatar)}">`:a.name_emoji||"\u{1F3AE}"}</div>
      </div>
      <div class="profile-info">
        <div class="name" style="color:${u}">${o(m)}</div>
        <div class="username">@${o(a.game_username||"unknown")}</div>
        <div class="bio">${o(a.bio||"---")}</div>
      </div>
      
      <div class="main-stats">
        <div class="stat-box"><div class="stat-value">${a.likes_count||0}</div><div class="stat-label">\u0647\u0627\u0644\u0647</div></div>
        <div class="stat-box"><div class="stat-value">${a.age_days||0}</div><div class="stat-label">\u0631\u0648\u0632</div></div>
        <div class="stat-box"><div class="stat-value">${a.wins||0}</div><div class="stat-label">\u0628\u0631\u062F</div></div>
      </div>

      <div class="action-buttons">
        <button class="btn-square blue" id="btn-edit-bio">\u270F\uFE0F</button>
        <button class="btn-square green">\u{1F48E} ${a.gems||0}</button>
        <button class="btn-square">\u{1FA99} ${a.coins||0}</button>
      </div>
      
      <div style="padding: 0 20px 20px;">
         <button class="btn-hype">\u0647\u0627\u06CC\u067E \u{1F680}</button>
      </div>
    </div>

    <div class="section-title" style="margin-bottom: 10px; color:#8b949e;">\u0622\u0645\u0627\u0631 \u0645\u0646</div>
    <div class="game-stats-card">
  `;b.forEach(n=>{l+=`
      <div class="game-item">
        <div class="game-icon">${n.icon}</div>
        <div class="game-details">
          <div class="game-header">
            <div><span class="game-name">${n.name}</span> <span class="game-rank">${n.rank}</span></div>
            <div class="game-level">\u0645\u0631\u062D\u0644\u0647 ${n.level}</div>
          </div>
          <div class="progress-bar"><div class="progress-fill" style="width: ${n.progress}%"></div></div>
          <div class="game-meta">
            <span>${n.played} \u0628\u0627\u0632\u06CC \u06A9\u0631\u062F\u0647</span>
            <span>${n.won} \u0628\u0631\u0646\u062F\u0647 \u0634\u062F\u0647</span>
          </div>
        </div>
        <div class="game-star">${n.level}</div>
      </div>
    `}),l+="</div></div>",e.innerHTML=l,document.getElementById("btn-edit-bio")?.addEventListener("click",t)}var i=document.getElementById("app"),s=null;async function f(){let a=new URLSearchParams(location.search).get("id");if(!a){i&&(i.innerHTML='<div style="text-align:center;padding:20px;">No ID</div>');return}try{let t=await(await fetch("/api/profile/"+a)).json();if(!t.ok){i&&(i.innerHTML='<div style="text-align:center;padding:20px;">Not Found</div>');return}s=t.user,i&&r(s,i,v,()=>{})}catch{i&&(i.innerHTML='<div style="text-align:center;padding:20px;">Error</div>')}}function v(){let a=prompt("\u0628\u06CC\u0648 \u062C\u062F\u06CC\u062F:",s?.bio||"");a!==null&&c({bio:a}).then(e=>{e&&s&&i&&(s.bio=a,r(s,i,v,()=>{}))})}f();})();
