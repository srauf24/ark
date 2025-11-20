import { initContract } from "@ts-rest/core";
import { healthContract } from "./health.js";
import { assetContract } from "./asset.js";
import { logContract } from "./log.js";

const c = initContract();

export const apiContract = c.router({
  System: healthContract,
  Assets: assetContract,
  Logs: logContract,
});
