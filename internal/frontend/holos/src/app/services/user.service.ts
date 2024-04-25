import { Inject, Injectable } from '@angular/core';
import { Observable, switchMap, of, shareReplay } from 'rxjs';
import { ObservableClient } from '../../connect/observable-client';
import { Claims, User } from '../gen/holos/v1alpha1/user_pb';
import { UserService as ConnectUserService } from '../gen/holos/v1alpha1/user_connect';

@Injectable({
  providedIn: 'root'
})
export class UserService {

  getClaims(): Observable<Claims | null> {
    return this.client.getCallerClaims({ request: {} }).pipe(
      switchMap(getCallerClaimsResponse => {
        if (getCallerClaimsResponse && getCallerClaimsResponse.claims) {
          return of(getCallerClaimsResponse.claims)
        } else {
          return of(null)
        }
      }),
      // Consolidate to one api call for all subscribers
      shareReplay(1)
    )
  }

  createUser(): Observable<User | null> {
    return this.client.createCallerUser({ request: {} }).pipe(
      switchMap(resp => {
        if (resp && resp.user) {
          return of(resp.user)
        } else {
          return of(null)
        }
      }),
      shareReplay(1)
    )
  }

  constructor(
    @Inject(ConnectUserService) private client: ObservableClient<typeof ConnectUserService>,
  ) { }
}
