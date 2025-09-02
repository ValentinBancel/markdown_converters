import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class HistoryService {
  private api = '/api/history';

  constructor(private http: HttpClient) {}

  list() {
    return this.http.get<Array<{ id: number; filename: string }>>(this.api);
  }

  download(id: number) {
    return this.http.get(`${this.api}/${id}`, { responseType: 'blob' });
  }
}
