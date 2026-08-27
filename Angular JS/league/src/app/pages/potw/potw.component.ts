import { Component, OnInit, ViewChild, ElementRef  } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { MatIconModule } from '@angular/material/icon';

interface PotwInfo {
  rnd: number;
  rnk: string;
  plyr: string;
}

@Component({
  selector: 'app-potw',
  imports: [CommonModule, HttpClientModule, FormsModule, ReactiveFormsModule, RouterModule,
    MatIconModule 
  ],
  templateUrl: './potw.component.html',
  styleUrl: './potw.component.css'
})

export class PotwComponent {
  title = 'potw';
  APIURL = "http://127.0.0.1:8080/";
  tournId: number = 0;
  potwInfo: PotwInfo[] = [];
  maxRnds: number = 0;

  constructor(
      private http: HttpClient,
      private route: ActivatedRoute,
    ) {
      this.route.params.subscribe(params => {
        this.tournId = params['tourn_id'];
      });
    }
  
    payload = {
      id: 0,
    }

    ngOnInit(): void {
      this.get_potw();
    }

    get_potw() {
      this.payload.id = Number(this.tournId);
      this.http.post<PotwInfo>(this.APIURL + "trns/Potw", this.payload).subscribe((res) => {
        if (Array.isArray(res)) {  // Ensure res is an array before iterating
          let plyr = ""
          let rnk = ""
          res.forEach((pInfo: { rnd: number; letters: string; plyr: string; }) => {
            if (this.maxRnds < pInfo.rnd) {
              if (this.maxRnds>0)
                this.potwInfo.push({ rnd: this.maxRnds, rnk: rnk.trim(), plyr: plyr });
              this.maxRnds = pInfo.rnd;
              rnk = "";
            }
            if (pInfo.letters == 'potw')
              plyr = pInfo.plyr;
            else
              rnk += pInfo.plyr + ' ';
          });
          if (this.maxRnds>0)
                this.potwInfo.push({ rnd: this.maxRnds, rnk: rnk.trim(), plyr: plyr });
        } else {
          console.error("Unexpected response format:", res);
        }
      });
    }
}
