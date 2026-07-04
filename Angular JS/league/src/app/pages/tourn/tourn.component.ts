import { HttpParams } from '@angular/common/http';
import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { MatIconModule } from '@angular/material/icon';

@Component({
  standalone: true,
  selector: 'app-tourn',
  imports: [CommonModule, HttpClientModule, FormsModule, ReactiveFormsModule, RouterModule, MatIconModule], 
  templateUrl: './tourn.component.html',
  styleUrl: './tourn.component.css'
})
export class TournComponent {
  title = 'Tourns';
  APIURL="http://127.0.0.1:8000/";
  leagueId: any=0;
  leagueName: any='';
  tourns:any=[];
  clickedChange = 0;
  changeTourn="";
  rows: { value: string }[] = [];

  constructor(private http: HttpClient, private route: ActivatedRoute) {
    this.route.params.subscribe(params => {
      this.leagueId = params['league_id'];
      this.leagueName = params['league_name']; 
    });
  }
  
  payload = {
    table: "tourn",
    column: "league_id",
    parent_id: this.leagueId,
    id: 0,
    descriptions: [] as string[]
  }
  
  ngOnInit(){
    this.get_tourns()
  }
  
  click_change(num:number){
    this.clickedChange=num;
    if (num === 0) {
      this.changeTourn=""
    }
  }

  get_tourns(){
    this.payload.parent_id=this.leagueId;
    this.http.post(this.APIURL+"tcd/get",this.payload).subscribe((res)=>{
      this.tourns=res;
    })
  }

  add_row(){
    this.rows.push({ value: '' }); 
  }

  delete_row(index: number){
    this.rows.splice(index, 1);
  }

  add_tourns(){
    // Make sure descriptions is not empty before adding it to the payload
    if (this.rows.length === 0 || this.rows.some(row => row.value.trim() === '')) {
      alert('Please fill in all the descriptions.');
      return;
    }
    this.payload.parent_id = this.leagueId;
    this.payload.descriptions = this.rows.map(row => row.value.trim());
    this.http.post(this.APIURL+"tcd/add",this.payload).subscribe((res)=>{
      alert(res);
      this.rows=[];
      this.payload.descriptions=[];
      this.get_tourns();
    })
  }

  change_tourn(id:number,description:string){
    this.payload.id = id;
    this.payload.parent_id=this.leagueId;
    this.payload.descriptions.push(description);
    this.leagueId;
    this.http.post(this.APIURL+"tcd/change",this.payload).subscribe((res)=>{
      alert(res);
      this.changeTourn="";
      this.payload.descriptions=[];
      this.clickedChange=0;
      this.get_tourns();
    })
  }

  delete_tourn(id:number,description:string){
    this.payload.id = id;
    this.payload.parent_id=this.leagueId;
    this.payload.descriptions.push(description);
    this.http.post(this.APIURL+"tcd/delete",this.payload).subscribe((res)=>{
      alert(res);
      this.payload.descriptions=[];
      this.get_tourns();
    })
  }

}
