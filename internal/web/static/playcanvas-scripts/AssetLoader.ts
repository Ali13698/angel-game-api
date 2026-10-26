// AssetLoader.ts - Async loader for 3D models (glTF) and textures
export class AssetLoader {
    private app: any;

    constructor(app: any) {
        this.app = app;
    }

    /**
     * Loads an asset into PlayCanvas and returns a Promise.
     * @param url Path to the asset (e.g., 'assets/models/token.glb')
     * @param type e.g., 'container', 'texture'
     */
    public load(url: string, type: string = 'container'): Promise<any> {
        return new Promise((resolve, reject) => {
            const asset = new (window as any).pc.Asset(url, type, { url: url });
            
            asset.on('load', () => resolve(asset));
            asset.on('error', (err: any) => {
                console.error(`[AssetLoader] Failed to load ${url}:`, err);
                reject(err);
            });

            this.app.assets.add(asset);
            this.app.assets.load(asset);
        });
    }

    /**
     * Instantiates a loaded glTF container.
     */
    public instantiate(asset: any): any {
        if (!asset || !asset.resource) return null;
        const entity = asset.resource.instantiateRenderEntity();
        return entity;
    }
}
