import { Component, OnInit, ViewChild, ElementRef  } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { MatIconModule } from '@angular/material/icon';

interface GrammarInfo {
  id: number;
  letters: string;
  dv1_id: number;
  dv2_id: number;
  divi1: string;
  divi2: string;
  conf1: string;
  conf2: string;
}

interface DiviInfo {
  did: number;
  divi: string; 
  conf: string;
}

interface GrammarResponse {
  grammar: GrammarInfo[];
  confdivi: DiviInfo[];
}

interface TournSelect {
  tourn_id: number;
  tourn: string;
}

@Component({
  selector: 'app-grammar',
  imports: [CommonModule, HttpClientModule, FormsModule, ReactiveFormsModule, RouterModule,
    MatIconModule 
  ],
  templateUrl: './grammar.component.html',
  styleUrl: './grammar.component.css'
})
export class GrammarComponent {
  title = 'tbrk';
  APIURL = "http://127.0.0.1:8000/";
  leagueId: number = 0;
  leagueName: string = '';
  tournId: number = 0;
  tournName: string = '';
  grammarInfo: GrammarInfo[] = [];
  origInfo: GrammarInfo[] = [];
  diviInfo: DiviInfo[] = [];
  gramEdit: number = -1;
  changeFreeze: boolean = false;
  letterAdd: string = "";
  divi1Add: number = 0;
  divi2Add: number = 0;
  tournSelect: TournSelect[] = [];
  fromTrnId: number = -1;
  fileName: string = "";

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
      letters: "",
      divi1_id: 0,
      divi2_id: 0,
      id: 0
    }

    payload_trn = {
      id: 0
    }

    payload_trn_copy = {
      TournIdFrom: 0,
      TournIdTo: 0
    }

    payload_file = {
      tourn_id: 0,
      filename: ""
    }

    ngOnInit(): void {
      this.get_grammars();
    }

    get_grammars() {
      this.payload.tourn_id = this.tournId;
      this.payload.letters = "";
      this.payload.divi1_id = 0;
      this.payload.divi2_id = 0;
      this.payload.id = 0;

      this.http.post<GrammarResponse>(this.APIURL + "grammar/get", this.payload).subscribe((res) => {
        if (res && res.grammar && res.confdivi) {
          this.grammarInfo = res.grammar;
          this.origInfo = JSON.parse(JSON.stringify(this.grammarInfo));
          this.diviInfo = res.confdivi;
          if (this.grammarInfo.length==0) {
            this.get_tourns_dropdown();
          }
        } else {
          console.error("Unexpected response format:", res);
        }
      });
    }

    get_tourns_dropdown() {    
      this.payload_trn.id = this.tournId;
      this.tournSelect = [];

      this.http.post("http://localhost:3000/grammarSch/GetCopyFrom", this.payload_trn).subscribe({
        next: (res) => { 
          if (Array.isArray(res)) {  // Ensure res is an array before iterating
            res.forEach((trn: { tournId: number; tourn: string }) => {
              this.tournSelect.push({ tourn_id: trn.tournId, tourn: trn.tourn });
            });
          }
        },
        error: (err) => alert("API Error: " +JSON.stringify(err)) 
      });
    }

    copy_from(trnIdFrom:number) {
      this.payload_trn_copy.TournIdFrom = trnIdFrom;
      this.payload_trn_copy.TournIdTo = this.tournId;

      this.http.post("http://localhost:3000/grammarSch/CopyFrom",this.payload_trn_copy).subscribe({
        next: (res) => {
          console.log("Success:", res);
          this.reset();
        },
        error: (err) => alert("API Error: " +JSON.stringify(err))
      });
    }

    edit_area(gramEdit: number) {
      if (!this.changeFreeze) {
        this.gramEdit = gramEdit;
        this.changeFreeze = true;
      }
    }

    add() {
      this.payload.tourn_id = this.tournId;
      this.payload.letters = this.letterAdd;
      this.payload.divi1_id = this.divi1Add;
      this.payload.divi2_id = this.divi2Add;
      this.payload.id = 0;

      this.http.post<GrammarResponse>(this.APIURL + "grammar/add", this.payload).subscribe((res) => {
        alert(res);
        this.reset();
      });
    }

    change(id:number,letters:string,d1:number,d2:number) {
      this.payload.tourn_id = this.tournId;
      this.payload.letters = letters;
      this.payload.divi1_id = d1;
      this.payload.divi2_id = d2;
      this.payload.id = id;

      this.http.post<GrammarResponse>(this.APIURL + "grammar/change", this.payload).subscribe((res) => {
        alert(res);
        this.reset();
      });
    }

    cancel() {
      this.grammarInfo = [];
      this.grammarInfo = JSON.parse(JSON.stringify(this.origInfo));
      this.gramEdit = -1;
      this.changeFreeze = false;  

    }

    delete(id: number) {
      this.payload.tourn_id = this.tournId;
      this.payload.letters = "";
      this.payload.divi1_id = 0;
      this.payload.divi2_id = 0;
      this.payload.id = id;

      this.http.post<GrammarResponse>(this.APIURL + "grammar/delete", this.payload).subscribe((res) => {
        alert(res);
        this.reset();
      });
    }

    reset() {
      this.grammarInfo = [];
      this.origInfo = [];
      this.diviInfo = [];
      this.letterAdd = "";
      this.divi1Add = 0;
      this.divi2Add = 0;
      this.changeFreeze = false;
      this.gramEdit = -1;
      this.get_grammars();
    }

    OpenInputFile() {

      const fileInput = document.getElementById('fileInput') as HTMLInputElement | null;

      if (fileInput) {
        fileInput.click();
      } else {
        console.error("Element with id 'fileInput' not found.");
      }
    }

    import_file(filename: string) {
      this.payload_file.tourn_id = Number(this.tournId);
      this.payload_file.filename = filename;

      this.http.post("http://localhost:8080/trns/Import",this.payload_file).subscribe({
        next: (res) => {
          console.log("Success:", res);
          this.reset();
        },
        error: (err) => alert("API Error: " +JSON.stringify(err))
      });
    }

}
