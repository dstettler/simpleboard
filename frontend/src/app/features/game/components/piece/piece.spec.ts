import { ComponentFixture, TestBed } from '@angular/core/testing';

import { Piece } from './piece';
import { Pawn } from '../../services/pieces/ChessPiece';

describe('Piece', () => {
  let component: Piece;
  let fixture: ComponentFixture<Piece>;
  const defaultPiece = new Pawn(0, false, 0, 0);

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Piece]
    })
    .compileComponents();

    fixture = TestBed.createComponent(Piece);
    component = fixture.componentInstance;
    
    const mockElement = document.createElement('div');

    component.piece = defaultPiece;
    component.boardRef = mockElement;

    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
