// UIManager.ts - Handles HTML overlay (Buttons, Loading, Errors)
export class UIManager {
    private loadingOverlay: HTMLElement | null;

    constructor() {
        this.loadingOverlay = document.getElementById('loading-overlay');
    }

    public hideLoading() {
        if (this.loadingOverlay) {
            this.loadingOverlay.classList.add('hidden');
        }
    }

    public showLoading(text: string = "در حال اتصال...") {
        if (this.loadingOverlay) {
            this.loadingOverlay.classList.remove('hidden');
            const textEl = this.loadingOverlay.querySelector('.loading-text');
            if (textEl) textEl.textContent = text;
        }
    }

    public showError(message: string) {
        // Basic non-blocking error popup
        const errDiv = document.createElement('div');
        errDiv.style.cssText = "position:fixed;top:20px;left:50%;transform:translateX(-50%);background:#ef4444;color:#fff;padding:10px 20px;border-radius:8px;z-index:10000;font-family:sans-serif;font-size:14px;box-shadow:0 4px 12px rgba(0,0,0,0.5);";
        errDiv.textContent = message;
        document.body.appendChild(errDiv);

        setTimeout(() => {
            errDiv.style.opacity = '0';
            errDiv.style.transition = 'opacity 0.5s';
            setTimeout(() => errDiv.remove(), 500);
        }, 3000);
    }
}
