import { ZAsset } from "@ark/zod";
import { z } from "zod";

type MyAsset = z.infer<typeof ZAsset>;
const a: MyAsset = {} as any;
const b: string = a.id; // Should be string

