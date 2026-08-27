import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { MatIconModule } from '@angular/material/icon';

interface XtraInfo {
  id: number;
  ordr: number;
  description: string;
}

interface TournSelect {
  tourn_id: number;
  tourn: string;
}

@Component({
  selector: 'app-lll-consolidated',
  imports: [CommonModule, HttpClientModule, FormsModule, ReactiveFormsModule, RouterModule,
    MatIconModule ],
  templateUrl: './lll-consolidated.component.html',
  styleUrl: './lll-consolidated.component.css'
})
export class LllConsolidatedComponent implements OnInit {
  title = 'sch builder';
  APIURL = "http://127.0.0.1:8000/";
  leagueId: number = 0;
  leagueName: string = '';
  tournId: number = 0;
  tournName: string = '';
  xtraInfo: XtraInfo[] = [];
  changeFreeze: boolean = false;
  xtraEdit: number = -1;
  xtraAdd: string = "";
  xtraDescOrig: string[] = [];
  tournSelect: TournSelect[] = [];
  fromTrnId: number = -1;
  
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
    tourn_id: this.tournId,
    id: 0,
    description: "",
    ordr: 0,
    up: true
  }

  payload_trn = {
    tournId: 0,
    tourn: ""
  }

  payload_trn_copy = {
    TournIdFrom: 0,
    TournIdTo: 0
  }

  ngOnInit(): void {
    this.get_xtra();
  }
  
  edit_area(xtra: number) {
    if (!this.changeFreeze) {
      this.xtraEdit = xtra;
    }
  }

  get_xtra() {
    this.payload.tourn_id = this.tournId;
    this.payload.id = 0;
    this.payload.description = "";
    this.payload.ordr = 0;
    this.payload.up = false;

    this.http.post(this.APIURL + "lllextra/get", this.payload).subscribe((res) => {
      if (Array.isArray(res)) {  // Ensure res is an array before iterating
        res.forEach((r1: { id: number, description: string, ordr: number }) => {
            this.xtraInfo.push({ id: r1.id, description: r1.description, ordr: r1.ordr });
            this.xtraDescOrig.push(r1.description);
          });
        } else {
          console.error("Unexpected response format:", res);
        }
        if (this.xtraInfo.length==0) {
          this.get_tourns_dropdown();
        }
    });
  }

  get_tourns_dropdown() {    
    this.payload_trn.tournId = this.tournId;
    this.payload_trn.tourn = "";
    this.tournSelect = [];

    this.http.post("http://localhost:5164/lllextra/GetCopyFrom", this.payload_trn).subscribe((res) => { 
      if (Array.isArray(res)) {  // Ensure res is an array before iterating
        res.forEach((trn: { tournId: number; tourn: string }) => {
          this.tournSelect.push({ tourn_id: trn.tournId, tourn: trn.tourn });
        });
      }
    });
  }

  copy_from(trnIdFrom:number) {
    this.payload_trn_copy.TournIdFrom = trnIdFrom;
    this.payload_trn_copy.TournIdTo = this.tournId;

    this.http.post("http://localhost:5164/lllextra/CopyFrom",this.payload_trn_copy).subscribe({
      next: (res) => {
        console.log("Success:", res);
        this.reset();
      },
      error: (err) => alert("API Error: " +JSON.stringify(err))
    });
  }

  cancel(indx:number) {
    if (indx<this.xtraInfo.length) {
      this.xtraInfo[indx].description = this.xtraDescOrig[indx];
    }
    this.xtraEdit = -1;
    this.xtraAdd = "";
    this.changeFreeze = false;
  }

  change_xtra(indx:number) {
    this.payload.tourn_id = this.tournId;
    this.payload.id = this.xtraInfo[indx].id
    this.payload.description = this.xtraInfo[indx].description;
    this.payload.ordr = this.xtraInfo[indx].ordr
    this.payload.up = false;

    this.http.post(this.APIURL + "lllextra/change", this.payload).subscribe((res) => {
        alert(res)
        this.reset();
    });
  }

  delete_xtra(indx:number) {
    this.payload.tourn_id = this.tournId;
    this.payload.id = this.xtraInfo[indx].id;
    this.payload.description = this.xtraInfo[indx].description;
    this.payload.ordr = this.xtraInfo[indx].ordr;
    this.payload.up = false;

    this.http.post(this.APIURL + "lllextra/delete", this.payload).subscribe((res) => {
      alert(res)
      this.reset();
    });
  }

  add_xtra() {
    this.payload.tourn_id = this.tournId;
    this.payload.id = 0;
    this.payload.description = this.xtraAdd;
    this.payload.ordr = 0;
    this.payload.up = false;
    
    this.http.post(this.APIURL + "lllextra/add", this.payload).subscribe((res) => {
      alert(res)
      this.reset();
    });
  }

  ordr_swap(indx:number,id:number,up:boolean) {
    this.payload.tourn_id = this.tournId;
    this.payload.id = id;
    this.payload.description = "";  
    this.payload.ordr = 0;
    this.payload.up = up;
    
    this.http.post(this.APIURL+"lllextra/rnk",this.payload).subscribe((res)=>{
      //alert(res);
      this.reset();
    });
  }

  reset() {
    this.xtraInfo = [];
    this.xtraAdd = "";
    this.xtraDescOrig = [];
    this.changeFreeze = false;
    this.xtraEdit = -1;
    this.get_xtra();
  }
}

