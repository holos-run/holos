import { Inject, Injectable } from '@angular/core';
import { OrganizationService as ConnectOrganizationService } from '../gen/holos/v1alpha1/organization_connect';
import { ObservableClient } from '../../connect/observable-client';
import { Observable, switchMap, of, shareReplay, catchError } from 'rxjs';
import { Organization } from '../gen/holos/v1alpha1/organization_pb';
import { UserService } from './user.service';
import { Code, ConnectError } from '@connectrpc/connect';

@Injectable({
  providedIn: 'root'
})
export class OrganizationService {
  getOrganizations(): Observable<Organization[] | null> {
    return this.client.getCallerOrganizations({ request: {} }).pipe(
      switchMap(resp => {
        if (resp && resp.organizations) {
          if (resp.organizations.length > 0) {
            return of(resp.organizations)
          } else {
            return this.addOrganization()
          }
        } else {
          return of(null)
        }
      }),
      catchError(err => {
        if (err instanceof ConnectError) {
          if (err.code == Code.NotFound) {
            return this.userService.createUser().pipe(
              switchMap(user => {
                if (user) {
                  return this.addOrganization()
                } else {
                  return of(null)
                }
              })
            )
          }
        }
        return of(null)
      }),
      shareReplay(1)
    )
  }

  addOrganization(): Observable<Organization[] | null> {
    return this.client.createCallerOrganization({ request: {} }).pipe(
      switchMap(resp => {
        if (resp && resp.organizations) {
          return of(resp.organizations)
        } else {
          return of(null)
        }
      })
    )
  }


  constructor(
    @Inject(ConnectOrganizationService) private client: ObservableClient<typeof ConnectOrganizationService>,
    private userService: UserService,
  ) { }
}
