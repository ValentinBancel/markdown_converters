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
  downloadUrl = ''; // For PDF download link
  pdfData = ''; // Base64 PDF data

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
          
          // Handle PDF file data
          if (response.format === 'pdf' && response.fileData) {
            this.pdfData = response.fileData;
            this.createDownloadUrl(response.fileData, 'document.pdf', 'application/pdf');
          } else {
            this.pdfData = '';
            this.downloadUrl = '';
          }
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
    this.pdfData = '';
    this.downloadUrl = '';
  }

  // Create download URL for binary files
  createDownloadUrl(base64Data: string, filename: string, mimeType: string) {
    try {
      const binaryString = atob(base64Data);
      const bytes = new Uint8Array(binaryString.length);
      for (let i = 0; i < binaryString.length; i++) {
        bytes[i] = binaryString.charCodeAt(i);
      }
      const blob = new Blob([bytes], { type: mimeType });
      this.downloadUrl = URL.createObjectURL(blob);
    } catch (error) {
      console.error('Error creating download URL:', error);
      this.error = 'Error preparing PDF download';
    }
  }

  // Trigger download
  downloadPdf() {
    if (this.downloadUrl) {
      const link = document.createElement('a');
      link.href = this.downloadUrl;
      link.download = 'document.pdf';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    }
  }
}
