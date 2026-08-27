import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { MatIconModule } from '@angular/material/icon';

interface TBDesc {
  id: number;
  description: string;
  rnk: number;
}

interface TbrkInfo {
  dv_cn_other: string;
  tbDesc: TBDesc[];
}

interface OrigDescs {
  descriptions: string[];
}

interface OrigDCO {
  dv_cn_other: string[];
  origDescs: OrigDescs[];
}

interface TournSelect {
  tourn_id: number;
  tourn: string;
}

@Component({
  standalone: true,
  selector: 'app-tbrks',
  imports: [ CommonModule, HttpClientModule, FormsModule, ReactiveFormsModule, RouterModule,
    MatIconModule 
  ],
  templateUrl: './tbrks.component.html',
  styleUrl: './tbrks.component.css'
})
export class TbrksComponent implements OnInit {
  title = 'tbrk';
  APIURL = "http://127.0.0.1:8000/";
  leagueId: number = 0;
  leagueName: string = '';
  tournId: number = 0;
  tournName: string = '';
  tbrks: TbrkInfo[] = [];
  changeFreeze: boolean = false;
  dcoAdd: string = "";
  descAdd: string = "";
  dcoEdit: number = -1;
  descDcoEdit: number = -1;
  descEdit: number = -1;
  origDco: OrigDCO = {
    dv_cn_other: [],
    origDescs: []
  };
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
    tourn_id: 0,
    league_id: 0,
    dvCnOth: "",
    id: 0,
    descriptions: [] as string[],
    rnk: 0,
    up: false
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
    this.get_tbrks();
  }

  edit_area(dco: number, descDco: number, desc: number) {
    if (!this.changeFreeze) {
      this.dcoEdit = dco
      this.descDcoEdit = descDco;
      this.descEdit = desc;
    }
  }

  get_tbrks() {
    this.payload.tourn_id = this.tournId;
    this.payload.league_id = 0;
    this.payload.dvCnOth = "";
    this.payload.id = 0;
    this.payload.descriptions = [];  
    this.payload.rnk = 0;
    this.payload.up = false;

    this.http.post(this.APIURL + "tbrk/get", this.payload).subscribe((res) => {
      if (Array.isArray(res)) {  // Ensure res is an array before iterating
        res.forEach((tb_rk: { dv_cn_other: string; descriptions: TBDesc[] }) => {
          this.tbrks.push({ dv_cn_other: tb_rk.dv_cn_other, tbDesc: tb_rk.descriptions });
          this.origDco.dv_cn_other.push(tb_rk.dv_cn_other);
          const origDEntry: OrigDescs = {
            descriptions: []
          };
          this.origDco.origDescs.push(origDEntry);
          const indx1 = this.origDco.origDescs.length-1;
          if (Array.isArray(tb_rk.descriptions)) {
            tb_rk.descriptions.forEach(descA => {
              this.origDco.origDescs[indx1].descriptions.push(descA.description);
            });
          }
        });
        if (this.tbrks.length==0) {
          this.get_tourns_dropdown();
        }
      } else {
        console.error("Unexpected response format:", res);
      }
    });
  }

  cancel(indx:number,indx2:number=-1) {
    if (indx<this.tbrks.length) {
      this.tbrks[indx].dv_cn_other = this.origDco.dv_cn_other[indx];
      if ((indx2>-1) && (indx2<this.tbrks[indx].tbDesc.length)) {
        this.tbrks[indx].tbDesc[indx2].description = this.origDco.origDescs[indx].descriptions[indx2];
      }
    }
    this.dcoEdit = -1;
    this.descDcoEdit = -1;
    this.descEdit = -1;
    this.changeFreeze = false;
  }

  change_dco(indx:number,dco:string) {
    if (dco == "") {
      alert("Please enter value");
      return;
    }
    this.payload.tourn_id = this.tournId;
    this.payload.league_id = 0;
    this.payload.dvCnOth = this.origDco.dv_cn_other[indx];
    this.payload.id = 0;
    this.payload.descriptions = [];
    this.payload.descriptions.push(dco);  
    this.payload.rnk = 0;
    this.payload.up = false

    this.http.post(this.APIURL + "tbrk/change-dco", this.payload).subscribe((res) => {
        alert(res)
        this.reset();
    });
  }

  get_tourns_dropdown() {    
    this.payload_trn.tournId = this.tournId;
    this.payload_trn.tourn = "";
    this.tournSelect = [];

    this.http.post("http://localhost:5164/tbrk/GetCopyFrom", this.payload_trn).subscribe((res) => { 
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

    this.http.post("http://localhost:5164/tbrk/CopyFrom",this.payload_trn_copy).subscribe({
      next: (res) => {
        console.log("Success:", res);
        this.reset();
      },
      error: (err) => alert("API Error: " +JSON.stringify(err))
    });
  }

  add_dco(dco:string,desc:string) {
    if ((dco == "") || (desc == "")) {
      alert("Please enter values for both");
      return;
    }
    this.payload.tourn_id = this.tournId;
    this.payload.league_id = 0;
    this.payload.dvCnOth = dco;
    this.payload.id = 0;
    this.payload.descriptions = [];
    this.payload.descriptions.push(desc);  
    this.payload.rnk = 0;
    this.payload.up = false

    this.http.post(this.APIURL + "tbrk/add", this.payload).subscribe((res) => {
        alert(res)
        this.reset();
    });
  }

  rnk_swap(indx:number,id:number,up:boolean) {
    this.payload.tourn_id = this.tournId;
    this.payload.league_id = 0;    
    this.payload.dvCnOth = this.origDco.dv_cn_other[indx];
    this.payload.id = id;
    this.payload.descriptions = [];
    this.payload.descriptions.push("");  
    this.payload.rnk = 0;
    this.payload.up = up;
    
    this.http.post(this.APIURL+"tbrk/rnk",this.payload).subscribe((res)=>{
      //alert(res);
      this.reset();
    });
  }

  change_desc(id:number,dvCnOth:string,description:string) {
    if (description == "") {
      alert("Please enter value");
      return;
    }
    this.payload.rnk = 0;
    this.payload.up = false

    this.payload.tourn_id = this.tournId;
    this.payload.league_id = 0;
    this.payload.dvCnOth = dvCnOth;
    this.payload.id = id;
    this.payload.descriptions = [];
    this.payload.descriptions.push(description);  
    this.payload.rnk = 0;
    this.payload.up = false;

    this.http.post(this.APIURL+"tbrk/change",this.payload).subscribe((res)=>{
      alert(res);
      this.reset();
    });
  }

  delete_desc(id:number,cvDnOth:string,description:string) {
    this.payload.tourn_id = this.tournId;
    this.payload.league_id = 0;
    this.payload.dvCnOth = cvDnOth;
    this.payload.id = id;
    this.payload.descriptions = [];
    this.payload.descriptions.push(description);  
    this.payload.rnk = 0;
    this.payload.up = false;
    
    this.http.post(this.APIURL+"tbrk/delete",this.payload).subscribe((res)=>{
      alert(res);
      this.reset();
    });
  }

  add_desc(cvDnOth:string,description:string) {
    if (description == "") {
      alert("Please enter value");
      return;
    }
    this.payload.tourn_id = this.tournId;
    this.payload.league_id = 0;
    this.payload.dvCnOth = cvDnOth;
    this.payload.id = 0;
    this.payload.descriptions = [];
    this.payload.descriptions.push(description);  
    this.payload.rnk = 0;
    this.payload.up = false;

    this.http.post(this.APIURL+"tbrk/add",this.payload).subscribe((res)=>{
      alert(res);
      this.reset();
    });
  }

  reset() {
    this.payload.descriptions = [];
    this.tbrks = [];
    this.origDco.dv_cn_other = [];
    this.origDco.origDescs = [];
    this.dcoAdd = "";
    this.descAdd = "";
    this.changeFreeze = false;
    this.dcoEdit = -1;
    this.descDcoEdit = -1;
    this.dcoEdit = -1;
    this.get_tbrks();
  }
}
