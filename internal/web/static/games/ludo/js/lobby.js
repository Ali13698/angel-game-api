const gameOrigin = 'https://game.nomi98.com';
const start = document.getElementById('start-game-btn');

function report(text) {
  window.startQuickMatch = () => window.alert(text);
  if (start) start.onclick = window.startQuickMatch;

  let message = document.getElementById('ludo-connection-status');
  if (!message) {
    message = document.createElement('p');
    message.id = 'ludo-connection-status';
    message.setAttribute('role', 'status');
    message.style.cssText = 'text-align:center;padding:12px';
    (document.querySelector('main') || document.body).append(message);
  }

  message.textContent = text;
}

report('در حال بارگذاری لابی منچ…');

if (location.origin !== gameOrigin) {
  const target = new URL('/games/ludo/lobby.html', gameOrigin);
  const source = new URLSearchParams(location.hash.slice(1));
  const launch = new URLSearchParams();

  for (const key of [
    'tgWebAppData',
    'tgWebAppVersion',
    'tgWebAppPlatform',
    'tgWebAppThemeParams',
    'tgWebAppStartParam',
  ]) {
    if (source.has(key)) launch.set(key, source.get(key));
  }

  const app = window.Telegram?.WebApp;
  if (app?.initData) launch.set('tgWebAppData', app.initData);

  target.hash = launch.toString();
  location.replace(target.href);
} else {
  try {
    await import('/games/ludo-light/assets/lobby.js');
    document.getElementById('ludo-connection-status')?.remove();
  } catch (error) {
    console.error('Ludo lobby failed', error);
    report(
      'فایل لابی منچ بارگذاری نشد؛ ایمیج یکپارچه باید روی میزبان منتشر شود.'
    );
  }
}
