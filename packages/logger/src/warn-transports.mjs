import { once } from 'node:events'
import { mkdir } from 'node:fs/promises'
import { dirname } from 'node:path'
import build from 'pino-abstract-transport'
import SonicBoom from 'sonic-boom'

export default async function (opts) {
  if (opts.mkdir) {
    const directoryPath = dirname(opts.destination)

    // Ensure the directory exists
    await mkdir(directoryPath, { recursive: true })
  }

  // SonicBoom is necessary to avoid loops with the main thread.
  // It is the same of pino.destination().
  const destination = new SonicBoom({ dest: opts.destination || 1, sync: false })
  await once(destination, 'ready')

  return build(async (source) => {
    for await (const obj of source) {
      if (obj.level === 40) {
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
