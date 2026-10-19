import { Bot, InlineKeyboard, Keyboard } from 'grammy';
import { DB } from '../db/index.js';
import { BALE_TOKEN } from '../config.js';

export function startBaleBot() {
  if (!BALE_TOKEN) {
    console.warn('[bale] no token');
    return null;
  }

  const bot = new Bot(BALE_TOKEN, {
    client: { apiRoot: 'https://tapi.bale.ai/' }
  });

  bot.command('start', async (ctx) => {
    const userId = ctx.from.id;
    const payload = (ctx.match || '').trim();
    const tg = {
      id: userId,
      username: ctx.from.username,
      first_name: ctx.from.first_name,
      last_name: ctx.from.last_name
    };
    await DB.users.upsertUser(tg);

    if (payload.startsWith('REF_')) {
      const code = payload.slice(4).toUpperCase();
      const r = await DB.economy.attachReferrer(userId, code);
      if (r.ok) {
        await ctx.reply('کد ثبت شد');
      }
    }

    await sendMainMenu(ctx);
  });

  bot.hears('🎮 شروع بازی', async (ctx) => {
    await ctx.reply('بازی:', {
      reply_markup: mainInlineKb()
    });
  });

  bot.command('help', async (ctx) => {
    await ctx.reply('ربات آنجل');
  });

  bot.command('me', async (ctx) => {
    const u = await DB.users.getUser(ctx.from.id);
    if (!u) {
      return ctx.reply('ابتدا /start بزن');
    }
    const text = 'نام: ' + (u.game_username || '-') +
      '\nسکه: ' + u.coins +
      '\nبرد: ' + u.wins +
      '\nباخت: ' + u.losses +
      '\nسطح: ' + u.level;
    await ctx.reply(text);
  });

  bot.catch((err) => {
    console.error(err);
  });

  bot.start({
    onStart: (info) => {
      console.log('[bale] @' + info.username + ' started');
    }
  });

  return bot;
}

function mainInlineKb() {
  return new InlineKeyboard().webApp(
    '🎮 بازی آنجل',
    'https://game.nomi98.com/'
  );
}

async function sendMainMenu(ctx) {
  const kb = new Keyboard().text('🎮 شروع بازی').resized();
  await ctx.reply('سلام!', {
    reply_markup: kb
  });
  await ctx.reply('بازی:', {
    reply_markup: mainInlineKb()
  });
}
