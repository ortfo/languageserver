/* --------------------------------------------------------------------------------------------
 * Copyright (c) Microsoft Corporation. All rights reserved.
 * Licensed under the MIT License. See License.txt in the project root for license information.
 * ------------------------------------------------------------------------------------------ */

import { workspace, ExtensionContext, RelativePattern } from "vscode"

import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  TransportKind,
} from "vscode-languageclient/node"
import { loadConfiguration } from "./ortfoconfig"
import * as path from "node:path"
import untildify from "untildify-commonjs"

let client: LanguageClient

export function activate(context: ExtensionContext) {
  const serverModule = "ortfols"

  const configurationFilepath: string = workspace
    .getConfiguration("ortfo")
    .get("globalConfigPath")

  const configurationHome = path.dirname(configurationFilepath)

  console.log(
    `Running ortfo extension with configuration file at: ${configurationFilepath}`
  )

  const configuration = loadConfiguration(configurationFilepath)
  if (!configuration || !configuration.tags || !configuration.technologies) {
    throw new Error(
      `Could not load configuration from ${configurationFilepath}: ${JSON.stringify(
        configuration
      )} is missing a value for tags.repository or technologies.repository`
    )
  }

  // If the extension is launched in debug mode then the debug server options are used
  // Otherwise the run options are used
  const serverOptions: ServerOptions = {
    run: {
      command: serverModule,
      options: { cwd: configurationHome },
      transport: TransportKind.stdio,
    },
    debug: {
      command: serverModule,
      options: { cwd: configurationHome },
      transport: TransportKind.stdio,
    },
  }

  const watchPattern = new RelativePattern(
    configurationHome,
    `{${relativePathsToRepositories(configurationHome, configuration).join(
      ","
    )}}`
  )

  // Options to control the language client
  const clientOptions: LanguageClientOptions = {
    // Register the server for plain text documents
    documentSelector: [
      { scheme: "file", pattern: "**/description.md" },
      { scheme: "file", language: "yaml" },
    ],
    outputChannelName: "ortfols",
    synchronize: {
      // Notify the server about file changes to '.clientrc files contained in the workspace
      fileEvents: workspace.createFileSystemWatcher(watchPattern),
    },
  }

  // Create the language client and start the client.
  client = new LanguageClient("ortfo", "Ortfo", serverOptions, clientOptions)

  // Start the client. This will also launch the server
  client.start()
}

function relativePathsToRepositories(
  configurationHome: string,
  configuration: ReturnType<typeof loadConfiguration>
) {
  return [
    configuration.tags.repository,
    configuration.technologies.repository,
  ].map((p) => path.relative(configurationHome, untildify(p)))
}

export function deactivate(): Thenable<void> | undefined {
  if (!client) {
    return undefined
  }
  return client.stop()
}
