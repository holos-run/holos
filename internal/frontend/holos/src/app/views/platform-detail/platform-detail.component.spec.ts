import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PlatformDetailComponent } from './platform-detail.component';

describe('PlatformDetailComponent', () => {
  let component: PlatformDetailComponent;
  let fixture: ComponentFixture<PlatformDetailComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PlatformDetailComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(PlatformDetailComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
