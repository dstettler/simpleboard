import { HttpBackend, HttpClient, HttpInterceptorFn } from '@angular/common/http';
import { catchError, map, of, switchMap, tap } from 'rxjs';

import { API_ENDPOINT } from '../../app.constants';
import { inject } from '@angular/core';
import { AuthStateService } from '../../core/services/auth-state.service';
type UserType = {
  guest_id: string;
}

type GuestResponse = {
  message: string;
  token: string;
  user: UserType;
}

export const authInterceptor: HttpInterceptorFn = (req, next) => {
  const http = new HttpClient(inject(HttpBackend));
  const authState = inject(AuthStateService);

  const token = localStorage.getItem('token');
  const guestToken = localStorage.getItem('guestToken');

  const addAuthHeader = (bearer: string) => {
    return req.clone({
      setHeaders: bearer ? { Authorization: `Bearer ${bearer}` } : {}
    });
  };

  const request$ = (() => {
    if (req.url === `${API_ENDPOINT}/api/guest`) {
      return of(req);
    }

    if (token) {
      return of(addAuthHeader(token));
    }

    if (guestToken) {
      return of(addAuthHeader(guestToken));
    }

    return http.get<GuestResponse>(`${API_ENDPOINT}/api/guest`).pipe(
      tap((response) => {
        localStorage.setItem('guestToken', response.token);
        localStorage.setItem('guestId', response.user.guest_id);
        authState.setGuestId(response.user.guest_id);
      }),
      map((response) => addAuthHeader(response.token)),
      catchError((err) => {
        console.error('Error registering guest token', err);
        return of(req);
      })
    );
  })();

  return request$.pipe(switchMap((authReq) => next(authReq)));
};
