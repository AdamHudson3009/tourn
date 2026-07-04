import { Component, Injector, ViewContainerRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { MatIconModule } from '@angular/material/icon';
import { Overlay, OverlayRef } from '@angular/cdk/overlay';
import { ComponentPortal } from '@angular/cdk/portal';
import { OverlayContentComponent } from './overlay-content.component';  // Assuming you already have the OverlayContentComponent
import { OVERLAY_DATA } from './overlay-data.token'; // Make sure this is imported
import { TbrksComponent } from '../tbrks/tbrks.component';
import { SchBuilderComponent} from "../sch-builder/sch-builder.component"
import { CmmComponent } from '../cmm/cmm.component';
import { LllConsolidatedComponent } from "../lll-consolidated/lll-consolidated.component"


interface ConDivPlr {
  cdp: string;
  id: number;
  description: string;
  description2: string;
  expand: boolean;
  parent_id: number;
  sds: number;
}

@Component({
  selector: 'app-cdp',
  standalone: true,
  imports: [CommonModule, HttpClientModule, FormsModule, ReactiveFormsModule, RouterModule, MatIconModule, 
    TbrksComponent, SchBuilderComponent, CmmComponent, LllConsolidatedComponent],
  templateUrl: './cdp.component.html',
  styleUrls: ['./cdp.component.css']
})
export class CdpComponent {
  private overlayRef!: OverlayRef;

  title = 'Confs';
  APIURL = "http://127.0.0.1:8000/";
  leagueId: number = 0;
  leagueName: string = '';
  tournId: number = 0;
  tournName: string = '';
  condivplr: ConDivPlr[] = [];
  clickedChange = 0;
  changeCdp = "";
  changeCdp2 = "";
  clickedChangeCdp = "";
  isOverlayVisible = false;
  preexpandnum: { cdp: string, id: number}[] = [];
  showRightNum: number = 0;
  
  rows: { value: string }[] = [];
  g_i = 0;
  g_j = 0;

  constructor(
    private http: HttpClient,
    private route: ActivatedRoute,
    private overlay: Overlay,
    private viewContainerRef: ViewContainerRef,
    private injector: Injector
  ) {
    this.route.params.subscribe(params => {
      this.leagueId = params['league_id'];
      this.leagueName = params['league_name'];
      this.tournId = params['tourn_id'];
      this.tournName = params['tourn_name'];
    });
  }

  payload = {
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
    this.get_confs()
  }

  openOverlay() {
    if (this.overlayRef) return; // Prevent multiple overlays

    // Define overlay position
    const positionStrategy = this.overlay.position()
      .global()
      .centerHorizontally()
      .centerVertically();

    console.log('Opening overlay');

    // Create overlay
    this.overlayRef = this.overlay.create({
      hasBackdrop: true,
      backdropClass: 'cdk-overlay-dark-backdrop',
      positionStrategy
    });

    // Define the data to pass
    const dataToPass = { tournId: this.tournId };

    // Create injector with data
    const injector = Injector.create({
      providers: [{ provide: OVERLAY_DATA, useValue: dataToPass }]
    });

    // Attach component with injector
    const overlayPortal = new ComponentPortal(OverlayContentComponent, null, injector);
    this.overlayRef.attach(overlayPortal);

    // Close overlay on backdrop click
    this.overlayRef.backdropClick().subscribe(() => this.closeOverlay());
  }

  closeOverlay() {
    if (this.overlayRef) {
      this.overlayRef.dispose();
      this.overlayRef = undefined!;
    }
  }

  // Your other functions remain unchanged, like get_confs, add_row, etc.
  get_confs() {
    this.payload.parent_id = this.tournId;
    this.payload.table = "conf";
    this.payload.column = "tourn_id";
    this.http.post(this.APIURL + "tcd/get", this.payload).subscribe((res) => {
      if (Array.isArray(res)) {  // Ensure res is an array before iterating
        res.forEach(cdp => {
          this.condivplr.push({ cdp: 'c', id: cdp.id, description: cdp.description, expand: false, parent_id: 0, description2: '', sds: 0 });
        });
      } else {
        console.error("Unexpected response format:", res);
      }
      this.autoexpand();
    });
  }

  expanding(cdp:string,indx:number,item:number,preexp:boolean) {
    this.condivplr[indx].expand = true;
    if (cdp == "c") {
      this.payload.table="divi";
      this.payload.column="conf_id";
    } else if (cdp == "d") {
      this.payload.table="plyrtm";
      this.payload.column="divi_id";
    } 
    this.payload.parent_id=item;
    this.http.post(this.APIURL+"tcd/get",this.payload).subscribe((res)=>{
      if (Array.isArray(res)) {  // Ensure res is an array before iterating
        if (cdp == "c") {
          res.sort((a, b) => a.description - b.description);
          res.forEach(cdpA => {
            this.condivplr.splice(indx+1, 0, { cdp: 'd', id: cdpA.id, description: cdpA.description, expand: false, parent_id: item, description2: '', sds: 0 });
            indx++;
          });
        } else {
          res.sort((a, b) => a.sds - b.sds);
          res.forEach(cdpA => {
            this.condivplr.splice(indx+1, 0, { cdp: 'p', id: cdpA.id, description: cdpA.description, description2: cdpA.description2, sds: cdpA.sds, expand: false, parent_id: item });
            indx++;
          });
        }
        if (preexp) {
          this.g_i++;
          this.g_j++;
          this.autoexpand();
        }
      } else {
        console.error("Unexpected response format:", res);
      }
    })
  }

  contracting(cdp:string,indx:number) {
    let indx2 = indx+1;
    let done = false;
    while (!done) {
      if (indx2>=this.condivplr.length) {
        done = true;
      } else if (this.condivplr[indx2].cdp == "c") {
        done = true;
      } else if (this.condivplr[indx2].cdp == cdp) {
        done = true;
      }  else {
        this.condivplr.splice(indx2, 1);
      }
    }
    this.condivplr[indx].expand = false;
  }

  click_change(cdp:string,num:number,desc:string,visible:boolean,desc2:string="") {
    if (visible && (this.clickedChange>0)) { //do not change if another is opened
      return;
    }
    this.clickedChangeCdp = cdp;
    this.changeCdp = desc;
    this.changeCdp2 = desc2;
    this.clickedChange=num;
    if (num === 0) {
      this.changeCdp=""
      this.changeCdp2=""
    }
  }

  change_conf(table:string,column:string,itemid:number) {
    this.payload.id = itemid;
    this.payload.parent_id = 0;
    this.payload.table = table;
    this.payload.column = column;
    this.payload.descriptions.push(this.changeCdp);
    this.http.post(this.APIURL+"tcd/change",this.payload).subscribe((res)=>{
      alert(res);
      this.changeCdp="";
      this.clickedChange=0;
      this.payload.descriptions.pop();
      this.condivplr = [];
      this.get_confs();
    })
  }

  delete_conf(table:string,column:string,itemid:number,description:string) {
    this.payload.id = itemid;
    this.payload.parent_id = 0;
    this.payload.table = table;
    this.payload.column = column;
    this.payload.descriptions.push(description);
    alert(this.payload.id);
    this.http.post(this.APIURL+"tcd/delete",this.payload).subscribe((res)=>{
      alert(res);
      this.payload.descriptions.pop();
      this.condivplr = [];
      this.get_confs();
    })
  }

  change_plyrtm(itemid:number) {
    this.payload_plyrtm.descriptions = [];
    this.payload_plyrtm.descriptions2 = [];
    this.payload_plyrtm.id = itemid;
    this.payload_plyrtm.parent_id = 0;
    this.payload_plyrtm.descriptions.push(this.changeCdp);
    if (this.changeCdp2==null) {
      this.payload_plyrtm.descriptions2.push("");
    } else {
      this.payload_plyrtm.descriptions2.push(this.changeCdp2);
    }
    this.payload_plyrtm.sds = 0;
    this.payload_plyrtm.up = true;
    this.http.post(this.APIURL+"plyr/change",this.payload_plyrtm).subscribe((res)=>{
      alert(res);
      this.changeCdp="";
      this.changeCdp2="";
      this.clickedChange=0;
      this.payload_plyrtm.descriptions.pop();
      this.payload_plyrtm.descriptions2.pop();
      this.condivplr = [];
      this.get_confs();
    })
  }

  delete_plyrtm(itemid:number,description:string,description2:string) {
    this.payload_plyrtm.descriptions = [];
    this.payload_plyrtm.descriptions2 = [];
    this.payload_plyrtm.id = itemid;
    this.payload_plyrtm.parent_id = 0;
    this.payload_plyrtm.descriptions.push(description);
    if (this.changeCdp2==null) {
      this.payload_plyrtm.descriptions2.push("");
    } else {
      this.payload_plyrtm.descriptions2.push(this.changeCdp2);
    }
    this.payload_plyrtm.sds = 0;
    this.payload_plyrtm.up = true;
    alert(this.payload_plyrtm);
    this.http.post(this.APIURL+"plyr/delete",this.payload_plyrtm).subscribe((res)=>{
      alert(res);
      this.payload_plyrtm.descriptions.pop();
      this.payload_plyrtm.descriptions2.pop();
      this.condivplr = [];
      this.get_confs();
    })
  }

  sds(itemid:number,up:boolean) {
    this.payload_plyrtm.id = itemid;
    this.payload_plyrtm.up = up;
    this.http.post(this.APIURL+"plyr/sds",this.payload_plyrtm).subscribe((res)=>{
      //alert(res);
      this.preexpand();
    })
  }

  preexpand() {
    this.preexpandnum = [];
    for (let i = 0; i < this.condivplr.length; i++) {
      if (this.condivplr[i].expand) {
        this.preexpandnum.push({ cdp: this.condivplr[i].cdp, id: this.condivplr[i].id});
      }
    }
    this.g_i=0;
    this.g_j=0;
    this.condivplr = [];
    this.get_confs();
  }

  autoexpand() {
    while ((this.g_j<this.preexpandnum.length) && (this.g_i<this.condivplr.length)) {   
      if ((this.condivplr[this.g_i].cdp == this.preexpandnum[this.g_j].cdp) && (this.condivplr[this.g_i].id == this.preexpandnum[this.g_j].id)) {
        this.expanding(this.condivplr[this.g_i].cdp,this.g_i,this.condivplr[this.g_i].id,true);
      }
      this.g_i++;
    }
  }

}
