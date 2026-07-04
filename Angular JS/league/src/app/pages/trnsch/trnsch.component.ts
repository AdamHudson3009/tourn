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

interface PlayerTeamInfo {
  plyrtm_id: number;
  plyrtm: string;
  cid: number;
  conf: string;
  did: number;
  divi: string;
  sds: number;
}

interface PlayerNameSchedule {
  plyrtm: string;
  sds: number;
  plyrtms: number[];
  plyrnms: string[];
  plyrnmesep: string[];
}

interface PlayerTeamSchedule {
  plyrtm_id: number;
  plyrtms: number[];
}

interface Division {
  divi_id: number;
  divi: string;
  plyrtms_sch: Record<string, PlayerTeamSchedule>;
  plyrnmsSchArr?: PlayerNameSchedule[];
}

interface Conference {
  conf_id: number;
  conf: string;
  divis: Record<string, Division>;
  diviArr?: Division[];
}

interface TournamentData {
  confs: Record<string, Conference>;
  plyrtm_info: Record<string, PlayerTeamInfo>;
}

@Component({
  selector: 'app-trnsch',
  imports: [
    CommonModule, HttpClientModule, FormsModule, ReactiveFormsModule, RouterModule,
    MatIconModule, TbrksComponent, SchBuilderComponent, CmmComponent, LllConsolidatedComponent,
    GrammarComponent
  ],
  templateUrl: './trnsch.component.html',
  styleUrl: './trnsch.component.css'
})
export class TrnschComponent {
  title = 'Trn Sch';
  APIURL = "http://localhost:8080/";
  leagueId: number = 0;
  leagueName: string = '';
  tournId: number = 0;
  tournName: string = "";
  data: TournamentData = {
    confs: {},
    plyrtm_info: {}
  };
  confsArr: Conference[] = [];
  playertmInfo: PlayerTeamInfo[] = [];
  maxRnd: number = 0;
  rndSelect: number = 0;
  rndType: number = 0;
  isChecked: boolean = false;
  buildNum: number = 0;
  strCheckBox: string = "Check for Div Cmm Oth";
  strCheckMsg1: string = "Check for Div Cmm Oth";
  strCheckMsg2: string = "Uncheck for Rnds";

  constructor(
    private http: HttpClient,
    private route: ActivatedRoute,
  ) {
    this.route.params.subscribe(params => {
      this.leagueId = params['league_id'];
      this.leagueName = params['league_name'];
      this.tournId = Number(params['tourn_id']);
      this.tournName = params['tourn_name'];
    });
  }

  payload_bld = {
    id: 0,
    build: 0
  }

  payload = {
    id: 0
  }
  
  build_sch() {
    this.payload_bld.id = this.tournId;
    this.payload_bld.build = this.buildNum
    this.http.post(this.APIURL + "trns/Build", this.payload_bld).subscribe((res) => {
      this.get_sch()
    });
  }

  get_sch() {
    this.payload.id = this.tournId;
  
    
    this.http.post<TournamentData>(this.APIURL + "trns/Get", this.payload).subscribe((res) => {
      this.data.confs = res.confs;
      this.data.plyrtm_info = res.plyrtm_info;
  
      // Convert to sorted player info array
      this.playertmInfo = Object.values(this.data.plyrtm_info).sort((a, b) => {
        const confCompare = a.conf.localeCompare(b.conf);
        if (confCompare !== 0) return confCompare;
        const diviCompare = a.divi.localeCompare(b.divi);
        if (diviCompare !== 0) return diviCompare;
        return a.sds - b.sds;
      });
  
      // Convert confs to array and sort
      this.confsArr = Object.values(this.data.confs).sort((a, b) =>
        a.conf.localeCompare(b.conf)
      );
  
      for (let i = 0; i < this.confsArr.length; i++) {
        const conf = this.confsArr[i];
        conf.diviArr = Object.values(conf.divis).sort((a, b) =>
          a.divi.localeCompare(b.divi)
        );
  
        for (let i2 = 0; i2 < conf.diviArr.length; i2++) {
          const div = conf.diviArr[i2];
  
          div.plyrnmsSchArr = [];
  
          const teams = Object.values(div.plyrtms_sch || {});
  
          for (let i3 = 0; i3 < teams.length; i3++) {
            const team = teams[i3];
            const playerName = this.playertmInfo.find(p => p.plyrtm_id === team.plyrtm_id)?.plyrtm || "(Unknown)";
            const sds = this.playertmInfo.find(p => p.plyrtm_id === team.plyrtm_id)?.sds || 0;
            if (team.plyrtms.length>this.maxRnd) {
              this.maxRnd = team.plyrtms.length;
            }
            div.plyrnmsSchArr.push({
              plyrtm: playerName,
              sds: sds,
              plyrtms: team.plyrtms,
              plyrnms: [],
              plyrnmesep: []
            });
          }
          div.plyrnmsSchArr.sort((a, b) => a.sds - b.sds);
          for (let i3= 0; i3 < div.plyrnmsSchArr.length; i3++) {
            let cmm: string[] = [];
            let oth: string[] = [];
            for (let i4 = 0; i4 < div.plyrnmsSchArr[i3].plyrtms.length; i4++) {
              const playerId = div.plyrnmsSchArr?.[i3]?.plyrtms?.[i4];
              const playerName = this.playertmInfo.find(p => p.plyrtm_id === playerId)?.plyrtm || "(Unknown)";
              div.plyrnmsSchArr[i3].plyrnms.push(playerName);
              let playerConf = this.playertmInfo.find(p => p.plyrtm_id === playerId)?.conf || "(Unknown)";
              let playerDivi = this.playertmInfo.find(p => p.plyrtm_id === playerId)?.divi || "(Unknown)";
              if ((playerConf === conf.conf) && (playerDivi === div.divi)) {
                div.plyrnmsSchArr[i3].plyrnmesep.push(playerName);
              } else if (playerConf !== conf.conf) {
                cmm.push(playerName);
              } else {
                oth.push(playerName);
              }
            }
            div.plyrnmsSchArr[i3].plyrnmesep.push("CMM");
            for (let i4 = 0; i4 < cmm.length; i4++) {
              div.plyrnmsSchArr[i3].plyrnmesep.push(cmm[i4]);
            }
            div.plyrnmsSchArr[i3].plyrnmesep.push("OTH");
            for (let i4 = 0; i4 < oth.length; i4++) {
              div.plyrnmsSchArr[i3].plyrnmesep.push(oth[i4]);
            }
          }
        }
      }
    });
  }
  
  change_check() {
    if (!this.isChecked) {
      this.strCheckBox = this.strCheckMsg2;
    } else {
      this.strCheckBox = this.strCheckMsg1;
    }
  }
}
