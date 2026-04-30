import { Injectable, computed, signal } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class AuthStateService {
  private readonly _isLoggedIn = signal(false);
  private _userId = signal<string>('');
  private _guestId = signal<string>('');
  private _isGuest = false;

  readonly isLoggedIn = computed(() => this._isLoggedIn());
  readonly userId = this._userId.asReadonly();
  readonly guestId = this._guestId.asReadonly();

  constructor() {
    const userId = localStorage.getItem('userId');
    const guestId = localStorage.getItem('guestId');
    if (userId) {
      this._userId.update(() => userId);
      this._isLoggedIn.update(() => true);
    } else if (guestId) {
      this._userId.update(() => guestId);
      this._guestId.update(() => guestId);
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
    this._userId.update(() => id);
    localStorage.setItem('userId', id);
  }

  setGuestId(id: string) {
    this._userId.update(() => id);
    this._guestId.update(() => id);
    this._isGuest = true;
    localStorage.setItem('guestId', id);
  }

  logout(): void {
    this._isLoggedIn.set(false);
    this._isGuest = false;
    this._userId.set('');
    this._guestId.set('');
    localStorage.removeItem('userId');
    localStorage.removeItem('guestId');
    localStorage.removeItem('token');
    localStorage.removeItem('guestToken');
  }
}
