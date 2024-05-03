import { Inject, Injectable } from '@angular/core';
import { PlatformService as ConnectPlatformService } from '../gen/holos/v1alpha1/platform_connect';
import { Observable, filter, of, switchMap } from 'rxjs';
import { ObservableClient } from '../../connect/observable-client';
import { UserDefinedSection, UserDefinedConfig, GetPlatformsRequest, Platform, PutPlatformConfigRequest } from '../gen/holos/v1alpha1/platform_pb';
import { Organization } from '../gen/holos/v1alpha1/organization_pb';
import { Value } from '@bufbuild/protobuf';

export interface Section {
  [field: string]: any;
}

export interface Model {
  [section: string]: Section;
}

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

  putConfig(id: string, model: Model): Observable<Platform> {
    // Round trip the model through JSON to strip multi select field
    // discriminators.  Otherwise, a select field looks like:
    //
    //     {"regions":[{"kind": {"case": "stringValue", "value": "us-east2"}}]}
    //
    // Afterwards, we end up with a plain json representation:
    //
    //     {"regions":["us-east2","us-west2"]}
    const jsonModel: Model = JSON.parse(JSON.stringify(model))
    const values = new UserDefinedConfig
    // Set string values from the model
    Object.keys(jsonModel).forEach(sectionName => {
      values.sections[sectionName] = new UserDefinedSection
      Object.keys(jsonModel[sectionName]).forEach(fieldName => {
        const val = new Value
        val.fromJson(jsonModel[sectionName][fieldName])
        values.sections[sectionName].fields[fieldName] = val
      })
    })

    const req = new PutPlatformConfigRequest({ platformId: id, values: values })
    return this.client.putPlatformConfig(req).pipe(
      switchMap(resp => { return of(resp.platform) }),
      filter((platform): platform is Platform => platform !== undefined),
    )
  }

  constructor(@Inject(ConnectPlatformService) private client: ObservableClient<typeof ConnectPlatformService>) { }
}
