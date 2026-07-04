import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TbrksComponent } from './tbrks.component';

describe('TbrksComponent', () => {
  let component: TbrksComponent;
  let fixture: ComponentFixture<TbrksComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TbrksComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(TbrksComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
