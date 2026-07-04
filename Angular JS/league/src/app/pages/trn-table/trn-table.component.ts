import { Component, Injector, ViewContainerRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { MatIconModule } from '@angular/material/icon';
import { TbrksComponent } from '../tbrks/tbrks.component';
import { SchBuilderComponent} from "../sch-builder/sch-builder.component"
import { CmmComponent } from '../cmm/cmm.component';
import { LllConsolidatedComponent } from "../lll-consolidated/lll-consolidated.component"
import { GrammarComponent } from '../grammar/grammar.component';

interface Plr {
  pid: number;
  plyrtm: string;
  plyrtm2: string;
  sds: number;
}

interface DivPlr {
  did: number;
  divi: string;
  plyrtms: Plr[];
}

interface ConDivPlr {
  cid: number;
  conf: string;
  divis: DivPlr[];
}

interface OrigP{
  plyrtms: string[];
  plyrtms2: string[];
}

interface OrigD{
  divis: string[];
  origP: OrigP[];
}

interface OrigC{
  confs: string[];
  origD: OrigD[];
}

@Component({
  selector: 'app-trn-table',
  imports: [
    CommonModule, HttpClientModule, FormsModule, ReactiveFormsModule, RouterModule,
    MatIconModule, TbrksComponent, SchBuilderComponent, CmmComponent, LllConsolidatedComponent,
    GrammarComponent
  ],
  templateUrl: './trn-table.component.html',
  styleUrl: './trn-table.component.css'
})
export class TrnTableComponent {
  title = 'Trn Table';
  APIURL = "http://127.0.0.1:8000/";
  leagueId: number = 0;
  leagueName: string = '';
  tournId: number = 0;
  tournName: string = '';
  conDivPlr: ConDivPlr[] = [];
  confAdd: string = "";
  diviAdd: string = "";
  plyrAdd: string = "";
  origC: OrigC = {
    confs: [],
    origD: []
  };
  confEdit: number = -1;
  confDiviEdit: number = -1;
  diviEdit: number = -1;
  confPlyrEdit: number = -1;
  diviPlyrEdit: number = -1;
  plyrEdit: number = -1;
  changeFreeze: boolean = false;
  showRightNum = 0;

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

    payload_t = {
      tourn_id: this.tournId
    }

    payload_cdp = {
      table: "conf",
      column: "tourn_id",
      parent_id: this.tournId,
      id: 0,
      descriptions: [] as string[]
    }

    payload_plyrtm = {
      parent_id: 0,
      id: 0,
      descriptions: [] as string[],
      descriptions2: [] as string[],
      sds: 0,
      up: true
    }
    
    ngOnInit() {
      this.get_cdps()
    }

    edit_area(cn: number, cndv: number, dv: number, cnpl: number = -1, dvpl: number = -1, pl: number = -1) {
      if (!this.changeFreeze) {
        this.confEdit = cn;
        this.confDiviEdit = cndv;
        this.diviEdit = dv;
        this.confPlyrEdit = cnpl;
        this.diviPlyrEdit = dvpl;
        this.plyrEdit = pl;  
      }
    }

    get_cdps() {
      this.payload_t.tourn_id = this.tournId;
      this.http.post(this.APIURL + "trn-table/get", this.payload_t).subscribe((res) => {
       if (Array.isArray(res)) {  // Ensure res is an array before iterating
        res.forEach((cdp: { cid: number; conf: string; divis: DivPlr[] }) => {
            this.conDivPlr.push({ cid: cdp.cid, conf: cdp.conf, divis: cdp.divis });
            this.origC.confs.push(cdp.conf);
            const origDEntry: OrigD = {
              divis: [],
              origP: []
            };
            this.origC.origD.push(origDEntry);
            const cindx = this.conDivPlr.length-1;
            if (Array.isArray(cdp.divis)) {
              cdp.divis.forEach(divA => {
                this.origC.origD[cindx].divis.push(divA.divi);
                const origPEntry: OrigP = {
                  plyrtms: [],
                  plyrtms2: []
                };
                this.origC.origD[cindx].origP.push(origPEntry);
                let dindx = this.origC.origD[cindx].origP.length-1;
                if (Array.isArray(this.conDivPlr[cindx].divis[dindx].plyrtms)) {
                  divA.plyrtms.forEach(plrt=> {
                    this.origC.origD[cindx].origP[dindx].plyrtms.push(plrt.plyrtm);
                    this.origC.origD[cindx].origP[dindx].plyrtms2.push(plrt.plyrtm2);
                  });
                }
              });
            }
          });
        } else {
          console.error("Unexpected response format:", res);
        }

        //this.autoexpand();*/
      });
    }

    change_cdp(table:string,column:string,confId:number,desc:string) {
      this.payload_cdp.id = confId;
      this.payload_cdp.parent_id = 0;
      this.payload_cdp.table = table;
      this.payload_cdp.column = column;
      this.payload_cdp.descriptions.push(desc);
      this.http.post(this.APIURL+"tcd/change",this.payload_cdp).subscribe((res)=>{
        alert(res);
        this.reset();
      })
    }

    cancel_cdp(table:string,cindx:number,dindx:number=0,pindx:number=0) {
      if (table == "conf") {
        this.conDivPlr[cindx].conf = this.origC.confs[cindx];
      } else if (table == "divi") {
        this.conDivPlr[cindx].divis[dindx].divi = this.origC.origD[cindx].divis[dindx];
      } else {
        this.conDivPlr[cindx].divis[dindx].plyrtms[pindx].plyrtm = this.origC.origD[cindx].origP[dindx].plyrtms[pindx];
        this.conDivPlr[cindx].divis[dindx].plyrtms[pindx].plyrtm2 = this.origC.origD[cindx].origP[dindx].plyrtms2[pindx];
      }
      this.confEdit = -1;
      this.confDiviEdit = -1;
      this.diviEdit = -1;
      this.confPlyrEdit = -1;
      this.diviPlyrEdit = -1;
      this.plyrEdit = -1
      this.changeFreeze = false;
    }
  
    delete_cdp(table:string,column:string,itemid:number,description:string) {
      this.payload_cdp.id = itemid;
      this.payload_cdp.parent_id = 0;
      this.payload_cdp.table = table;
      this.payload_cdp.column = column;
      this.payload_cdp.descriptions.push(description);
      this.http.post(this.APIURL+"tcd/delete",this.payload_cdp).subscribe((res)=>{
        alert(res);
        this.reset();
      })
    }
  
    add_cdp(table:string,column:string,conf:string,parentId:number) {
      this.payload_cdp.id = 0;
      this.payload_cdp.parent_id = parentId;
      this.payload_cdp.table = table;
      this.payload_cdp.column = column;
      this.payload_cdp.descriptions.push(conf);
      alert(conf);
      this.http.post(this.APIURL+"tcd/add",this.payload_cdp).subscribe((res)=>{
        alert(res);
        this.reset();
      })
    }

    cancel_add() {
      this.confEdit = -1;
      this.confDiviEdit = -1;
      this.diviEdit = -1;
      this.confPlyrEdit = -1;
      this.diviPlyrEdit = -1;
      this.plyrEdit = -1;
      this.changeFreeze = false;
      this.plyrAdd = "";
      this.confAdd = "";
      this.diviAdd = "";
    }

    sds_swap(id:number,up:boolean) {
      this.payload_plyrtm.descriptions = [];
      this.payload_plyrtm.descriptions2 = [];
      this.payload_plyrtm.id = id;
      this.payload_plyrtm.parent_id = 0;
      this.payload_plyrtm.descriptions.push("");
      this.payload_plyrtm.descriptions2.push("");
      this.payload_plyrtm.sds = 0;
      this.payload_plyrtm.up = up;
      this.http.post(this.APIURL+"plyr/sds",this.payload_plyrtm).subscribe((res)=>{
        alert(res);
        this.payload_plyrtm.descriptions.pop();
        this.payload_plyrtm.descriptions2.pop();
        this.reset();
      })
    }

    add_plyrtm(parent_id:number,plyrtm:string,plyrtm2:string="") {
      this.payload_plyrtm.descriptions = [];
      this.payload_plyrtm.descriptions2 = [];
      this.payload_plyrtm.id = 0;
      this.payload_plyrtm.parent_id = parent_id;
      this.payload_plyrtm.descriptions.push(plyrtm);
      if (plyrtm2=="") {
        this.payload_plyrtm.descriptions2.push("");
      } else {
        this.payload_plyrtm.descriptions2.push(plyrtm2);
      }
      this.payload_plyrtm.sds = 0;
      this.payload_plyrtm.up = true;
      this.http.post(this.APIURL+"plyr/add",this.payload_plyrtm).subscribe((res)=>{
        alert(res);
        this.payload_plyrtm.descriptions.pop();
        this.payload_plyrtm.descriptions2.pop();
        this.reset();
      })
    }
    
    change_plyrtm(id:number,parent_id:number,plyrtm:string,plyrtm2:string="") {
      this.payload_plyrtm.descriptions = [];
      this.payload_plyrtm.descriptions2 = [];
      this.payload_plyrtm.id = id;
      this.payload_plyrtm.parent_id = parent_id;
      this.payload_plyrtm.descriptions.push(plyrtm);
      if (plyrtm2=="") {
        this.payload_plyrtm.descriptions2.push("");
      } else {
        this.payload_plyrtm.descriptions2.push(plyrtm2);
      }
      this.payload_plyrtm.sds = 0;
      this.payload_plyrtm.up = true;
      this.http.post(this.APIURL+"plyr/change",this.payload_plyrtm).subscribe((res)=>{
        alert(res);
        this.payload_plyrtm.descriptions.pop();
        this.payload_plyrtm.descriptions2.pop();
        this.reset();
      })
    }

    delete_plyrtm(id:number,plyrtm:string,plyrtm2:string="") {
      this.payload_plyrtm.descriptions = [];
      this.payload_plyrtm.descriptions2 = [];
      this.payload_plyrtm.id = id;
      this.payload_plyrtm.parent_id = 0;
      this.payload_plyrtm.descriptions.push(plyrtm);
      if (plyrtm2=="") {
        this.payload_plyrtm.descriptions2.push("");
      } else {
        this.payload_plyrtm.descriptions2.push(plyrtm2);
      }
      this.payload_plyrtm.sds = 0;
      this.payload_plyrtm.up = true;
      this.http.post(this.APIURL+"plyr/delete",this.payload_plyrtm).subscribe((res)=>{
        alert(res);
        this.payload_plyrtm.descriptions.pop();
        this.payload_plyrtm.descriptions2.pop();
        this.reset();
      })
    }

    reset() {
      this.payload_cdp.descriptions = [];
      this.conDivPlr = [];
      this.origC.confs = [];
      this.origC.origD = [];
      this.confAdd = "";
      this.diviAdd = "";
      this.plyrAdd = "";
      this.changeFreeze = false;
      this.confEdit = -1;
      this.confDiviEdit = -1;
      this.diviEdit = -1
      this.confPlyrEdit = -1;
      this.diviPlyrEdit = -1;
      this.plyrEdit = -1
      this.get_cdps();
    }
}