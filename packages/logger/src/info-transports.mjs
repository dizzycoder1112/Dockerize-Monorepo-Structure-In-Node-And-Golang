import { once } from 'node:events'
import { dirname } from 'node:path'
import { mkdir } from 'node:fs/promises'
import build from 'pino-abstract-transport'
import SonicBoom from 'sonic-boom'

export default async function (opts) {
  const { mkdir: createDir, sync, destination: dest } = opts
  if (createDir) {
    const directoryPath = dirname(opts.destination)

    // Ensure the directory exists
    await mkdir(directoryPath, { recursive: true })
  }

  // SonicBoom is necessary to avoid loops with the main thread.
  // It is the same of pino.destination().

  const destination = new SonicBoom({ dest: dest || 1, sync })
  await once(destination, 'ready')

  return build(async (source) => {
    for await (const obj of source) {
      if (obj.level === 30) {
        const toDrain = !destination.write(`${JSON.stringify(obj)}\n`)
        // This block will handle backpressure
        if (toDrain) {
          await once(destination, 'drain')
        }
      }
    }
  }, {
    async close() {
      destination.end()
      await once(destination, 'close')
    },
  })
}
