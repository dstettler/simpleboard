import { Injectable, computed, signal } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class AuthStateService {
  private readonly _isLoggedIn = signal(false);
  private _userId = signal<string>('');
  private _isGuest = false;

  readonly isLoggedIn = computed(() => this._isLoggedIn());
  readonly userId = this._userId.asReadonly();

  constructor() {
    const userId = localStorage.getItem("userId");
    const guestId = localStorage.getItem("guestId");
    if (userId) {
      this._userId.update(_ => userId);
      this._isLoggedIn.update(_ => true);
    } else if (guestId) {
      this._userId.update(_ => guestId);
      this._isGuest = true;
    }
  }

  isGuest(): boolean {
    return this._isGuest;
  }

  setGuest(b: boolean): void {
    this._isGuest = b;
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
