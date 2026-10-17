export function logInfo(msg: string, data?: any) { console.log(`[INFO] ${msg}`, data || ''); }
export function logError(msg: string, err?: any) { console.error(`[ERROR] ${msg}`, err || ''); }
export function logWarn(msg: string, data?: any) { console.warn(`[WARN] ${msg}`, data || ''); }