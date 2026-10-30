import test from 'node:test'
import assert from 'node:assert/strict'
import { access, mkdir, mkdtemp, readFile, rm, symlink, writeFile } from 'node:fs/promises'
import { join } from 'node:path'
import { tmpdir } from 'node:os'
import { createHash } from 'node:crypto'
import { prepareHost } from './prepare-host.mjs'

async function fixture(t) {
  const root = await mkdtemp(join(tmpdir(), 'angel-import-'))
  t.after(() => rm(root, { recursive: true, force: true }))
  const source = join(root, 'package')
  const output = join(root, 'overlay')
  const web = join(source, 'web/games/sample')
  await mkdir(join(source, 'games/sample'), { recursive: true })
  await mkdir(join(web, 'assets'), { recursive: true })
  const wasm = Buffer.from([0, 97, 115, 109, 1, 0, 0, 0])
  await writeFile(join(source, 'games/sample/game.wasm'), wasm)
  await writeFile(join(web, 'index.html'), '<h1>Sample</h1>')
  await writeFile(join(web, 'assets/app.js'), 'console.log(1)')
  const manifest = {
    schemaVersion: 1,
    games: [{
      id: 'sample', version: '1.0.0', abi: 'angel-apply-v1',
      wasm: 'games/sample/game.wasm',
      wasmSha256: createHash('sha256').update(wasm).digest('hex'),
      entry: 'web/games/sample/index.html',
    }],
  }
  const save = () => writeFile(
    join(source, 'host-manifest.json'), JSON.stringify(manifest),
  )
  await save()
  return { root, source, output, web, manifest, save }
}

test('creates host overlay with nested assets and manifest', async t => {
  const f = await fixture(t)
  assert.deepEqual(await prepareHost(f.source, f.output), ['sample'])
  await access(join(f.output, 'games/sample/game.wasm'))
  assert.equal(await readFile(
    join(f.output, 'internal/web/static/games/sample/assets/app.js'), 'utf8',
  ), 'console.log(1)')
  assert.deepEqual(JSON.parse(await readFile(
    join(f.output, 'internal/web/static/host-manifest.json'), 'utf8',
  )), f.manifest)
})

test('rejects corrupt WASM before creating output', async t => {
  const f = await fixture(t)
  await writeFile(join(f.source, 'games/sample/game.wasm'), 'bad')
  await assert.rejects(prepareHost(f.source, f.output), /integrity/)
  await assert.rejects(access(f.output), { code: 'ENOENT' })
})

test('rejects duplicate IDs', async t => {
  const f = await fixture(t)
  f.manifest.games.push(structuredClone(f.manifest.games[0]))
  await f.save()
  await assert.rejects(prepareHost(f.source, f.output), /duplicate/)
})

test('rejects path traversal', async t => {
  const f = await fixture(t)
  f.manifest.games[0].wasm = '../secret.wasm'
  await f.save()
  await assert.rejects(prepareHost(f.source, f.output), /Invalid/)
})

test('rejects symlinked web assets', async t => {
  const f = await fixture(t)
  await symlink(join(f.web, 'index.html'), join(f.web, 'link.html'))
  await assert.rejects(prepareHost(f.source, f.output), /Symlink/)
})

test('rejects symlinked game directories', async t => {
  const f = await fixture(t)
  await rm(join(f.source, 'games/sample'), { recursive: true })
  await symlink(f.web, join(f.source, 'games/sample'))
  await assert.rejects(prepareHost(f.source, f.output), /real directory/)
})

test('does not expose server WASM as a web asset', async t => {
  const f = await fixture(t)
  await writeFile(join(f.web, 'leak.wasm'), 'bad')
  await assert.rejects(prepareHost(f.source, f.output), /must not be in web/)
})

test('does not overwrite existing output', async t => {
  const f = await fixture(t)
  await prepareHost(f.source, f.output)
  await assert.rejects(prepareHost(f.source, f.output), { code: 'EEXIST' })
})

test('rejects empty manifest', async t => {
  const f = await fixture(t)
  f.manifest.games = []
  await f.save()
  await assert.rejects(prepareHost(f.source, f.output), /Non-empty/)
})
