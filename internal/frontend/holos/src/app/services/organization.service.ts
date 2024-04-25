import { Inject, Injectable } from '@angular/core';
import { OrganizationService as ConnectOrganizationService } from '../gen/holos/v1alpha1/organization_connect';
import { ObservableClient } from '../../connect/observable-client';
import { Observable, switchMap, of, shareReplay, catchError, BehaviorSubject } from 'rxjs';
import { GetCallerOrganizationsResponse, Organization } from '../gen/holos/v1alpha1/organization_pb';
import { UserService } from './user.service';
import { Code, ConnectError } from '@connectrpc/connect';

@Injectable({
  providedIn: 'root'
})
export class OrganizationService {
  private callerOrganizationsTrigger$ = new BehaviorSubject<void>(undefined);
  private callerOrganizations$: Observable<GetCallerOrganizationsResponse>;

  private fetchCallerOrganizations(): Observable<GetCallerOrganizationsResponse> {
    return this.client.getCallerOrganizations({ request: {} }).pipe(
      switchMap(resp => {
        if (resp && resp.organizations.length > 0) {
          return of(resp)
        }
        return this.client.createCallerOrganization({ request: {} })
      }),
      catchError(err => {
        if (err instanceof ConnectError) {
          if (err.code == Code.NotFound) {
            return this.userService.createUser().pipe(
              switchMap(user => this.client.createCallerOrganization({ request: {} }))
            )
          }
        }
        console.error('Error fetching data:', err);
        throw err;
      }),
    )
  }

  getOrganizations(): Observable<Organization[]> {
    return this.callerOrganizations$.pipe(
      switchMap(resp => of(resp.organizations))
    )
  }

  activeOrg(): Observable<Organization | undefined> {
    return this.callerOrganizations$.pipe(
      switchMap(resp => of(resp.organizations.at(-1)))
    )
  }

  refreshOrganizations(): void {
    this.callerOrganizationsTrigger$.next()
  }

  constructor(
    @Inject(ConnectOrganizationService) private client: ObservableClient<typeof ConnectOrganizationService>,
    private userService: UserService,
  ) {
    this.callerOrganizations$ = this.callerOrganizationsTrigger$.pipe(
      switchMap(() => this.fetchCallerOrganizations()),
      shareReplay(1),
      catchError(err => {
        console.error('Error fetching data:', err);
        throw err;
      })
    );
  }
}
