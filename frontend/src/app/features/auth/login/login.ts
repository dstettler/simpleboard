import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './login.html',
  styleUrls: ['./login.css']
})
export class LoginComponent {
  loginForm: FormGroup;
  registerForm: FormGroup;

  authError: string = '';
  successMessage: string = '';
  isLoading: boolean = false;
  isRegisterMode: boolean = false;

  private testUser = {
    email: 'test@chess.com',
    password: '123456'
  };

  constructor(private fb: FormBuilder, private router: Router) {
    this.loginForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(6)]]
    });

    this.registerForm = this.fb.group({
      userId: ['', [Validators.required, Validators.minLength(3)]],
      name: ['', [Validators.required]],
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

    const { email, password } = this.loginForm.value;

    setTimeout(() => {
      if (email === this.testUser.email && password === this.testUser.password) {
        this.router.navigate(['/dashboard']);
      } else {
        this.authError = 'Invalid email or password';
      }

      this.isLoading = false;
    }, 1000);
  }

  onRegister() {
  if (this.registerForm.invalid) {
    this.registerForm.markAllAsTouched();
    return;
  }

  this.authError = '';
  this.successMessage = '';
  this.isLoading = true;

  setTimeout(() => {
    console.log('Registered user:', this.registerForm.value);

    this.isLoading = false;
    this.successMessage = 'Successfully Registered';

    this.registerForm.reset();
  }, 1000);
}
}