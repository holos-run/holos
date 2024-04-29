import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PlatformsComponent } from './platforms.component';

describe('PlatformsComponent', () => {
  let component: PlatformsComponent;
  let fixture: ComponentFixture<PlatformsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PlatformsComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(PlatformsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
