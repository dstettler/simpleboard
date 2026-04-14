import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AuthStateService } from './auth-state.service';

export const authGuard: CanActivateFn = () => {
  const authState = inject(AuthStateService);
  const router = inject(Router);

  if (authState.isLoggedIn()) {
    return true;
  }

  return router.createUrlTree(['/login']);
};