import { Injectable, Inject } from '@angular/core';
import { UserService } from './gen/holos/v1alpha1/user_connect';
import { ObservableClient } from '../connect/observable-client';

@Injectable({
  providedIn: 'root'
})
export class HolosUserService {
  getUser() {
    // returns an Observable
    return this.client.getCallerClaims({ request: {} })
  }

  constructor(
    @Inject(UserService)
    private client: ObservableClient<typeof UserService>
  ) { }
}
