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
  authError: string = '';
  isLoading: boolean = false;

  private testUser = {
    email: 'test@chess.com',
    password: '123456'
  };

  constructor(private fb: FormBuilder, private router: Router) {
    this.loginForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(6)]]
    });
  }

  onSubmit() {

    if (this.loginForm.invalid) return;

    this.authError = '';
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
}
