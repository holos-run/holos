import { ComponentFixture, TestBed } from '@angular/core/testing';

import { VersionButtonComponent } from './version-button.component';

describe('VersionButtonComponent', () => {
  let component: VersionButtonComponent;
  let fixture: ComponentFixture<VersionButtonComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [VersionButtonComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(VersionButtonComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
