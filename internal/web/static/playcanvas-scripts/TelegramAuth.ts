// TelegramAuth.ts - Handles secure extraction of Telegram initData

export class TelegramAuth {
    /**
     * Retrieves the raw initData string from Telegram WebApp SDK.
     * Fallback to URL parameter for local testing (development only).
     */
    public static getInitData(): string {
        // 1. Try Telegram SDK
        const tg = (window as any).Telegram?.WebApp;
        if (tg && tg.initData) {
            return tg.initData;
        }

        // 2. Fallback for local development/testing
        const urlParams = new URLSearchParams(window.location.search);
        const testData = urlParams.get('initData');
        if (testData) {
            console.warn("[Auth] Using test initData from URL. DO NOT USE IN PRODUCTION.");
            return testData;
        }

        return "";
    }

    /**
     * Extract User ID safely if needed for UI display.
     * NOTE: This is NOT for authentication, only for visual display.
     */
    public static getUserIdForDisplay(): string {
        const tg = (window as any).Telegram?.WebApp;
        if (tg && tg.initDataUnsafe?.user?.id) {
            return tg.initDataUnsafe.user.id.toString();
        }
        return "0";
    }
}
