import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PotwComponent } from './potw.component';

describe('PotwComponent', () => {
  let component: PotwComponent;
  let fixture: ComponentFixture<PotwComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PotwComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(PotwComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
