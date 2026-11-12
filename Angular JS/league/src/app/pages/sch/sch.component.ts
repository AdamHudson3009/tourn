import { Component, Injector, ViewContainerRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { PotwComponent } from '../potw/potw.component';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { MatSortModule } from '@angular/material/sort';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatInputModule } from '@angular/material/input';
import { finalize } from 'rxjs';
//import { setCookie, getCookie, clearCookie } from '../../../common/cookies';

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
  mscAddAll1: string;
  msc2: string[];
  mscAddAll2: string;
  mscAddMode: boolean;
  letter1: string;
  letter2: string;
  arch_nm: string;
  arch_pts: string;
  ordnum: number;
  has_stnd: string;
}

/*
interface PlfSelect {
  indx: number;
  name: string;
}*/


interface EditPayload {
  gid: number;
  plyrwn: number;
  pf: number;
  pa: number;
  notes: string;
  ordnum: number;
  plyr1: number,
  plyr2: number,
  oldMscAddAll1: string;
  oldMscAddAll2: string;
  mscAddAll1: string;
  mscAddAll2: string;
}

@Component({
  selector: 'app-sch',
  imports: [
    CommonModule, HttpClientModule, FormsModule, ReactiveFormsModule, RouterModule,
    MatIconModule, MatTableModule, MatSortModule, MatPaginatorModule, MatInputModule,
    PotwComponent 
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
  mscAdd1: number = 1;
  mscAdd2: string = "";
  changeFreeze: boolean = false;
  loading: boolean = false;

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
      this.rndSelect = this.rnd;
      this.rndType = 1;
    });
  }

  payload = {
    id: 0,
    rnd: 0
  }

  payload_edits: EditPayload[] = [];

  /*
  payload_msc = {
    tourn_id: 0,
    rnd: 0, 
    plyrtm_id: 0,
    letters: "",
    add_delete: 0
  }*/

  ngOnInit(){
    this.loading = false;   // <-- add this with your other fields
    this.get_sch();
  }

  async get_sch() {
    this.loading = true;
    this.payload.id = this.tournId;
    this.payload.rnd = Number(this.rnd);
  
    await this.http.post<Mtch[]>(this.APIURL + "trns/GetSch", this.payload)
    .pipe(finalize(() => this.loading = false))
    .subscribe({
      next: (res) => {
        this.mtchs = res;
        this.Oldmtchs = res.map(m => JSON.parse(JSON.stringify(m)));
        this.msc_joining(this.mtchs);
        this.msc_joining(this.Oldmtchs);
        this.loading = false;
      },
      error: (err) => {
        alert("HTTP ERROR: " + JSON.stringify(err));
        console.error('HTTP POST failed:', err);
        this.loading = false;
      }
    });
  }

  msc_joining(mtchs: Mtch[]) {
    mtchs.forEach(mtch => {
      mtch.change = false //on reload, so no changes yet
      mtch.mscAddMode = false
      mtch.mscAddAll1 = ""
      mtch.mscAddAll2 = ""
      if (mtch.has_stnd === "N") { //if has_stnd is "N", playoffs, so only one comparison is needed
        mtch.letter1 = "";
        mtch.letter2 = "";
      } else {
        if (mtch.msc1) {
          mtch.mscAddAll1 = mtch.msc1.join(',');
        }
        if (mtch.msc2) {
          mtch.mscAddAll2 = mtch.msc2.join(',');
        }
      }
    });
  }
  
  swap_row(indx:number, up:boolean) {
    let other_indx = up ? indx - 1 : indx + 1;
    const temp_mtch = this.mtchs[indx]
    this.mtchs[indx] = this.mtchs[other_indx];
    this.mtchs[other_indx] = temp_mtch;
  }

 
  async edit_row(indx: number, gid:number, plyrwn:number, plyrwnnm: string, pf: number, pa: number, notes: string) {
    this.payload_edits = [];
    this.payload_edits.push({
      gid: Number(gid), 
      plyrwn: Number(plyrwn), 
      pf: Number(pf), 
      pa: Number(pa), 
      notes: notes ?? "", 
      ordnum: 0,
      plyr1: Number(this.mtchs[indx].plyr1),
      plyr2: Number(this.mtchs[indx].plyr2),
      oldMscAddAll1: this.Oldmtchs[indx].mscAddAll1,
      oldMscAddAll2: this.Oldmtchs[indx].mscAddAll2,
      mscAddAll1: this.mtchs[indx].mscAddAll1,
      mscAddAll2: this.mtchs[indx].mscAddAll2
    }) 

    await this.http.post<any>(this.APIURL + "trns/EditSchs", this.payload_edits).subscribe({
      next: (res) => {
        this.Oldmtchs[indx].plyrwn = plyrwn;
        this.Oldmtchs[indx].plyrwnnm = plyrwnnm;
        this.Oldmtchs[indx].pf = pf ?? 0;
        this.Oldmtchs[indx].pa = pa ?? 0;
        this.Oldmtchs[indx].notes = notes;
        this.Oldmtchs[indx].msc1 = this.mtchs[indx].mscAddAll1.split(","); 
        this.Oldmtchs[indx].msc2 = this.mtchs[indx].mscAddAll2.split(","); 
        this.Oldmtchs[indx].mscAddAll1 = this.mtchs[indx].mscAddAll1;
        this.Oldmtchs[indx].mscAddAll2 = this.mtchs[indx].mscAddAll2; 
        this.mtchs[indx].change = false;
      },
      error: (err) => {
        alert("HTTP ERROR: " + JSON.stringify(err));
        console.error('HTTP POST failed:', err);
      }
    });
  }

  add_msc2(indx: number, oneOrTwo: number, letters: string) {
    if (Number(oneOrTwo) === 1) {
      this.mtchs[indx].mscAddAll1 = letters;
      this.mtchs[indx].msc1 = []
    } else {
      this.mtchs[indx].mscAddAll2 = letters;
      this.mtchs[indx].msc2 = []; //placeholder for existing, use mscAddAll to analyze 
    }
    this.mtchs[indx].mscAddMode = false;
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


  async save_changes() {
    if (!confirm("Do you want to save changes?")) return;

    this.payload_edits = this.mtchs.map((match, idx) => ({
      gid: Number(match.gid),
      plyrwn: Number(match.plyrwn),
      pf: Number(match.pf),
      pa: Number(match.pa),
      ordnum: idx + 1,
      notes: match.notes ?? "",
      plyr1: Number(this.mtchs[idx].plyr1),
      plyr2: Number(this.mtchs[idx].plyr2),
      oldMscAddAll1: this.Oldmtchs[idx].mscAddAll1,
      oldMscAddAll2: this.Oldmtchs[idx].mscAddAll2,
      mscAddAll1: this.mtchs[idx].mscAddAll1,
      mscAddAll2: this.mtchs[idx].mscAddAll2
    }));

    await this.http.post<any>(this.APIURL + "trns/EditSchs", this.payload_edits).subscribe({
      next: (res) => {
        this.mtchs.forEach(mtch => mtch.change = false);
        this.reload(this.rnd);
      },
      error: (err) => {
        alert("HTTP ERROR: " + JSON.stringify(err));
        console.error('HTTP POST failed:', err);
      }
    });

  }

  cancel_changes() {
    if (confirm("Are you sure you want to cancel changes?")) {
      this.mtchs = JSON.parse(JSON.stringify(this.Oldmtchs));
    }
  }

}

