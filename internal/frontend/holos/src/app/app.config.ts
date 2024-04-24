import { ApplicationConfig, importProvidersFrom } from '@angular/core';
import { provideRouter } from '@angular/router';
// import { provideHttpClient, withFetch } from '@angular/common/http';

import { routes } from './app.routes';
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';
import { UserService } from './gen/holos/v1alpha1/user_connect';
import { provideClient } from "../connect/client.provider";
import { ConnectModule } from '../connect/connect.module';

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),
    provideAnimationsAsync(),
    // provideHttpClient(withFetch()),
    provideClient(UserService),
    importProvidersFrom(
      ConnectModule.forRoot({
        baseUrl: "https://jeff.app.dev.k2.holos.run"
      }),
    ),
  ]
};
