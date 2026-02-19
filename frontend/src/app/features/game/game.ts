import { Component } from '@angular/core';
import { Board } from './components/board/board';

@Component({
  selector: 'app-game',
  imports: [Board],
  templateUrl: './game.html',
  styleUrl: './game.css',
})
export class Game {

}
