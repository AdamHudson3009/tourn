import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TrnschComponent } from './trnsch.component';

describe('TrnschComponent', () => {
  let component: TrnschComponent;
  let fixture: ComponentFixture<TrnschComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TrnschComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(TrnschComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
