import { Component } from '@angular/core';
import { RouterModule } from '@angular/router'; // ✅ Import RouterModule

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterModule], // ✅ Add RouterModule to imports
  template: `
    <router-outlet></router-outlet>
  `,
  styles: [`
    nav {
      display: flex;
      gap: 10px;
      margin-bottom: 20px;
    }
  `]
})
export class AppComponent { }
