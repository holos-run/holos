import { Inject, Injectable } from '@angular/core';
import { Observable, of, switchMap } from 'rxjs';
import { ObservableClient } from '../../connect/observable-client';
import { Version } from '../gen/holos/system/v1alpha1/system_pb';
import { SystemService as ConnectSystemService } from '../gen/holos/system/v1alpha1/system_service_connect';
import { GetVersionRequest } from '../gen/holos/system/v1alpha1/system_service_pb';
import { FieldMask } from '@bufbuild/protobuf';

@Injectable({
  providedIn: 'root'
})
export class SystemService {
  getVersion(): Observable<Version | undefined> {
    const fieldMask = new FieldMask({ paths: ["version", "git_commit", "go_version", "os", "arch"] })
    const req = new GetVersionRequest({ fieldMask: fieldMask })
    return this.client.getVersion(req).pipe(
      switchMap(resp => { return of(resp.version) })
    )
  }

  constructor(@Inject(ConnectSystemService) private client: ObservableClient<typeof ConnectSystemService>) { }
}
