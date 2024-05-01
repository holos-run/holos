import { Inject, Injectable } from '@angular/core';
import { PlatformService as ConnectPlatformService } from '../gen/holos/v1alpha1/platform_connect';
import { Observable, filter, of, switchMap } from 'rxjs';
import { ObservableClient } from '../../connect/observable-client';
import { ConfigValues, GetPlatformsRequest, Platform, PutPlatformConfigRequest } from '../gen/holos/v1alpha1/platform_pb';
import { Organization } from '../gen/holos/v1alpha1/organization_pb';

@Injectable({
  providedIn: 'root'
})
export class PlatformService {
  getPlatforms(org: Observable<Organization>): Observable<Platform[]> {
    return org.pipe(
      switchMap(org => {
        const req = new GetPlatformsRequest({ orgId: org.id })
        return this.client.getPlatforms(req).pipe(
          switchMap(resp => { return of(resp.platforms) })
        )
      })
    )
  }

  getPlatform(id: string): Observable<Platform> {
    return this.client.getPlatform({ platformId: id }).pipe(
      switchMap(resp => {
        return of(resp.platform);
      }),
      filter((platform): platform is Platform => platform !== undefined),
    )
  }

  putConfig(id: string, values: ConfigValues): Observable<Platform> {
    const req = new PutPlatformConfigRequest({ platformId: id, values: values })
    return this.client.putPlatformConfig(req).pipe(
      switchMap(resp => { return of(resp.platform) }),
      filter((platform): platform is Platform => platform !== undefined),
    )
  }

  constructor(@Inject(ConnectPlatformService) private client: ObservableClient<typeof ConnectPlatformService>) { }
}
