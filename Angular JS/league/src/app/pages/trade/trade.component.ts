import { Component } from '@angular/core';
import { RouterOutlet, RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { MatIconModule } from '@angular/material/icon';

@Component({
  standalone: true,
  selector: 'app-root',
  imports: [RouterOutlet,CommonModule, HttpClientModule, FormsModule, ReactiveFormsModule, RouterModule, MatIconModule],
  templateUrl: './league.component.html',
  styleUrls: ['./league.component.css']
})
export class LeagueComponent {
  title = 'Leagues';
  APIURL="http://127.0.0.1:8000/";
  newLeague="";
  leagues:any=[];
  changeLeague="";
  clickedChange=0;

  constructor(private http:HttpClient){}

  payload = {
    id: 0,
    league: ""
  }

  ngOnInit(){
    this.get_leagues()
  }
  
  get_leagues(){
    this.http.get(this.APIURL+"league/get").subscribe((res)=>{
      this.leagues=res;
    })
  }

  add_league(){
    this.payload.id = 0;
    this.payload.league = this.newLeague;
    this.http.post(this.APIURL+"league/add",this.payload).subscribe((res)=>{
      alert(res);
      this.newLeague="";
      this.get_leagues();
    })
  }

  click_change(num:number){
    this.clickedChange=num;
    if (num === 0) {
      this.changeLeague=""
    }
  }

  change_league(id:any){
    this.payload.id = id;
    this.payload.league = this.changeLeague;
    this.http.post(this.APIURL+"league/change",this.payload).subscribe((res)=>{
      alert(res);
      this.newLeague="";
      this.changeLeague="";
      this.get_leagues();
      this.clickedChange=0;
    })
  }

  delete_league(id:any,description:any){
    this.payload.id = id;
    this.payload.league = description;
    this.http.post(this.APIURL+"league/delete",this.payload).subscribe((res)=>{
      alert(res);
      this.get_leagues();
    })
  }
 
}
