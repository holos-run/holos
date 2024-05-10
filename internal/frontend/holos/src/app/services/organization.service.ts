import { Inject, Injectable } from '@angular/core';
import { Code, ConnectError } from '@connectrpc/connect';
import { BehaviorSubject, Observable, catchError, of, shareReplay, switchMap } from 'rxjs';
import { ObservableClient } from '../../connect/observable-client';
import { OrganizationService as ConnectOrganizationService } from '../gen/holos/organization/v1alpha1/organization_service_connect';
import { Organization } from '../gen/holos/organization/v1alpha1/organization_pb';
import { UserService } from './user.service';

@Injectable({
  providedIn: 'root'
})
export class OrganizationService {
  private callerOrganizationsTrigger$ = new BehaviorSubject<void>(undefined);
  private callerOrganizations$: Observable<Organization[]>;

  private fetchCallerOrganizations(): Observable<Organization[]> {
    return this.client.listOrganizations({}).pipe(
      switchMap(resp => {
        if (resp && resp.organizations.length > 0) {
          return of(resp.organizations)
        }
        return this.client.createOrganization({}).pipe(
          switchMap(resp => {
            if (resp.organization !== undefined) {
              return of([resp.organization])
            } else {
              return of([])
            }
          })
        )
      }),
      catchError(err => {
        if (err instanceof ConnectError) {
          if (err.code == Code.NotFound) {
            return this.userService.createUser().pipe(
              switchMap(() => this.client.createOrganization({}).pipe(
                switchMap(resp => {
                  if (resp.organization !== undefined) {
                    return of([resp.organization])
                  } else {
                    return of([])
                  }
                })
              ))
            )
          }
        }
        console.error('Error fetching data:', err);
        throw err;
      }),
    )
  }

  activeOrg(): Observable<Organization | undefined> {
    return this.callerOrganizations$.pipe(
      switchMap(orgs => of(orgs.at(0)))
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
