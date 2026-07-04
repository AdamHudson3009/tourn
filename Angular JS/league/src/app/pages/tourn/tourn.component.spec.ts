import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TournComponent } from './tourn.component';

describe('TournComponent', () => {
  let component: TournComponent;
  let fixture: ComponentFixture<TournComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TournComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(TournComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
