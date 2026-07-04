import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SchResComponent } from './sch-res.component';

describe('SchResComponent', () => {
  let component: SchResComponent;
  let fixture: ComponentFixture<SchResComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SchResComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SchResComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
