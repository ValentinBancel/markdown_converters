import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { tap } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private api = '/api/auth';

  constructor(private http: HttpClient) {}

  login(email: string, password: string) {
    return this.http
      .post<{ token: string }>(`${this.api}/login`, { email, password })
      .pipe(tap(res => localStorage.setItem('token', res.token)));
  }

  register(email: string, password: string) {
    return this.http
      .post<{ token: string }>(`${this.api}/register`, { email, password })
      .pipe(tap(res => localStorage.setItem('token', res.token)));
  }
}
