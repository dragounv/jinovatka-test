import { Scoop } from "@harvard-lil/scoop";
import Valkey from "iovalkey";
import fs from "fs/promises";
import "process";
import path from "path";

// Global constants
const requestQueueKey = "queue:requests";
const resultQueueKey = "queue:results";

async function main() {
  // Prepare config
  const confData = await fs.readFile("config.json", "utf-8");
  /** @type { WorkerConfig } */
  const config = JSON.parse(confData);

  // Prepare settings for capture
  const captureSettings = prepareCaptureSettings(config);

  // Normalize command line arguments
  const args = process.argv.slice(2);
  if (args.length < 1) {
    console.log("Not enough aruments!");
    printUsage();
    process.exitCode = 1;
    return; // Return out of main instead of calling exit explicitly.
  }

  // Call the command
  switch (args[0]) {
    case "run": {
      await run(captureSettings, config); // This call might never return.
      break;
    }
    case "print-settings": {
      printScoopSettings(captureSettings);
      break;
    }
    default: {
      console.log("Unkonwn command!");
      printUsage();
      process.exitCode = 1;
      break; // Return out of main instead of calling exit explicitly.
    }
  }
}

// prettier-ignore
function printUsage() {
  console.log("First argument must be one of these commands:");
  console.log("  run            - start the worker and listen for requests");
  console.log("  print-settings - print the capture settings that will be used with current configuration");
}

/**
 * Uses default scoop settings and overrides them with values set in "config.captureSettings".
 * List of available options: https://github.com/harvard-lil/scoop/blob/main/options.types.js
 *
 * @param { WorkerConfig } config Configuraton object
 *
 * @returns { ScoopOptions }
 */
function prepareCaptureSettings(config) {
  const settings = Scoop.defaults;
  if (!config.captureSettings) {
    return settings;
  }

  Object.assign(settings, config.captureSettings);

  return settings;
}

/**
 * This function logs capture settings to console and exits.
 * Used for debuging.
 *
 * @param { ScoopOptions } settings Capture settings
 */
function printScoopSettings(settings) {
  console.log(settings);
}

/**
 * This function starts listening for requests in the valkey database.
 * When request is recieved it will capture the requested page.
 *
 * @param {ScoopOptions} captureSettings
 * @param {WorkerConfig} config
 */
async function run(captureSettings, config) {
  // Initialize valkey client
  // TODO: Pass config to Valkey
  const valkey = new Valkey();

  // Run forever and handle requests
  while (true) {
    /** @type { CaptureResult } */
    const result = {
      seedShadowID: "",
      done: false,
      errorMessages: [],
    };

    let request;
    try {
      request = await fetchRequest(valkey);
    } catch (err) {
      console.error("Fetch error: " + err.message);
      console.log(request);
      continue; // Don't fail. continue to another request.
    }

    console.log(request);
    result.seedShadowID = request.seedShadowID;

    try {
      await captureRequest(request, captureSettings, config);
    } catch (err) {
      console.error("Capture error: " + err.message);
      console.log(request);
      continue;
    }

    result.done = true;

    await enqueueResult(valkey, result);
  }
}

/**
 * See entities/capture.go for request format.
 *
 * @param { Valkey } valkey Valkey client
 *
 * @returns { Promise<CaptureRequest> } Returns deserialized request object
 */
async function fetchRequest(valkey) {
  const data = await valkey.blpop(requestQueueKey, 0);
  if (data === null) {
    throw new Error("Valkey operation timed out. This should never happen.");
  }
  // Data is list where data[0] is key of the list, data[1] is the returned value.
  /** @type {CaptureRequest} */
  const request = JSON.parse(data[1]);
  // TODO: add some validation
  return request;
}

/**
 * Run scoop capture.
 * @param { CaptureRequest } request
 * @param { ScoopOptions } captureSettings
 * @param { WorkerConfig } config
 */
async function captureRequest(request, captureSettings, config) {
  const capture = await Scoop.capture(request.seedURL, captureSettings);
  if (capture.state === Scoop.states.FAILED) {
    throw new Error("Capture failed. The URL may not exist.");
  }
  // @ts-ignore Typescript type checker is very unhappy about this. The definition and jsdoc annotation for this function needs some love.
  const wacz = await capture.toWACZ(false);
  const filename = request.seedShadowID + ".wacz";
  await fs.writeFile(path.join(config.outputDir, filename), Buffer.from(wacz));
}

/**
 *
 * @param { Valkey } valkey
 * @param { CaptureResult } result
 */
async function enqueueResult(valkey, result) {
  const data = JSON.stringify(result);
  await valkey.rpush(resultQueueKey, data);
}

// --- Type definitions ---
/**
 * @typedef { object } CaptureRequest
 * @property { string } seedURL
 * @property { string } seedShadowID
 * @property { RequestState } state
 */

/**
 * @typedef {("NewRequest" | "Pending" | "DoneSuccess" | "DoneFailure")} RequestState
 */

/**
 * @typedef { object } WorkerConfig
 * @property { string } outputDir Path to directory where WACZ files will be stored be scoop
 * @property { string } valkeyUrl Adress and port of the valkey database used for request queue
 * @property { object | undefined } captureSettings Overrides for default CaptureOptions used in scoop capture
 */

/**
 * @typedef { object } CaptureResult
 * @property { string } seedShadowID
 * @property {boolean} done
 * @property {string[]} errorMessages
 */

// ------------------------

await main();
