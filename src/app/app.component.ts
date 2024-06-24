import { Component } from '@angular/core';
import { KmeansService, Cluster } from './kmeans.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'test';
  response: string;

  constructor(private kmeansService: KmeansService) {
    this.response = '';
    this.clusters = [];
  }

  clusters: Cluster[];

  startProcess(): void {
    this.kmeansService.startProcess().subscribe({
      next: (clusters) => {
        this.clusters = clusters;
      },
      error: (error) => {
        console.error('Error starting process: ' + error);
      }
    });
  }
}
