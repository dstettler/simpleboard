import { Injectable, computed, signal } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class AuthStateService {
  private readonly _isLoggedIn = signal(false);

  readonly isLoggedIn = computed(() => this._isLoggedIn());

  setLoggedIn(value: boolean): void {
    this._isLoggedIn.set(value);
  }

  logout(): void {
    this._isLoggedIn.set(false);
  }
}