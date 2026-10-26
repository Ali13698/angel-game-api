import * as esbuild from 'esbuild';
import path from 'path';
import { fileURLToPath } from 'url';

// Get the directory name in ES modules
const __dirname = path.dirname(fileURLToPath(import.meta.url));

const isWatch = process.argv.includes('--watch');

const buildOptions = {
    // Use absolute paths to prevent any path resolution issues
    entryPoints: [path.resolve(__dirname, '../../internal/web/static/playcanvas-scripts/App.ts')],
    bundle: true,
    minify: true, // Reduces payload size for fast loading
    sourcemap: isWatch,
    format: 'esm',
    target: ['es2020'],
    outfile: path.resolve(__dirname, '../../internal/web/static/dist/main.js'),
    // ROOT CAUSE FIX: Tell esbuild exactly where the node_modules folder is located
    nodePaths: [path.resolve(__dirname, 'node_modules')],
    logLevel: 'info',
};

async function run() {
    if (isWatch) {
        let ctx = await esbuild.context(buildOptions);
        await ctx.watch();
        console.log('👀 Watching for changes...');
    } else {
        await esbuild.build(buildOptions);
        console.log('✅ Build complete.');
    }
}

run().catch(() => process.exit(1));
