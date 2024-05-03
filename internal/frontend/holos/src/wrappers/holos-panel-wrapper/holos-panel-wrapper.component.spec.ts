import { ComponentFixture, TestBed } from '@angular/core/testing';

import { HolosPanelWrapperComponent } from './holos-panel-wrapper.component';

describe('HolosPanelWrapperComponent', () => {
  let component: HolosPanelWrapperComponent;
  let fixture: ComponentFixture<HolosPanelWrapperComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [HolosPanelWrapperComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(HolosPanelWrapperComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
