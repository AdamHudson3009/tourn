import { ComponentFixture, TestBed } from '@angular/core/testing';

import { LllConsolidatedComponent } from './lll-consolidated.component';

describe('LllConsolidatedComponent', () => {
  let component: LllConsolidatedComponent;
  let fixture: ComponentFixture<LllConsolidatedComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [LllConsolidatedComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(LllConsolidatedComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
