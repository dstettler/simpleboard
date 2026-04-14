import { Routes } from '@angular/router';
import { LoginComponent } from './features/auth/login/login';
import { HomeComponent } from './features/home/home';
import { DashboardComponent } from './dashboard/dashboard';
import { Game } from './features/game/game';
import { BoardStateService } from './features/game/services/board-state-service';

export const routes: Routes = [
  { path: '', component: HomeComponent },
  { path: 'login', component: LoginComponent },
  { path: 'dashboard', component: DashboardComponent },
  { path: 'game', component: Game, providers: [BoardStateService] },
  { path: '**', redirectTo: '' }
];
