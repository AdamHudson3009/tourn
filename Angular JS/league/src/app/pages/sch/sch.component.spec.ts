import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SchComponent } from './sch.component';

describe('SchComponent', () => {
  let component: SchComponent;
  let fixture: ComponentFixture<SchComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SchComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SchComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
