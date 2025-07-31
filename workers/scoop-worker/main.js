// import { Scoop } from "@harvard-lil/scoop";
import Valkey from "iovalkey";
import fs from "fs";

// Global constants
const noTimeout = 0
const requestQueueKey = "queue:requests"

async function main() {
  // Initialization
  const confData = fs.readFileSync("config.json", "utf-8");
  const config = JSON.parse(confData);
  const valkey = new Valkey();

  // Run forever
  while (true) {
    const request = await fetchRequest(valkey, noTimeout)
    console.log("Recieved request ---")
    console.log(request.SeedURL)
    console.log(request.SeedShadowID)
    console.log(request.Status)
  }
}

/**
 * See entities/capture.go for request format.
 * 
 * @param {Valkey} valkey Valkey client
 * @param {number} timeout Number of seconds to block before returning null. Zero means no timeout.
 *
 * @returns {Promise<object | null>} Returns deserialized request object or null if the operation timed out
 */
async function fetchRequest(valkey, timeout) {
  const data = await valkey.blpop(requestQueueKey, timeout);
  if (data === null) {
    return null
  }
  // Data is list where data[0] is key of the list, data[1] is the returned value.
  const request = JSON.parse(data[1])
  // TODO: add some validation
  return request
}

await main();