import { Component, inject, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Board } from './components/board/board';

@Component({
  selector: 'app-game',
  standalone: true,
  imports: [Board],
  templateUrl: './game.html',
  styleUrl: './game.css'
})
export class Game implements OnInit {
  private route = inject(ActivatedRoute);
  private router = inject(Router);

  gameId: string = "";

  game: any;

  contructor() {

  }

  ngOnInit(): void {
    const paramId = this.route.snapshot.paramMap.get('id');
    if (paramId) {
      this.gameId = paramId;
    }
    if (!this.gameId) {
      console.warn('No game id found');
      this.router.navigate(['/']);
      return;
    }
  }
}
