import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { MatIconModule } from '@angular/material/icon';

interface SchSdsInfo {
  w2id: number;
  sd1: number;
  sd2: number;
}

interface SchInfo {
  rnd: number;
  rpt: string;
  gid: number;
  w1id: number;
  letters: string;
  HasStnd: boolean;
  schSdsInfo: SchSdsInfo[];
}

interface OrigSchSds {
  w2id: number[];
  sd1: number[];
  sd2: number[];
}

interface OrigSch {
  rnd: number[];
  rpt: string[];
  gid: number[];
  HasStnd: boolean[];
  origSchSds: OrigSchSds[];
}

interface LettersInfo {
  gid: number;
  letters: string; 
}

interface SchResponse {
  grammar: { letters: string; gid: number }[];
  rnds: {
    rnd: number;
    entries: {
      w1id: number;
      rpt: string;
      gid: number;
      HasStnd: boolean;
      w2Info: { w2id: number; sd1: number; sd2: number }[];
    }[];
  }[];
}

@Component({
  selector: 'app-sch-builder',
  imports: [CommonModule, HttpClientModule, FormsModule, ReactiveFormsModule, RouterModule,
    MatIconModule ],
  templateUrl: './sch-builder.component.html',
  styleUrl: './sch-builder.component.css'
})
export class SchBuilderComponent implements OnInit {
  title = 'sch builder';
  APIURL = "http://127.0.0.1:8000/";
  leagueId: number = 0;
  leagueName: string = '';
  tournId: number = 0;
  tournName: string = '';
  schInfo: SchInfo[] = [];
  changeFreeze: boolean = false;
  schEdit: number = -1;
  rndAdd: number = 0;
  rptAdd: string = "";
  HasStndAdd: boolean = true;
  gidAdd: number = 0;
  sdsAdd: SchSdsInfo[] = [];
  lettersInfo: LettersInfo[] = [];
  origSch: OrigSch= {
    rnd: [],
    rpt: [],
    gid: [],
    HasStnd: [],
    origSchSds: []
  };

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
    gid: 0,
    rnd: 0,
    rpt: "",
    w1id: 0,
    HasStnd: true,
    w2id: [] as number[],
    sd1: [] as number[],
    sd2: [] as number[]
  }
     
  ngOnInit(): void {
    this.get_sch();
  }
  
  edit_area(sch: number) {
    if (!this.changeFreeze) {
      this.schEdit = sch;
    }
  }

  get_sch() {
    this.payload.tourn_id = this.tournId;
    this.payload.gid = 0;
    this.payload.rnd = 0;
    this.payload.rpt = "";
    this.payload.w1id = 0;
    this.payload.HasStnd = true;
    this.payload.w2id = [];
    this.payload.sd1 = [];
    this.payload.sd2 = [];

    this.http.post<SchResponse>(this.APIURL + "lllsch/get", this.payload).subscribe((res) => {
      if (res) {
        if (res.grammar) {
          this.lettersInfo = res.grammar;
        } else {
          console.error("Missing grammar in response:", res);
        } 
        res.rnds.forEach(rndGroup => {
          rndGroup.entries.forEach(entry => {
            this.schInfo.push({ rnd: rndGroup.rnd, rpt: entry.rpt, gid: entry.gid, w1id: entry.w1id, HasStnd: entry.HasStnd, letters: "", schSdsInfo: []});
            this.origSch.gid.push(entry.gid);
            this.origSch.rnd.push(rndGroup.rnd);
            this.origSch.rpt.push(entry.rpt);
            this.origSch.HasStnd.push(entry.HasStnd);
            this.origSch.origSchSds.push({w2id: [], sd1: [], sd2: []});
            let last = this.schInfo.length-1;
            this.schInfo[last].letters = this.IdToLetters(entry.gid)
            entry.w2Info.forEach(w2 => {
              this.schInfo[last].schSdsInfo.push({w2id:w2.w2id, sd1:w2.sd1, sd2:w2.sd2});
              this.origSch.origSchSds[last].w2id.push(w2.w2id);
              this.origSch.origSchSds[last].sd1.push(w2.sd1);
              this.origSch.origSchSds[last].sd2.push(w2.sd2);
            });
          });
        });
      } else {
        console.error("Unexpected response format:", res);
      }
    });
  }
  
  cancel(indx:number) {
    if (indx<this.schInfo.length) {
      this.schInfo[indx].rnd = this.origSch.rnd[indx];
      this.schInfo[indx].rpt = this.origSch.rpt[indx];
      this.schInfo[indx].HasStnd = this.origSch.HasStnd[indx];
      this.schInfo[indx].gid = this.origSch.gid[indx];
      this.schInfo[indx].letters = this.IdToLetters(this.schInfo[indx].gid);
      this.schInfo[indx].schSdsInfo = [];
      for (let indx2=0; indx2<this.origSch.origSchSds[indx].sd1.length; indx2++) {
        this.schInfo[indx].schSdsInfo.push({w2id: this.origSch.origSchSds[indx].w2id[indx2], 
          sd1: this.origSch.origSchSds[indx].sd1[indx2], 
          sd2: this.origSch.origSchSds[indx].sd2[indx2]} );
      }
    }
    this.schEdit = -1;
    this.changeFreeze = false;
  }

  change_sch(indx:number) {
    this.payload.tourn_id = this.tournId;
    this.payload.gid = this.schInfo[indx].gid;
    this.payload.rnd = this.schInfo[indx].rnd;
    this.payload.rpt = this.schInfo[indx].rpt;
    this.payload.w1id = this.schInfo[indx].w1id;
    this.payload.HasStnd = this.schInfo[indx].HasStnd;
    this.payload.w2id = [];
    this.payload.sd1 = [];
    this.payload.sd2 = [];
    for (let indx2=0; indx2<this.schInfo[indx].schSdsInfo.length; indx2++) {
      this.payload.w2id.push(this.schInfo[indx].schSdsInfo[indx2].w2id);
      this.payload.sd1.push(this.schInfo[indx].schSdsInfo[indx2].sd1);
      this.payload.sd2.push(this.schInfo[indx].schSdsInfo[indx2].sd2);
    }
    this.http.post(this.APIURL + "lllsch/change", this.payload).subscribe((res) => {
        alert(res)
        this.reset();
    });
  }

  delete_sch(indx:number) {
    this.payload.tourn_id = this.tournId;
    this.payload.gid = this.schInfo[indx].rnd;;
    this.payload.rnd = this.schInfo[indx].rnd;
    this.payload.rpt = this.schInfo[indx].rpt;
    this.payload.w1id = this.schInfo[indx].w1id;
    this.payload.HasStnd = this.schInfo[indx].HasStnd;
    this.payload.w2id = [];
    this.payload.sd1 = [];
    this.payload.sd2 = [];

    this.http.post(this.APIURL + "lllsch/delete", this.payload).subscribe((res) => {
      alert(res)
      this.reset();
    });
  }

  add_sch() {
    this.payload.tourn_id = this.tournId;
    this.payload.gid = this.gidAdd;
    this.payload.rnd = this.rndAdd;
    this.payload.rpt = this.rptAdd;
    this.payload.HasStnd = this.HasStndAdd;
    this.payload.w1id = 0;
    this.payload.w2id = [];
    this.payload.sd1 = [];
    this.payload.sd2 = [];
    for (let indx=0; indx<this.sdsAdd.length; indx++) {
      this.payload.w2id.push(this.sdsAdd[indx].w2id);
      this.payload.sd1.push(this.sdsAdd[indx].sd1);
      this.payload.sd2.push(this.sdsAdd[indx].sd2);
    }

    this.http.post(this.APIURL + "lllsch/add", this.payload).subscribe((res) => {
      alert(res)
      this.reset();
    });
  }

  add_sds(indx:number) {
    if (indx < this.schInfo.length) {
      this.schInfo[indx].schSdsInfo.push({w2id:0, sd1:1, sd2:1});
    } else {
      this.sdsAdd.push({w2id:0, sd1:1, sd2:1});
    }
    this.changeFreeze = true;
  }

  delete_sds(indx:number,indx2:number) {
    if (indx< this.schInfo.length) {
      this.schInfo[indx].schSdsInfo.splice(indx2,1);
    } else {
      this.sdsAdd.splice(indx2,1);
    }
    this.changeFreeze = true;
  }

  IdToLetters(gid:number): string {
    for (let indx=0; indx<this.lettersInfo.length; indx++) {
      if (this.lettersInfo[indx].gid == gid) {
        return this.lettersInfo[indx].letters;
      }
    }
    return "";
  }

  lettersSelected(indx:number) {
    this.changeFreeze = true;
    this.schInfo[indx].letters = this.IdToLetters(this.schInfo[indx].gid);
  }


  reset() {
    this.payload.w2id = [];
    this.payload.sd1 = [];
    this.payload.sd2 = [];
    this.schInfo = [];
    this.lettersInfo = [];
    this.rndAdd = 0;
    this.rptAdd = "";
    this.HasStndAdd = true;
    this.gidAdd = 0;
    this.sdsAdd = [];
    this.origSch.rnd = [];
    this.origSch.rpt = [];
    this.origSch.gid = [];
    this.origSch.origSchSds = [];
    this.changeFreeze = false;
    this.schEdit = -1;
    this.get_sch();
  }
}
