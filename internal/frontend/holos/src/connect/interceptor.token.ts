import { InjectionToken } from '@angular/core'
import type { Interceptor } from '@connectrpc/connect'

export const INTERCEPTORS = new InjectionToken<Interceptor[]>(
    'connect.interceptors',
    {
        factory: () => [],
    }
)
