import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-search',
  templateUrl: './search.component.html',
  styleUrls: ['./search.component.css']
})
export class SearchComponent implements OnInit {
	_ref:any;
  procedureCode: string;
  constructor() { }

  removeObject(){
    this._ref.destroy();
  }
  
  // save(){
    // if(this.procedureCode)
      // alert(`procedureCode: ${this.procedureCode}`);
    // else
      // alert('Please enter value to save');
  // }

  ngOnInit() {
  }

}
