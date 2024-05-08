import { Inject, Injectable } from '@angular/core';
import { JsonValue, Struct, } from '@bufbuild/protobuf';
import { Observable, of, switchMap } from 'rxjs';
import { ObservableClient } from '../../connect/observable-client';
import { Organization } from '../gen/holos/v1alpha1/organization_pb';
import { PlatformService as ConnectPlatformService } from '../gen/holos/v1alpha1/platform_connect';
import { PlatformServiceGetFormResponse, ListPlatformsRequest, Platform, PlatformServicePutModelRequest, PlatformServicePutModelResponse } from '../gen/holos/v1alpha1/platform_pb';

@Injectable({
  providedIn: 'root'
})
export class PlatformService {
  listPlatforms(org: Observable<Organization>): Observable<Platform[]> {
    return org.pipe(
      switchMap(org => {
        const req = new ListPlatformsRequest({ orgId: org.id })
        return this.client.listPlatforms(req).pipe(
          switchMap(resp => { return of(resp.platforms) })
        )
      })
    )
  }

  getForm(id: string): Observable<PlatformServiceGetFormResponse> {
    return this.client.getForm({ platformId: id })
  }

  putModel(id: string, model: JsonValue): Observable<PlatformServicePutModelResponse> {
    const req = new PlatformServicePutModelRequest({
      platformId: id,
      // "We recommend to use fromJson() to construct Struct literals" refer to
      // https://github.com/bufbuild/protobuf-es/blob/main/docs/runtime_api.md#struct
      model: Struct.fromJson(model),
    })
    return this.client.putModel(req)
  }


  constructor(@Inject(ConnectPlatformService) private client: ObservableClient<typeof ConnectPlatformService>) { }
}
