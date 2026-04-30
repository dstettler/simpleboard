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
    console.log('running auth state constructor');
    const userId = localStorage.getItem("userId");
    const guestId = localStorage.getItem("guestId");
    console.log(userId);
    console.log(guestId);
    if (userId) {
      console.log(`setting uid ${userId}`);
      this._userId.update(_ => userId);
      this._isLoggedIn.update(_ => true);
    } else if (guestId) {
      console.log(`setting guid ${guestId}`);
      this._userId.update(_ => guestId);
      this._isGuest = true;
    }
  }

  isGuest(): boolean {
    return this._isGuest;
  }

  setGuest(b: boolean): void {
    console.log(`guest ${b}`)
    this._isGuest = b;
  }

  setLoggedIn(value: boolean): void {
    this._isLoggedIn.set(value);
  }

  setUserId(id: string) {
    console.log(`setting id ${id}`);
    this._userId.update(_ => id);
  }

  logout(): void {
    this._isLoggedIn.set(false);
  }
}
