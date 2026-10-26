import * as esbuild from 'esbuild';

const isWatch = process.argv.includes('--watch');

const buildOptions = {
    entryPoints: ['../../internal/web/static/playcanvas-scripts/App.ts'],
    bundle: true,
    minify: true, // Reduces payload size for fast loading
    sourcemap: isWatch,
    format: 'esm',
    target: ['es2020'],
    outfile: '../../internal/web/static/dist/main.js',
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
