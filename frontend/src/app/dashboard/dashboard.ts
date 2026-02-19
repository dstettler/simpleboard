import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="dashboard">
      <h1>Welcome to the Chess Arena ♟️</h1>
      <p>You are successfully logged in.</p>
    </div>
  `
})
export class DashboardComponent {}
