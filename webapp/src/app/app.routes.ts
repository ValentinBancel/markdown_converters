import { Routes } from '@angular/router';
import { Auth } from './auth/auth';
import { ConvertForm } from './convert-form/convert-form';
import { HistoryList } from './history-list/history-list';

export const routes: Routes = [
  { path: 'login', component: Auth },
  { path: 'convert', component: ConvertForm },
  { path: 'history', component: HistoryList },
  { path: '', pathMatch: 'full', redirectTo: 'login' }
];
