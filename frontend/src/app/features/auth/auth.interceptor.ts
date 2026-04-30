import { HttpClient, HttpInterceptorFn } from '@angular/common/http';

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
  const http = inject(HttpClient);
  const authState = inject(AuthStateService);

  const token = localStorage.getItem('token');
  const guestToken = localStorage.getItem('guestToken');
  let bearer = '';

  if (!token && !guestToken) {
    // Get guest token and return with that
    http.get<GuestResponse>(`${API_ENDPOINT}/api/guest`).subscribe({
      next: (response) => {
        localStorage.setItem("guestToken", response.token);
        authState.setUserId(response.user.guest_id);
        authState.setGuest(true);
        bearer = response.token;
      },
      error: (err) => {
        console.log("Error registering guest token");
      }
    });
  } else if (!token && guestToken) {
    bearer = guestToken;
  } else if (token) {
    bearer = token;
  }

  const authReq = req.clone({
    setHeaders: { Authorization: `Bearer ${bearer}` }
  });

  return next(authReq);
};
