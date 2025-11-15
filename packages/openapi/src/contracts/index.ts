import { initContract } from "@ts-rest/core";
import { healthContract } from "./health.js";
import { plantContract } from "./plant.js";
import { observationContract } from "./observation.js";

const c = initContract();

export const apiContract = c.router({
  Health: healthContract,
  Plant: plantContract,
  Observation: observationContract,
});
