import { Component, Injector, ViewContainerRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { MatIconModule } from '@angular/material/icon';
import { PotwComponent } from '../potw/potw.component';
import { MatTableModule } from '@angular/material/table';
import { MatSortModule } from '@angular/material/sort';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatInputModule } from '@angular/material/input';
import { error } from 'console';

export interface SchResPlyrtm {
  rnk: number;
  plyrtm: string;
  columns: string[];
}

export interface SchResDivi {
  divi: string;
  plyrtms: SchResPlyrtm[];
  headers: string[];
}

export interface SchResConf {
  conf: string;
  divis: SchResDivi[];
}

export interface PlyrtmOOO {
  plyrtm: string;
  wn: string[];
  ls: string[];
  tie: string[];
}

export interface DiviOOO {
  divi: string;
  plyrtm_ooo: PlyrtmOOO[];
}

export interface ConfOOO {
  conf: string;
  divi_ooo: DiviOOO[];
}

export interface RnkStnd {
  rnk: number;
  rnkAlt: number;
  plyrId: number;
  name: string;
  rec: string;
  g_fields: string[];
  dv_rsn: string[];
  cn_rsn: string[];
}

export interface ConfStnd {
  conf: string;
  confId: number;
  cn: RnkStnd[];
  wc: RnkStnd[];
  dv: RnkStnd[];
}

export interface PlyrtmGMS {
  plyrtm: string;
  wn: number;
  wns: number[];
  remain: string[];
  w: string[];
  l: string[];
  t: string[];
}

export interface ConfGMS {
  conf: string;
  plyrtm_gms: PlyrtmGMS[];
}

export interface SchRes {
  confs: SchResConf[];
  conf_ooo: ConfOOO[];
  conf_stnd: ConfStnd[];
  conf_gms: ConfGMS[];
}

@Component({
  selector: 'app-sch-res',
  imports: [
    CommonModule, HttpClientModule, FormsModule, ReactiveFormsModule, RouterModule,
    MatIconModule, MatTableModule, MatSortModule, MatPaginatorModule, MatInputModule,
    PotwComponent 
  ],
  templateUrl: './sch-res.component.html',
  styleUrl: './sch-res.component.css'
})
export class SchResComponent {
  title = 'Sch';
  APIURL = "http://localhost:8080/";
  CSharpURL = "http://localhost:5151/";
  leagueId: number = 0;
  leagueName: string = '';
  tournId: number = 0;
  tournName: string = "";
  rnd: number= 0;
  maxRnd: number = 0;
  rndSelect: number = 0;
  rndType: number = 0;
  resStndOoo: string = "";
  potw: boolean = false;
  hasRemain: boolean = false;

  data: SchRes = {
    confs: [],
    conf_ooo: [],
    conf_stnd: [],
    conf_gms: []
  };


  constructor(
    private http: HttpClient,
    private route: ActivatedRoute,
  ) {
    this.route.params.subscribe(params => {
      this.leagueId = params['league_id'];
      this.leagueName = params['league_name'];
      this.tournId = Number(params['tourn_id']);
      this.tournName = params['tourn_name'];
      this.rnd = Number(params['rnd']);
      this.maxRnd = Number(params['maxRnd']);
      this.resStndOoo = params['resStndOoo'];
      this.rndSelect = this.rnd;
      if (this.resStndOoo === 'schres') {
        this.rndType = 2;
      } else if (this.resStndOoo === 'ooo') {
        this.rndType = 4;
      } else if (this.resStndOoo === 'gms') {
        this.rndType = 5;
      }
    });
  }

  payload = {
    id: 0,
    rnd: 0
  }

  payload_wc = {
    trnId: 0,
    confId: [] as number[],
    plyrsConIds: [] as number[][],
  };

  ngOnInit(){
    this.get_schres();
  }

  get_schres() {
    this.payload.id = this.tournId;
    this.payload.rnd = Number(this.rnd);
    this.data = { confs: [], conf_ooo: [], conf_stnd: [], conf_gms: [] }; // Clear data

    this.http.post<SchRes>(this.APIURL + "trns/SchResByRnd", this.payload).subscribe(res => {
      this.data = res;
      this.hasRemain = this.has_remain()
    }, error => {
      alert("Error fetching schedule results: " + JSON.stringify(error));  
    });
  }

  has_remain() {
    for (const one_conf_gms of this.data.conf_gms) {
      for (const one_plyrtm_gms of one_conf_gms.plyrtm_gms) {
        if (one_plyrtm_gms.remain.length > 0) {
          return true;   // here, return exits has_remain itself
        }
      }
    }
    return false;
  }

  reload(rndSelect: number, resStndOoo: string) {
    this.resStndOoo = resStndOoo;
    if (this.rnd != rndSelect) { 
      this.rnd = rndSelect;
      this.get_schres();
    }
  }

  wc_build() {
    this.payload_wc.trnId = this.tournId;
    this.payload_wc.confId = [];
    this.payload_wc.plyrsConIds = [];

    for (const one_conf_stnd of this.data.conf_stnd) {
      this.payload_wc.confId.push(one_conf_stnd.confId);
      const group: number[] = [];
      for (const one_plyrtm_stnd of one_conf_stnd.wc) {
        if (one_plyrtm_stnd.rnkAlt > 0) {
          group.push(one_plyrtm_stnd.plyrId);
        }
      }
      this.payload_wc.plyrsConIds.push(group);
    }

    this.http.post(this.CSharpURL + "trn/wc", this.payload_wc).subscribe(res => {
        alert("WC Build Response: " + JSON.stringify(res));
    });
    alert("WC Build: " + JSON.stringify(this.payload_wc));
  }
}
