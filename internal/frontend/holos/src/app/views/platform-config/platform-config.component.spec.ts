import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PlatformConfigComponent } from './platform-config.component';

describe('PlatformConfigComponent', () => {
  let component: PlatformConfigComponent;
  let fixture: ComponentFixture<PlatformConfigComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PlatformConfigComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(PlatformConfigComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
