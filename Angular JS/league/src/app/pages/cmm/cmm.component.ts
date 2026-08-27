import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { FormsModule, ReactiveFormsModule } from "@angular/forms";

interface DiviInfo {
  did: number;
  divi: string;
  conf_num: number;
  i_num: number;
}

interface CMMInfo {
  conf: string;
  diviInfo: DiviInfo[];
}

interface OrigCMM {
  conf_num: number[];
  i_num: number[];
}

@Component({
  standalone: true,
  selector: 'app-cmm',
  imports: [ CommonModule, HttpClientModule, FormsModule, ReactiveFormsModule, RouterModule
  ],
  templateUrl: './cmm.component.html',
  styleUrl: './cmm.component.css'
})
export class CmmComponent implements OnInit {
  title = 'cmm';
  APIURL = "http://127.0.0.1:8000/";
  leagueId: number = 0;
  leagueName: string = '';
  tournId: number = 0;
  tournName: string = '';
  changeFreeze: boolean = false;
  cmmInfo: CMMInfo[] = [];
  origCMM: OrigCMM[] = [];

  constructor(
      private http: HttpClient,
      private route: ActivatedRoute,
    ) {
      this.route.params.subscribe(params => {
        this.leagueId = params['league_id'];
        this.leagueName = params['league_name'];
        this.tournId = params['tourn_id'];
        this.tournName = params['tourn_name'];
      });
    }

  payload = {
    tourn_id: 0,
    league_id: 0,
    divi_id: [] as number[],
    num1: [] as number[],
    num2: [] as number[]
  }

  payload_lt = {
    tourn_id: 0,
    league_id: 0
  }

  ngOnInit(): void {
    this.get_cmms();
  }

  get_cmms() {
    this.payload_lt.tourn_id = this.tournId;
    this.payload_lt.league_id = 0;
    
    this.http.post(this.APIURL + "cmm/get", this.payload_lt).subscribe((res) => {
      if (Array.isArray(res)) {  // Ensure res is an array before iterating
        res.forEach((cmmA: { conf: string; descriptions: DiviInfo[] }) => {
          this.cmmInfo.push({ conf: cmmA.conf, diviInfo: cmmA.descriptions });
          this.origCMM.push({ conf_num: [], i_num: [] });
          const indx1 = this.origCMM.length-1;
          if (Array.isArray(cmmA.descriptions)) {
            cmmA.descriptions.forEach(descA => {
              this.origCMM[indx1].conf_num.push(descA.conf_num);
              this.origCMM[indx1].i_num.push(descA.i_num);
            });
          }
        });
      } else {
        console.error("Unexpected response format:", res);
      }
    });
  }

  change() {
    this.payload.tourn_id = this.tournId;
    this.payload.league_id = 0;
    this.payload.divi_id = [];
    this.payload.num1 = [];
    this.payload.num2 = [];

    this.cmmInfo.forEach((co) => {
      co.diviInfo.forEach((dv) => {
        this.payload.divi_id.push(dv.did);
        this.payload.num1.push(dv.conf_num);
        this.payload.num2.push(dv.i_num);
      });
    });

    this.http.post(this.APIURL+"cmm/change",this.payload).subscribe((res)=>{
      alert(res);
      this.reset();
    });
  }

  cancel() {
    for (let indx1=0; indx1 < this.cmmInfo.length; indx1++) {
      for (let indx2=0; indx2 < this.cmmInfo[indx1].diviInfo.length; indx2++) {
        this.cmmInfo[indx1].diviInfo[indx2].conf_num = this.origCMM[indx1].conf_num[indx2];
        this.cmmInfo[indx1].diviInfo[indx2].i_num = this.origCMM[indx1].i_num[indx2];
      }
    }
  }

  reset() {
    this.payload.divi_id = [];
    this.payload.num1 = [];
    this.payload.num2 = [];
    this.cmmInfo = [];
    this.origCMM = [];
    this.get_cmms();
  }
}
