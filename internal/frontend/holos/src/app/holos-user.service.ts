import { Injectable, Inject } from '@angular/core';
import { UserService } from './gen/holos/v1alpha1/user_connect';
import { ObservableClient } from '../connect/observable-client';
import { shareReplay, switchMap } from 'rxjs/operators';
import { Observable, of } from 'rxjs';
import { Claims } from './gen/holos/v1alpha1/user_pb';

@Injectable({
  providedIn: 'root'
})
export class HolosUserService {

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

  constructor(
    @Inject(UserService)
    private client: ObservableClient<typeof UserService>
  ) { }
}
