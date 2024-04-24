import { TestBed } from '@angular/core/testing';

import { HolosUserService } from './holos-user.service';

describe('HolosUserService', () => {
  let service: HolosUserService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(HolosUserService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
