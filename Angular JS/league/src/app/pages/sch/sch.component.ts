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
import { MatTableModule } from '@angular/material/table';
import { MatSortModule } from '@angular/material/sort';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatInputModule } from '@angular/material/input';

interface Mtch {
	gid: number;
	plyr1: number;
	plyr2: number;
	ordr: number;
	plyrwn: number;
	pf: number;
	pa: number;
	plyr1nm: string;
	plyr2nm: string;
	plyrwnnm: string;
	notes: string;
  rec1: number;
  tie1: number;
  rec2: number;
  tie2: number;
  change: boolean;
  msc1: string[];
  msc2: string[];
  letter1: string;
  letter2: string;
  arch_nm: string;
  arch_pts: string;
}

interface PlfSelect {
  indx: number;
  name: string;
}

@Component({
  selector: 'app-sch',
  imports: [
    CommonModule, HttpClientModule, FormsModule, ReactiveFormsModule, RouterModule,
    MatIconModule, TbrksComponent, SchBuilderComponent, CmmComponent, LllConsolidatedComponent,
    GrammarComponent, MatTableModule, MatSortModule, MatPaginatorModule, MatInputModule
  ],
  templateUrl: './sch.component.html',
  styleUrl: './sch.component.css'
})
export class SchComponent {
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
  mtchs: Mtch[] = [];
  Oldmtchs: Mtch[] = [];
  displayedColumns: string[] = ['plyr1nm', 'plyr2nm', 'plyrwnnm', 'actions'];
  mscRow: number = -1;
  mscEdit: number = -1;
  mscAdd1: number = -1;
  mscAdd2: string = "";
  changeFreeze: boolean = false;
  
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
    });
  }

  payload = {
    id: 0,
    rnd: 0
  }

  payload_swp = {
    tourn_id: 0,
    rnd: 0, 
    gid: 0, 
    up: false  
  }

  payload_edit = {
    gid: 0, 
    plyrwn: 0,
    pf: 0,
    pa: 0,
    notes: ""  
  }

  payload_msc = {
    tourn_id: 0,
    rnd: 0, 
    plyrtm_id: 0,
    letters: "",
    add_delete: 0
  }

  ngOnInit(){
    this.get_sch();
  }

  get_sch() {
    this.payload.id = this.tournId;
    this.payload.rnd = Number(this.rnd);
  
    this.http.post<Mtch[]>(this.APIURL + "trns/GetSch", this.payload).subscribe({
      next: (res) => {
        this.mtchs = res;
        this.Oldmtchs = res.map(m => JSON.parse(JSON.stringify(m)));
      },
      error: (err) => {
        alert("HTTP ERROR: " + JSON.stringify(err));
        console.error('HTTP POST failed:', err);
      }
    });
  }
  
  swap_row(gid:number, up:boolean) {
    this.payload_swp.tourn_id = this.tournId;
    this.payload_swp.rnd = Number(this.rnd);
    this.payload_swp.gid = gid
    this.payload_swp.up = up

    this.http.post<any>(this.APIURL + "trns/SwapOrdr", this.payload_swp).subscribe({
      next: (res) => {
        this.mtchs = [];
        this.Oldmtchs = [];
        this.get_sch();
      },
      error: (err) => {
        alert("HTTP ERROR: " + JSON.stringify(err));
        console.error('HTTP POST failed:', err);
      }
    });
  }

  edit_row(i: number, gid:number, plyrwn:number, plyrwnnm: string, pf: number, pa: number, notes: string) {
    this.payload_edit.gid = gid;
    this.payload_edit.plyrwn = Number(plyrwn);
    this.payload_edit.pf = pf;
    this.payload_edit.pa = pa;
    this.payload_edit.notes = notes

    this.http.post<any>(this.APIURL + "trns/EditSch", this.payload_edit).subscribe({
      next: (res) => {
        alert(res.msg);
        if (res.msg === "ok") {
          this.Oldmtchs[i].plyrwn = plyrwn;
          this.Oldmtchs[i].plyrwnnm = plyrwnnm;
          this.Oldmtchs[i].pf = pf;
          this.Oldmtchs[i].pa = pa;
          this.Oldmtchs[i].notes = notes; 
          this.mtchs[i].change = false;
        }
      },
      error: (err) => {
        alert("HTTP ERROR: " + JSON.stringify(err));
        console.error('HTTP POST failed:', err);
      }
    });
  }

  add_msc(indx: number, plyrtm: number, letters: string) {
    this.payload_msc.tourn_id =Number(this.tournId);
    this.payload_msc.rnd = Number(this.rnd);
    this.payload_msc.plyrtm_id = Number(plyrtm); 
    this.payload_msc.letters = letters;
    this.payload_msc.add_delete = 1;

    this.http.post<any>(this.APIURL + "trns/AddDeleteMsc", this.payload_msc).subscribe({
      next: (res) => {
        alert(res.msg);
        this.changeFreeze = false;
        if (plyrtm == this.mtchs[indx].plyr1) {
          this.mtchs[indx].msc1.push(letters);
        } else {
          this.mtchs[indx].msc2.push(letters);
        }
      },
      error: (err) => {
        alert("HTTP ERROR: " + JSON.stringify(err));
        console.error('HTTP POST failed:', err);
      }
    });
  }

  edit_msc(indx1:number,indx2:number) {
    if (!this.changeFreeze) {
      this.mscRow = indx1;
      this.mscEdit = indx2;
      this.changeFreeze = true;
    }
  }

  delete_msc(indx1: number, indx2: number, plyrtm:number, letters:string) {
    this.payload_msc.tourn_id = Number(this.tournId);
    this.payload_msc.rnd = Number(this.rnd);
    this.payload_msc.plyrtm_id = Number(plyrtm); 
    this.payload_msc.letters = letters;
    this.payload_msc.add_delete = 0;
    this.http.post<any>(this.APIURL + "trns/AddDeleteMsc", this.payload_msc).subscribe({
      next: (res) => {
        alert(res.msg);
        this.changeFreeze = false;
        if (indx2 < this.mtchs[indx1].msc1.length) {
          this.mtchs[indx1].msc1.splice(indx2, 1);
        } else {
          this.mtchs[indx1].msc2.splice(indx2-this.mtchs[indx1].msc1.length, 1);
        }
      },
      error: (err) => {
        alert("HTTP ERROR: " + JSON.stringify(err));
        console.error('HTTP POST failed:', err);
      }
    });
  }

  cancel_msc() {
    this.mscRow = -1;
    this.mscEdit = -1;
    this.changeFreeze = false;
  }

  cancel_row(i: number) {
    this.mtchs[i].plyrwn = this.Oldmtchs[i].plyrwn;
    this.mtchs[i].plyrwnnm  = this.Oldmtchs[i].plyrwnnm;
    this.mtchs[i].pf = this.Oldmtchs[i].pf;
    this.mtchs[i].pa = this.Oldmtchs[i].pa;
    this.mtchs[i].notes = this.Oldmtchs[i].notes; 
    this.mtchs[i].change = false;
  }

  reload(rndSelect: number) {
    this.rnd = rndSelect;
    this.get_sch();
  }
}

