import { Component } from '@angular/core';
import { AuthService } from './auth.service';

@Component({
  selector: 'app-login',
  template: `
    <form (ngSubmit)="submit()">
      <input [(ngModel)]="email" name="email" placeholder="Email" />
      <input [(ngModel)]="password" name="password" type="password" placeholder="Password" />
      <button type="submit">Login</button>
    </form>
  `
})
export class LoginComponent {
  email = '';
  password = '';

  constructor(private auth: AuthService) {}

  submit() {
    this.auth.login(this.email, this.password).subscribe(res => {
      this.auth.saveToken(res.token);
    });
  }
}
