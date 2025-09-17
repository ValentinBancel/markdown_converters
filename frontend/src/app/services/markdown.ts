import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface MarkdownRequest {
  content: string;
  format: string;
}

export interface MarkdownResponse {
  convertedContent: string;
  format: string;
  success: boolean;
  message: string;
}

export interface HealthResponse {
  message: string;
  status: string;
}

export interface FormatsResponse {
  formats: string[];
  message: string;
}

@Injectable({
  providedIn: 'root'
})
export class MarkdownService {
  private apiUrl = 'http://localhost:8080/api';

  constructor(private http: HttpClient) {}

  checkHealth(): Observable<HealthResponse> {
    return this.http.get<HealthResponse>(`${this.apiUrl}/health`);
  }

  getAvailableFormats(): Observable<FormatsResponse> {
    return this.http.get<FormatsResponse>(`${this.apiUrl}/formats`);
  }

  convertMarkdown(request: MarkdownRequest): Observable<MarkdownResponse> {
    return this.http.post<MarkdownResponse>(`${this.apiUrl}/convert`, request);
  }
}
