// config.js — case-insensitive env resolver
import os from 'os';

const HOST = String(os.hostname() || '').toLowerCase();

function detectRegion() {
  const explicit = String(process.env.REGION || '').toLowerCase().trim();
  if (explicit === 'ir' || explicit === 'com') return explicit;
  if (/angel-ir|sdl|iran/.test(HOST)) return 'ir';
  return 'com';
}

const REGION = detectRegion();
const RU = REGION.toUpperCase();

// Try all case variants + region suffix
function getEnv(...names) {
  for (const n of names) {
    if (!n) continue;
    const variants = [
      `${n}_${REGION}`, `${n}_${RU}`,
      n, n.toUpperCase(), n.toLowerCase()
    ];
    for (const k of variants) {
      const v = process.env[k];
      if (v && String(v).length > 0) return String(v).trim(); // <--- trim() اضافه شد
    }
  }
  return '';
}

// ---- BOTS ----
export const BOT_TOKEN = getEnv('TELEGRAM', 'telegram', 'BOT_TOKEN', 'token');
export const BALE_TOKEN = getEnv('BALE', 'bale', 'BALE_TOKEN');

// ---- WEBAPP URLs ----
export const WEBAPP_URL = getEnv('WEBAPP') ||
  (REGION === 'ir' ? 'https://nomi98.ir/' : 'https://nomi98.com/');

export const BALE_WEBAPP_URL = getEnv('BALEWEB', 'BALE_WEBAPP_URL') ||
  (REGION === 'ir' ? 'https://nomi98.ir/' : 'https://nomi98.com/');

// ---- SERVER ----
export const PORT = parseInt(process.env.PORT || '3000', 10);

// ---- DATABASE ----
export const DATABASE_URL = getEnv('DATABASE_URL', 'dburl');
const DB_TYPE_RAW = String(process.env.DB_TYPE || (DATABASE_URL ? 'postgres' : 'd1')).toLowerCase();
export const usePg = DB_TYPE_RAW === 'postgres' || DB_TYPE_RAW === 'pg';

// Cloudflare
export const CF_TOKEN   = getEnv('CF', 'CFTOKEN', 'CF_API_TOKEN', 'cftoken');
export const CF_ACCOUNT = getEnv('CFACC', 'CF_ACCOUNT_ID', 'CF_ACCOUNT', 'cfaccount');
export const CF_DB      = getEnv('CFDB', 'CF_D1_ID', 'CF_DB', 'cfdb');

// Postgres
export const PGHOST = process.env.PGHOST || '127.0.0.1';
export const PGPORT = parseInt(process.env.PGPORT || '5432', 10);
export const PGUSER = process.env.PGUSER || 'angel';
export const PGPASSWORD = process.env.PGPASSWORD || '';
export const PGDATABASE = process.env.PGDATABASE || 'angel_db';

// ---- BRIDGE (to game-api Go) ----
export const BRIDGE_URL = process.env.BRIDGE_URL || '';

// ---- REDIS ----
export const REDIS_HOST = process.env.ANGEL_REDIS_BUJ_SERVICE_SERVICE_HOST || '127.0.0.1';
export const REDIS_PORT = parseInt(process.env.ANGEL_REDIS_BUJ_SERVICE_SERVICE_PORT || '6379', 10);
export const REDIS_PASSWORD = process.env.REDIS || '';

// ---- MINIO ----
export const MINIO_ENDPOINT = process.env.ANGEL_CACHE_YTI_SERVICE_SERVICE_HOST || '127.0.0.1';
export const MINIO_PORT = parseInt(process.env.ANGEL_CACHE_YTI_SERVICE_SERVICE_PORT_API || '9000', 10);
export const MINIO_ACCESS_KEY = process.env.MINIO_ROOT_USER || '';
export const MINIO_SECRET_KEY = process.env.MINIO_ROOT_PASSWORD || '';
export const MINIO_USE_SSL = false;

export const DB_TYPE_EXPORT = usePg ? 'postgres' : 'd1';
export const REQUIRED_CHANNEL = process.env.REQUIRED_CHANNEL || null;

export const REFERRAL_REWARD        = 50;
export const DAILY_GAMES_FOR_REWARD = 3;
export const DAILY_REWARD_COINS     = 30;
export const SPIN_COOLDOWN          = 86400;
export const LIKE_COIN_REWARD       = 1;
export const LIKE_DAILY_CAP         = 20;
export const RENAME_COST            = 50;

export const RANKS = [
  { min: 25, name: 'Captain'  }, { min: 20, name: 'Expert'   },
  { min: 15, name: 'Master'   }, { min: 10, name: 'Pro'      },
  { min: 7,  name: 'Gamer'    }, { min: 5,  name: 'Beginner' },
  { min: 3,  name: 'Rookie'   }, { min: 1,  name: 'Newbie'   }
];

export const SPIN_PRIZES = [
  { coins: 5,   w: 30 }, { coins: 10,  w: 25 }, { coins: 25,  w: 20 },
  { coins: 50,  w: 15 }, { coins: 100, w: 7  }, { coins: 200, w: 2  }, { coins: 500, w: 1 }
];

export function rankForLevel(l) {
  for (const r of RANKS) if (l >= r.min) return r.name;
  return 'Newbie';
}

export function xpForNextLevel(l) {
  return Math.floor(100 * Math.pow(l, 1.5));
}

export function randomCode(len = 8) {
  const c = 'ABCDEFGHJKLMNPQRSTUVWXYZ23456789';
  let s = '';
  for (let i = 0; i < len; i++) s += c[Math.floor(Math.random() * c.length)];
  return s;
}

const ADJ  = ['بی‌ادب','بی‌رحم','بی‌فکر','خوش‌شانس','خوش‌قلب','سرزنده','دیوانه','مغرور','مکار','شیطون'];
const NOUN = ['شاهزاده','شوالیه','ستاره','قلب','دنیای','کهکشان','افسانه','قهرمان','جنگجو','جادوگر'];

export function randomGName() {
  const a = ADJ[Math.floor(Math.random() * ADJ.length)];
  const n = NOUN[Math.floor(Math.random() * NOUN.length)];
  return `${a}_${n}${Math.floor(1000 + Math.random() * 9000)}`;
}

// ---- DIAGNOSTIC ----
console.log('[config] hostname:', HOST);
console.log('[config] region:', REGION);
console.log('[config] telegram:', BOT_TOKEN ? '✅' : '❌');
console.log('[config] bale:', BALE_TOKEN ? '✅' : '❌');
console.log('[config] cf:', CF_TOKEN ? '✅' : '❌');
console.log('[config] cfacc:', CF_ACCOUNT ? '✅' : '❌');
console.log('[config] cfdb:', CF_DB ? '✅' : '❌');
console.log('[config] loaded');
