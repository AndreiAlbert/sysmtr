import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { SystemStats } from '../models/stats.model';
import { webSocket, WebSocketSubject } from 'rxjs/webSocket';

@Injectable({
  providedIn: 'root'
})
export class MonitorService {
  private baseUrl = 'http://localhost:8080';
  private wsUrl = 'ws://localhost:8080/ws';

  constructor(private http: HttpClient) { }

  getHistory(): Observable<SystemStats[]> {
    return this.http.get<SystemStats[]>(`${this.baseUrl}/history`);
  }

  getRealTimeStats(): WebSocketSubject<SystemStats> {
    return webSocket(this.wsUrl);
  }
}
