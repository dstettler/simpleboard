import { Component, inject, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { GameService } from '../../core/services/game.service';
import { Board } from './components/board/board';

@Component({
  selector: 'app-game',
  standalone: true,
  imports: [Board], 
  templateUrl: './game.html',
  styleUrl: './game.css'
})
export class Game implements OnInit {
  private gameService = inject(GameService);
  private router = inject(Router);

  game: any;

  ngOnInit() {
    this.game = this.gameService.getGame();

    
    if (!this.game) {
      console.warn('No game found, redirecting...');
      this.router.navigate(['/']);
      return;
    }

    console.log('Loaded game:', this.game);
  }
}