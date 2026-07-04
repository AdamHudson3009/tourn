import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SchBuilderComponent } from './sch-builder.component';

describe('SchBuilderComponent', () => {
  let component: SchBuilderComponent;
  let fixture: ComponentFixture<SchBuilderComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SchBuilderComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SchBuilderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
