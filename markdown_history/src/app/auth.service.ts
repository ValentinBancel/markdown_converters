import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Router } from '@angular/router';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private tokenKey = 'token';
  constructor(private http: HttpClient, private router: Router) {}

  register(email: string, password: string) {
    return this.http.post('/auth/register', { email, password });
  }

  login(email: string, password: string) {
    return this.http.post<{ token: string }>('/auth/login', { email, password });
  }

  saveToken(token: string) {
    localStorage.setItem(this.tokenKey, token);
  }

  getToken(): string | null {
    return localStorage.getItem(this.tokenKey);
  }

  logout() {
    localStorage.removeItem(this.tokenKey);
    this.router.navigate(['/login']);
  }

  authHeaders(): HttpHeaders {
    const token = this.getToken();
    return new HttpHeaders(token ? { Authorization: `Bearer ${token}` } : {});
  }

  isLoggedIn(): boolean {
    return !!this.getToken();
  }
}
