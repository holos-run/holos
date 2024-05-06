import { Inject, Injectable } from '@angular/core';
import { Code, ConnectError } from '@connectrpc/connect';
import { BehaviorSubject, Observable, catchError, of, shareReplay, switchMap } from 'rxjs';
import { ObservableClient } from '../../connect/observable-client';
import { OrganizationService as ConnectOrganizationService } from '../gen/holos/v1alpha1/organization_connect';
import { ListCallerOrganizationsResponse, Organization } from '../gen/holos/v1alpha1/organization_pb';
import { UserService } from './user.service';

@Injectable({
  providedIn: 'root'
})
export class OrganizationService {
  private callerOrganizationsTrigger$ = new BehaviorSubject<void>(undefined);
  private callerOrganizations$: Observable<ListCallerOrganizationsResponse>;

  private fetchCallerOrganizations(): Observable<ListCallerOrganizationsResponse> {
    return this.client.listCallerOrganizations({ request: {} }).pipe(
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
              switchMap(() => this.client.createCallerOrganization({ request: {} }))
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
