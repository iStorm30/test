import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Cluster {
  Centroid: number[];
  Sum: number[];
  Count: number;
}

@Injectable({
  providedIn: 'root'
})
export class KmeansService {
  constructor(private http: HttpClient) { }

  startProcess(): Observable<Cluster[]> {
    return this.http.get<Cluster[]>('http://localhost:8081/start');
  }
}
