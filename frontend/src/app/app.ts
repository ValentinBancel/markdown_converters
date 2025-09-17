import { Component, signal } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { MarkdownService, MarkdownRequest, MarkdownResponse } from './services/markdown';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, FormsModule, CommonModule],
  templateUrl: './app.html',
  styleUrl: './app.css'
})
export class App {
  protected readonly title = signal('Markdown Converters');
  
  markdownContent = '';
  selectedFormat = 'html';
  convertedContent = '';
  isLoading = false;
  error = '';
  availableFormats: string[] = [];
  apiStatus = '';

  constructor(private markdownService: MarkdownService) {
    this.checkApiHealth();
    this.loadAvailableFormats();
  }

  checkApiHealth() {
    this.markdownService.checkHealth().subscribe({
      next: (response) => {
        this.apiStatus = response.status;
        console.log('API Health:', response);
      },
      error: (error) => {
        this.apiStatus = 'offline';
        console.error('API Health Check Failed:', error);
      }
    });
  }

  loadAvailableFormats() {
    this.markdownService.getAvailableFormats().subscribe({
      next: (response) => {
        this.availableFormats = response.formats;
        console.log('Available formats:', response);
      },
      error: (error) => {
        console.error('Failed to load formats:', error);
        this.availableFormats = ['html', 'txt']; // fallback
      }
    });
  }

  convertMarkdown() {
    if (!this.markdownContent.trim()) {
      this.error = 'Please enter some markdown content';
      return;
    }

    this.isLoading = true;
    this.error = '';
    this.convertedContent = '';

    const request: MarkdownRequest = {
      content: this.markdownContent,
      format: this.selectedFormat
    };

    this.markdownService.convertMarkdown(request).subscribe({
      next: (response: MarkdownResponse) => {
        this.isLoading = false;
        if (response.success) {
          this.convertedContent = response.convertedContent;
        } else {
          this.error = response.message;
        }
      },
      error: (error) => {
        this.isLoading = false;
        this.error = 'Failed to convert markdown. Make sure the API server is running.';
        console.error('Conversion error:', error);
      }
    });
  }

  clearContent() {
    this.markdownContent = '';
    this.convertedContent = '';
    this.error = '';
  }
}
