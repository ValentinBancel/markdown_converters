import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HistoryService } from '../services/history';

@Component({
  selector: 'app-history-list',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './history-list.html',
  styleUrl: './history-list.css'
})
export class HistoryList implements OnInit {
  history: Array<{ id: number; filename: string }> = [];

  constructor(private historyService: HistoryService) {}

  ngOnInit() {
    this.historyService.list().subscribe(data => (this.history = data));
  }

  download(item: { id: number; filename: string }) {
    this.historyService.download(item.id).subscribe(blob => {
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = item.filename;
      a.click();
      window.URL.revokeObjectURL(url);
    });
  }
}
