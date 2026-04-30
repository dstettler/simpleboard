import { Component, inject } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { HttpErrorResponse } from '@angular/common/http';
import { AuthService } from '../auth.service';
import { AuthStateService } from '../../../core/services/auth-state.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './login.html',
  styleUrls: ['./login.css']
})
export class LoginComponent {
  private fb = inject(FormBuilder);
  private router = inject(Router);
  private authService = inject(AuthService);
  private authState = inject(AuthStateService);

  loginForm: FormGroup;
  registerForm: FormGroup;

  authError = '';
  successMessage = '';
  isLoading = false;
  isRegisterMode = false;

  constructor() {
    this.loginForm = this.fb.group({
      username: ['', [Validators.required, Validators.minLength(3)]],
      password: ['', [Validators.required, Validators.minLength(6)]]
    });

    this.registerForm = this.fb.group({
      username: ['', [Validators.required, Validators.minLength(3)]],
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(6)]]
    });
  }

  toggleMode() {
    this.isRegisterMode = !this.isRegisterMode;
    this.authError = '';
    this.successMessage = '';
    this.isLoading = false;
  }

  onSubmit() {
    if (this.loginForm.invalid) {
      this.loginForm.markAllAsTouched();
      return;
    }

    this.authError = '';
    this.successMessage = '';
    this.isLoading = true;

    const payload = {
      username: this.loginForm.value.username ?? '',
      password: this.loginForm.value.password ?? ''
    };

    this.authService.login(payload).subscribe({
      next: (response) => {
        this.isLoading = false;
        this.authState.setLoggedIn(true);
        this.authState.setUserId(response.user.user_id.toString());
        const token = response.token;
        localStorage.setItem("token", token);
        localStorage.setItem("userId", response.user.user_id.toString());
        this.successMessage = response.message;
        this.router.navigate(['/dashboard']);
      },
      error: (error: HttpErrorResponse) => {
        this.isLoading = false;
        this.authError = error.error?.message || 'Login failed';
      }
    });
  }

  onRegister() {
    if (this.registerForm.invalid) {
      this.registerForm.markAllAsTouched();
      return;
    }

    this.authError = '';
    this.successMessage = '';
    this.isLoading = true;

    const payload = {
      username: this.registerForm.value.username ?? '',
      email: this.registerForm.value.email ?? '',
      password: this.registerForm.value.password ?? ''
    };

    this.authService.register(payload).subscribe({
      next: (response) => {
        this.isLoading = false;
        this.successMessage = response.message;
        this.isRegisterMode = false;
        this.registerForm.reset();
      },
      error: (error: HttpErrorResponse) => {
        this.isLoading = false;
        this.authError = error.error?.message || 'Registration failed';
      }
    });
  }
}
