import { Component, inject } from '@angular/core';
import { Router, RouterLink, RouterLinkActive } from '@angular/router';
import { CommonModule } from '@angular/common';
import { GameService } from '../../core/services/game.service';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [CommonModule, RouterLink, RouterLinkActive],
  templateUrl: './navbar.html',
  styleUrl: './navbar.css'
})
export class NavbarComponent {
  private router = inject(Router);
  private gameService = inject(GameService);

  isCreating = false;

  onStartGame() {
    if (this.isCreating) return;

    this.isCreating = true;

    this.gameService.createGame().subscribe({
      next: (res) => {
        this.router.navigate(['/game', res.gameId]);
      },
      error: () => {
        this.isCreating = false;
      }
    });
  }
}