import { Injectable, computed, signal } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class AuthStateService {
  private readonly _isLoggedIn = signal(false);
  private _userId = signal<string>('');

  readonly isLoggedIn = computed(() => this._isLoggedIn());
  readonly userId = this._userId.asReadonly();

  constructor() {
    const userId = localStorage.getItem("userId");
    if (userId) {
      this._userId.update(_ => userId);
    }
  }

  setLoggedIn(value: boolean): void {
    this._isLoggedIn.set(value);
  }

  setUserId(id: string) {
    this._userId.update(_ => id);
    localStorage.setItem("userId", id);
  }

  logout(): void {
    this._isLoggedIn.set(false);
  }
}
