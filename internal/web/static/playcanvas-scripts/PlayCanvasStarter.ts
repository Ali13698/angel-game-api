// PlayCanvasStarter.ts - Boots the PlayCanvas Engine
import { Config } from './Config.js';

export class PlayCanvasStarter {
    public app: any; // pc.Application

    public init() {
        // Check if PlayCanvas is loaded globally
        const pc = (window as any).pc;
        if (!pc) {
            console.error("[PlayCanvas] Engine not found! Check index.html script tags.");
            return;
        }

        // Create Canvas dynamically if not exists
        let canvas = document.getElementById('application-canvas') as HTMLCanvasElement;
        if (!canvas) {
            canvas = document.createElement('canvas');
            canvas.id = 'application-canvas';
            document.body.appendChild(canvas);
        }

        // Create PlayCanvas Application
        this.app = new pc.Application(canvas, {
            mouse: new pc.Mouse(canvas),
            touch: new pc.TouchDevice(canvas),
            keyboard: new pc.Keyboard(window)
        });

        // Fill the window and start
        this.app.setCanvasFillMode(pc.FILLMODE_FILL_WINDOW);
        this.app.setCanvasResolution(pc.RESOLUTION_AUTO);
        this.app.start();
        
        window.addEventListener('resize', () => this.app.resizeCanvas());

        // Setup basic lighting and camera
        this.setupScene();
        
        if (Config.DEBUG) console.log("[PlayCanvas] Engine Started.");
    }

    private setupScene() {
        const pc = (window as any).pc;
        this.app.scene.ambientLight = new pc.Color(0.2, 0.2, 0.2);

        // Create an Entity for the board (placeholder)
        const boardEntity = new pc.Entity('LudoBoard');
        boardEntity.addComponent('model', { type: 'box' });
        boardEntity.setLocalScale(10, 0.5, 10);
        boardEntity.setPosition(0, 0, 0);
        
        // Add a default material with Plato colors
        const mat = new pc.StandardMaterial();
        mat.diffuse = new pc.Color(0.1, 0.1, 0.15);
        mat.update();
        boardEntity.model.material = mat;

        this.app.root.addChild(boardEntity);

        // Create Camera
        const camera = new pc.Entity('camera');
        camera.addComponent('camera', { clearColor: new pc.Color(0.05, 0.05, 0.08) });
        camera.setPosition(0, 15, 10);
        camera.lookAt(0, 0, 0);
        this.app.root.addChild(camera);
        
        // Directional Light
        const light = new pc.Entity('light');
        light.addComponent('light');
        light.setEulerAngles(45, 30, 0);
        this.app.root.addChild(light);
    }
}
