import { InjectionToken } from "@angular/core";
import type { Transport } from "@connectrpc/connect";

export const TRANSPORT = new InjectionToken<Transport>("connect.transport");
