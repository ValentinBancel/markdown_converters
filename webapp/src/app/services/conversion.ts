import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class ConversionService {
  private api = '/api/convert';

  constructor(private http: HttpClient) {}

  convert(markdown: string, format: string) {
    return this.http.post(`${this.api}`, { markdown, format }, { responseType: 'blob' });
  }
}
