import { lstat, mkdir, readFile, readdir, realpath, writeFile } from 'node:fs/promises'
import { dirname, join, resolve } from 'node:path'
import { createHash } from 'node:crypto'
import { fileURLToPath } from 'node:url'

const header = Buffer.from([0, 97, 115, 109, 1, 0, 0, 0])
const sha256 = bytes => createHash('sha256').update(bytes).digest('hex')

async function regular(path) {
  const stat = await lstat(path)
  if (!stat.isFile() || stat.isSymbolicLink()) {
    throw new Error(`Expected regular file: ${path}`)
  }
  return readFile(path)
}

async function directory(path) {
  const stat = await lstat(path)
  if (!stat.isDirectory() || stat.isSymbolicLink()) {
    throw new Error(`Expected real directory: ${path}`)
  }
}

async function webFiles(source, destination, files) {
  await directory(source)
  for (const name of (await readdir(source)).sort()) {
    const from = join(source, name)
    const to = join(destination, name)
    const stat = await lstat(from)
    if (stat.isSymbolicLink()) throw new Error(`Symlink forbidden: ${from}`)
    if (stat.isDirectory()) {
      await webFiles(from, to, files)
    } else if (stat.isFile()) {
      if (name.toLowerCase().endsWith('.wasm')) {
        throw new Error('Server WASM must not be in web assets')
      }
      files.push([to, await regular(from)])
    } else {
      throw new Error(`Unsupported web asset: ${from}`)
    }
  }
}

// Produces a NEW overlay; never modifies the host checkout itself.
export async function prepareHost(source, output) {
  source = await realpath(source)
  output = resolve(output)
  const manifestBytes = await regular(join(source, 'host-manifest.json'))
  const manifest = JSON.parse(manifestBytes.toString('utf8'))
  if (manifest.schemaVersion !== 1 || !Array.isArray(manifest.games) ||
      manifest.games.length === 0) {
    throw new Error('Non-empty version-1 host manifest required')
  }

  await directory(join(source, 'games'))
  await directory(join(source, 'web'))
  await directory(join(source, 'web/games'))
  const ids = new Set()
  const files = []

  for (const game of manifest.games) {
    if (typeof game.id !== 'string' ||
        !/^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(game.id) || ids.has(game.id) ||
        game.abi !== 'angel-apply-v1' ||
        game.wasm !== `games/${game.id}/game.wasm` ||
        game.entry !== `web/games/${game.id}/index.html`) {
      throw new Error('Invalid or duplicate manifest game')
    }
    ids.add(game.id)
    await directory(join(source, 'games', game.id))
    const wasmPath = join(source, game.wasm)
    if ((await lstat(wasmPath)).size > 16 * 1024 * 1024) {
      throw new Error('WASM exceeds host module size limit')
    }
    const wasm = await regular(wasmPath)
    if (!wasm.subarray(0, 8).equals(header) || sha256(wasm) !== game.wasmSha256) {
      throw new Error(`WASM integrity mismatch: ${game.id}`)
    }
    files.push([game.wasm, wasm])
    const webRoot = join(source, 'web/games', game.id)
    await directory(webRoot)
    const html = await regular(join(webRoot, 'index.html'))
    if (!html.length) throw new Error('Empty web entry')
    await webFiles(webRoot, `internal/web/static/games/${game.id}`, files)
  }

  files.push(['internal/web/static/host-manifest.json', manifestBytes])
  // Validate every source before creating any output. Existing output is rejected.
  await mkdir(dirname(output), { recursive: true })
  await mkdir(output)
  for (const [name, bytes] of files) {
    const target = join(output, name)
    await mkdir(dirname(target), { recursive: true })
    await writeFile(target, bytes, { flag: 'wx' })
  }
  return [...ids].sort()
}

if (process.argv[1] && resolve(process.argv[1]) === fileURLToPath(import.meta.url)) {
  if (process.argv.length !== 4) {
    console.error('Usage: node scripts/prepare-host.mjs SOURCE NEW_OUTPUT')
    process.exitCode = 1
  } else {
    prepareHost(process.argv[2], process.argv[3]).then(ids => {
      console.log(`Prepared games: ${ids.join(', ')}`)
    }).catch(error => {
      console.error(error)
      process.exitCode = 1
    })
  }
}
