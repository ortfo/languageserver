import * as YAML from "yaml"
import type { Configuration } from "@ortfo/db/dist/configuration"
import { readFileSync } from "node:fs"

export function loadConfiguration(at: string): Configuration {
  return YAML.parse(readFileSync(at, "utf-8"))
}
