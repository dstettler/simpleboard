import { ComponentFixture, TestBed } from '@angular/core/testing';

import { GameEntryModalComponent } from './game-entry-modal';

describe('GameEntryModalComponent', () => {
  let component: GameEntryModalComponent;
  let fixture: ComponentFixture<GameEntryModalComponent>;

  beforeEach(async () => {
    const localStorage = {
      getItem: (_key: string) => {
        return "1";
      }
    }

    Object.defineProperty(window, 'localStorage', { value:  localStorage });

    await TestBed.configureTestingModule({
      imports: [GameEntryModalComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(GameEntryModalComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
