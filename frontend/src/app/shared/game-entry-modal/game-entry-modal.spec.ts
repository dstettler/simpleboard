import { ComponentFixture, TestBed } from '@angular/core/testing';

import { GameEntryModal } from './game-entry-modal';

describe('GameEntryModal', () => {
  let component: GameEntryModal;
  let fixture: ComponentFixture<GameEntryModal>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [GameEntryModal]
    })
    .compileComponents();

    fixture = TestBed.createComponent(GameEntryModal);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
