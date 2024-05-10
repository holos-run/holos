import { Inject, Injectable } from '@angular/core';
import { Observable, switchMap, of, shareReplay } from 'rxjs';
import { ObservableClient } from '../../connect/observable-client';
import { User } from '../gen/holos/user/v1alpha1/user_pb';
import { UserService as ConnectUserService } from '../gen/holos/user/v1alpha1/user_service_connect';

@Injectable({
  providedIn: 'root'
})
export class UserService {

  getUser(): Observable<User | null> {
    return this.client.getUser({}).pipe(
      switchMap(resp => {
        if (resp && resp.user) {
          return of(resp.user)
        } else {
          return this.createUser()
        }
      }),
      // Consolidate to one api call for all subscribers
      shareReplay(1)
    )
  }

  createUser(): Observable<User | null> {
    return this.client.createUser({}).pipe(
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
