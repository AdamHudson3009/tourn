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
  name: string;
  rec: string;
  g_fields: string[];
  dv_rsn: string[];
  cn_rsn: string[];
}

export interface ConfStnd {
  conf: string;
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
    });
  }

  payload = {
    id: 0,
    rnd: 0
  }

  ngOnInit(){
    this.get_schres();
  }

  get_schres() {
    this.payload.id = this.tournId;
    this.payload.rnd = Number(this.rnd);
    this.data = { confs: [], conf_ooo: [], conf_stnd: [], conf_gms: [] }; // Clear data

    this.http.post<SchRes>(this.APIURL + "trns/SchResByRnd", this.payload).subscribe(res => {
      this.data = res;
    });

  }

  reload(rndSelect: number, resStndOoo: string) {
    this.resStndOoo = resStndOoo;
    if (this.rnd != rndSelect) { 
      this.rnd = rndSelect;
      this.get_schres();
    }
    
  }
}
