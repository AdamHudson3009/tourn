import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TrnTableComponent } from './trn-table.component';

describe('TrnTableComponent', () => {
  let component: TrnTableComponent;
  let fixture: ComponentFixture<TrnTableComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TrnTableComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(TrnTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
