import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { AuthService } from '../services/auth';

@Component({
  selector: 'app-auth',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './auth.html',
  styleUrl: './auth.css'
})
export class Auth {
  isLogin = true;
  email = '';
  password = '';

  constructor(private auth: AuthService, private router: Router) {}

  toggle() {
    this.isLogin = !this.isLogin;
  }

  submit() {
    const obs = this.isLogin
      ? this.auth.login(this.email, this.password)
      : this.auth.register(this.email, this.password);
    obs.subscribe(() => this.router.navigate(['/convert']));
  }
}
