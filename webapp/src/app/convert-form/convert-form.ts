import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ConversionService } from '../services/conversion';

@Component({
  selector: 'app-convert-form',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './convert-form.html',
  styleUrl: './convert-form.css'
})
export class ConvertForm {
  markdown = '';
  format = 'html';
  formats = ['html', 'pdf', 'docx'];

  constructor(private conversion: ConversionService) {}

  submit() {
    this.conversion.convert(this.markdown, this.format).subscribe(blob => {
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `output.${this.format}`;
      a.click();
      window.URL.revokeObjectURL(url);
    });
  }
}
