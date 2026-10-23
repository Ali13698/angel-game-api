(()=>{function d(t,n){console.error(`[ERROR] ${t}`,n||"")}function c(){let t=new URLSearchParams(location.search).get("id");if(t&&/^\d+$/.test(t))return parseInt(t,10);let e=window.Telegram?.WebApp?.initDataUnsafe?.user?.id;return e?Number(e):1001}function o(t){return document.getElementById(t)}function l(){let t=o("app");!t||o("gameRoom")||(t.insertAdjacentHTML("beforeend",`
    <div class="game-room active" id="gameRoom">
      <div style="background: #22c55e; color: white; padding: 10px; text-align: center; font-weight: bold; font-size: 14px; z-index: 9999;">
        \u2705 \u0627\u062A\u0635\u0627\u0644 \u0628\u0647 \u0647\u0633\u062A\u0647 Go \u0628\u0631\u0642\u0631\u0627\u0631 \u0634\u062F! (\u0646\u0633\u062E\u0647 \u0628\u06CC\u0644\u062F \u062E\u0648\u062F\u06A9\u0627\u0631: v1.0.5)
      </div>
      <div class="game-header">
        <button id="closeRoomBtn" style="background:none;border:none;color:white;">\u062E\u0631\u0648\u062C</button>
        <div id="roomStatus">\u062F\u0631 \u062D\u0627\u0644 \u0633\u0627\u062E\u062A \u0627\u062A\u0627\u0642...</div>
      </div>
      <div class="board-area"><div id="gameStateView">\u062F\u0631 \u0627\u0646\u062A\u0638\u0627\u0631 \u0634\u0631\u0648\u0639...</div></div>
      <div class="chat-area">
        <div class="chat-messages" id="chatMsgs"></div>
        <div class="chat-input-area">
          <input type="text" id="chatInput" placeholder="\u067E\u06CC\u0627\u0645..." />
          <button id="sendChatBtn">\u0627\u0631\u0633\u0627\u0644</button>
        </div>
      </div>
    </div>`),o("closeRoomBtn")?.addEventListener("click",()=>o("gameRoom")?.remove()),o("sendChatBtn")?.addEventListener("click",a),o("chatInput")?.addEventListener("keypress",n=>{n.key==="Enter"&&a()}))}function a(){let t=o("chatInput"),n=t?.value.trim();if(!n)return;let e=window.sendIntent;typeof e=="function"?(e("chat",{text:n}),t.value=""):d("sendIntent \u062F\u0631 \u062F\u0633\u062A\u0631\u0633 \u0646\u06CC\u0633\u062A")}async function u(t){let n=c();if(!n){alert("\u0627\u0628\u062A\u062F\u0627 \u0648\u0627\u0631\u062F \u0634\u0648\u06CC\u062F");return}l();let e=o("roomStatus");try{let i=await fetch("/rooms",{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify({game_id:t,players:[n.toString()]})});if(!i.ok)throw new Error(await i.text());let r=await i.json();e&&(e.innerText=`\u0627\u062A\u0627\u0642: ${r.room_id}`),setTimeout(()=>{let s=window.sendIntent;typeof s=="function"&&s("join",{room_id:r.room_id})},500)}catch(i){e&&(e.innerText=`\u062E\u0637\u0627: ${i.message}`)}}window.joinGame=u;window.sendChat=a;})();
