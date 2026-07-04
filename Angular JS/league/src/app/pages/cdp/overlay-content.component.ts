import { Component, Inject, EventEmitter, Output } from '@angular/core'; 
import { CommonModule } from '@angular/common';
import { OVERLAY_DATA } from './overlay-data.token';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { MatDialogRef } from '@angular/material/dialog';
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { MatIconModule } from '@angular/material/icon';
import { FormBuilder, FormGroup, FormArray, Validators } from '@angular/forms';
import { ChangeDetectorRef } from '@angular/core';
//import { OverlayRef } from '@angular/cdk/overlay';

@Component({
  selector: 'app-overlay-content',
  templateUrl: './overlay-content.component.html',
  standalone: true,
  imports: [CommonModule, HttpClientModule, FormsModule, ReactiveFormsModule, RouterModule, MatIconModule],
  styleUrls: ['./overlay-content.component.css']
})
export class OverlayContentComponent {
  @Output() closeOverlay = new EventEmitter<void>();

  cdpForm: FormGroup;
  APIURL = "http://127.0.0.1:8000/";
  cdpControl:string='';
  dropdownConf: any[] = [];
  dropdownDivi: any[] = [];
  
  constructor(@Inject(OVERLAY_DATA) public data: any, private http: HttpClient, private fb: FormBuilder,
  private cdr: ChangeDetectorRef) { //,
    this.cdpForm = this.fb.group({
      cdpControl: [{ value: '', disabled: false }], // Disable the radio button initially
      rows: this.fb.array([]), // initialize the FormArray
      selectedConfs: [0],
      selectedDivis: [0]
    });
    
  }

  onSelectChange(event: Event) {
    this.loadDropdownDivi(this.cdpForm.get('selectedConfs')?.value);
  }

  payload = {
    table: "conf",
    column: "tourn_id",
    parent_id: 0,
    id: 0,
    descriptions: [] as string[]
  }

  payload_plyrtm = {
    parent_id: 0,
    id: 0,
    descriptions: [] as string[],
    descriptions2: [] as string[],
    sds: 0,
    up: false
  }

  ngOnInit(): void {
    this.loadDropdownConf();
  }


  loadDropdownConf() {
    this.payload.table = "conf";
    this.payload.column = "tourn_id";
    this.payload.parent_id = this.data?.tournId ;
    this.http.post(this.APIURL + "tcd/get", this.payload).subscribe((res) => {
      if (Array.isArray(res)) {  // Ensure res is an array before iterating
        res.forEach(cdp => {
          this.dropdownConf.push({ id: cdp.id, description: cdp.description });
        });
      } else {
        console.error("Unexpected response format:", res);
      }
    });
  }

  loadDropdownDivi(parent_id: number) {
    this.dropdownDivi = [];
    this.payload.table = "divi";
    this.payload.column = "conf_id";
    this.payload.parent_id = parent_id;
    this.http.post(this.APIURL + "tcd/get", this.payload).subscribe((res) => {
      if (Array.isArray(res)) {  // Ensure res is an array before iterating
        res.forEach(cdp => {
          this.dropdownDivi.push({ id: cdp.id, description: cdp.description });
        });
      } else {
        console.error("Unexpected response format:", res);
      }
    });
  }

  get rows(): FormArray {
    return this.cdpForm.get('rows') as FormArray;
  }

  close() {
    this.closeOverlay.emit(); // Emit an event to close the overlay
  }

  add_row(table: string){
    if (table != "plyrtm") {
      this.rows.push(this.fb.group({ value: '' }));
    } else {
      this.rows.push(this.fb.group({ value: '', value2: '' }));
    }
    this.cdpForm.controls['cdpControl']?.disable();
    this.cdr.detectChanges();  // Manually trigger change detection if needed
  }

  delete_row(index: number){
    this.rows.removeAt(index);
    if (this.rows.length==0) {
      this.cdpForm.controls['cdpControl']?.enable();
    }
    this.cdr.detectChanges();  // Manually trigger change detection if needed
  }

  add_cdp(table: string, column: string, parentId: number){
      const allRowsHaveValues = this.rows.controls.every(control => {
      const value = control.get('value')?.value;
      return value !== null && value !== undefined && value.trim() !== '';
    });
  
    if (!allRowsHaveValues) {
      alert('Please fill in all the descriptions.');
      return;
    }

   if (table != "plyrtm") {
    this.payload.table = table;
    this.payload.column = column;
    this.payload.id = 0;
    this.payload.parent_id = parentId;
    this.payload.descriptions = this.rows.controls.map(row => row.get('value')?.value.trim());
    this.http.post(this.APIURL+"tcd/add",this.payload).subscribe((res)=>{
      alert(res);
      this.rows.controls=[];
      this.payload.descriptions=[];
    })
    } else {
      this.payload_plyrtm.parent_id = parentId;
      this.payload_plyrtm.descriptions = this.rows.controls.map(row => row.get('value')?.value?.trim() || "");
      this.payload_plyrtm.descriptions2 = this.rows.controls.map(row => row.get('value2')?.value?.trim() || "");
      this.payload_plyrtm.sds = 0;
      this.payload_plyrtm.up = true;
      this.http.post(this.APIURL+"plyr/add",this.payload_plyrtm).subscribe((res)=>{
        alert(res);
        this.rows.controls=[];
        this.payload_plyrtm.descriptions=[];
        this.payload_plyrtm.descriptions2=[];
        //this.get_tourns();
      })
    }
  }
}
