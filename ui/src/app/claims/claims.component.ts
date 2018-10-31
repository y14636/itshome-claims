import { Component, OnInit, ViewContainerRef } from '@angular/core';
import { ClaimsService, Claims } from '../claims.service';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { EditModalComponent } from '../edit-modal/edit-modal.component';
import { ErrorModalComponent } from '../error-modal/error-modal.component';
import { FormGroup, FormArray, FormBuilder, FormControl, Validators } from '@angular/forms';
import { HttpClient } from '@angular/common/http';
import { HttpErrorResponse } from '@angular/common/http';
// import { Subject } from 'rxjs';
// import { DataTableDirective } from 'angular-datatables';
import * as $ from 'jquery';
import 'datatables.net';
import 'datatables.net-bs4';

@Component({
  selector: 'app-claims',
  templateUrl: './claims.component.html',
  styleUrls: ['./claims.component.css']
})
export class ClaimsComponent implements OnInit {
	
  viewMode = 'tab1';
  
	instSearchForm: FormGroup; 
	profSearchForm: FormGroup;
  hiddenInputItems: FormArray;
  instInputItems: FormArray;
  instSelectItems: FormArray;
  profInputItems: FormArray;
  profSelectItems: FormArray;
	
  activeInstitutionalClaims: Claims[];
  activeProfessionalClaims: Claims[];
  modifiedClaims: Claims[];
  
  selectedActiveInstitutionalClaimIds: Array<any> = [];
	selectedActiveProfessionalClaimIds: Array<any> = [];
	options: string [];
  selectedInstOption: number;
  selectedProfOption: number;
	selectedType: string;
	showInstButton: boolean;
	showProfButton: boolean;
	dataTable: any;
	public table1: any;
	public table2: any;
	public table3: any;

  constructor(
  private claimsService: ClaimsService, 
  private modalService: NgbModal, 
	private httpService: HttpClient,
	private formBuilder: FormBuilder
	) { 
		this.createInstForm("");
		this.createProfForm("");						
	}
	
	ngAfterViewInit() {
    
	}
	
	ngOnInit() {		

		this.httpService.get('../assets/options.json').subscribe(
      data => {
				this.options = data["options"];	 // FILL THE ARRAY WITH DATA.
					this.findAndRemove(this.options, 'type', 'ClaimType');// remove hidden fields from dropdown
					this.createInstForm(this.options[0]);
					this.createProfForm(this.options[0]);						
      },
      (err: HttpErrorResponse) => {
        console.log (err.message);
      }
		);

		this.showInstButton = true;
		this.showProfButton = true;

		this.getAll();
	  this.getModifiedClaims();
	}
	
	private findAndRemove(array, property, value) {
		array.forEach(function(result, index) {
			if(result[property] === value) {
				//Remove from array
				array.splice(index, 1);
			}    
		});
	}
	private initDatatable1(): void {
    let table1: any = $('#table1');
    this.table1 = table1.DataTable({
			searching: false,
			"pagingType": "full_numbers"
    });
	}
	
	private initDatatable2(): void {
		let table2: any = $('#table2');
		this.table2 = table2.DataTable({
			searching: false,
			"pagingType": "full_numbers"
		});
	}
	
	private initDatatable3(): void {
		let table3: any = $('#table3');
		this.table3 = table3.DataTable({
			searching: false,
			"pagingType": "full_numbers"
		});
	}
	
	private reInitDatatable1(): void {
    if (this.table1) {
      this.table1.destroy()
      this.table1=null
    }
    setTimeout(() => this.initDatatable1(),0)
	}
	
	private reInitDatatable2(): void {
    if (this.table2) {
      this.table2.destroy()
      this.table2=null
    }
    setTimeout(() => this.initDatatable2(),0)
	}
	
	private reInitDatatable3(): void {
    if (this.table3) {
      this.table3.destroy()
      this.table3=null
    }
    setTimeout(() => this.initDatatable3(),0)
	}
	
  private createInstForm(option:string) {
    this.instSearchForm = this.formBuilder.group({
	  instSelectItems: this.formBuilder.array([ this.createSelectItems(option) ]),
	  instInputItems: this.formBuilder.array([ this.createInputItem(0, "default") ]),
	  hiddenInputItems: this.formBuilder.array([ this.createHiddenInputItem(0, "ClaimType", "11"), this.createHiddenInputItem(1, "ClaimType", "12") ])
    });
	
		console.log("Inside createInstForm()");
		var arrayControl = this.instSearchForm.get('instInputItems') as FormArray;
		var item = arrayControl.at(0);
		this.instSearchForm.get('instSelectItems').valueChanges.subscribe(data => {
			this.selectedInstOption = data[0].category.value;
			item.get("type").setValue(data[0].category.type);
		})
  }
  private createProfForm(option:string) {
    this.profSearchForm = this.formBuilder.group({
	  profSelectItems: this.formBuilder.array([ this.createSelectItems(option) ]),
	  profInputItems: this.formBuilder.array([ this.createInputItem(0, "default") ]),
	  hiddenInputItems: this.formBuilder.array([ this.createHiddenInputItem(0, "ClaimType", "20")])
    });
	
		console.log("Inside createProfForm()");
		var arrayControl = this.profSearchForm.get('profInputItems') as FormArray;
		var item = arrayControl.at(0);
		this.profSearchForm.get('profSelectItems').valueChanges.subscribe(data => {
			this.selectedProfOption = data[0].category.value;
			item.get("type").setValue(data[0].category.type);
			console.log("select prof event registered...");
		})
  }
  
  createHiddenInputItem(index:number, type:string, value:string): FormGroup {
	  var inputName = "inputName" + index;
	  console.log("hiddenInputName=", inputName);
	  return this.formBuilder.group({
		[inputName]: [value],
		type: [type]
	  });
  }
  
  createInputItem(index:number, type:string): FormGroup {
	  var inputName = "inputName" + index;
	  console.log("inputName=", inputName);
	  return this.formBuilder.group({
		[inputName]: [''],
		type: [type]
	  });
  }
  
  createSelectItems(option:string): FormGroup {
	  return this.formBuilder.group({
			category: [option]
	  });
  }
  
  addInputItem(index:number, type:string, form:FormGroup, claimType:string): void {
		console.log("addInputItem() index=", index);
		if (claimType === 'Institutional') {
			this.instInputItems = form.get('instInputItems') as FormArray;
			this.instInputItems.push(this.createInputItem(index, type));
		} else {
			this.profInputItems = form.get('profInputItems') as FormArray;
			this.profInputItems.push(this.createInputItem(index, type));
		}
  }
  
  addSelectItems(form:FormGroup, claimType:string, option:string): void {
		var arrayControl;
		if (claimType === 'Institutional') {
			this.instSelectItems = form.get('instSelectItems') as FormArray;
			if (this.instSelectItems.length < 6) {
				this.instSelectItems.push(this.createSelectItems(option));
				this.instInputItems = form.get('instInputItems') as FormArray;
				this.addInputItem(this.instInputItems.length, "default", form, "Institutional");
				arrayControl = this.getControls(form, 'instSelectItems');
			}
		} else {
			this.profSelectItems = form.get('profSelectItems') as FormArray;
			if (this.profSelectItems.length < 6) {
				this.profSelectItems.push(this.createSelectItems(option));
				this.profInputItems = form.get('profInputItems') as FormArray;
				this.addInputItem(this.profInputItems.length, "default", form, "Professional");
				arrayControl = this.getControls(form, 'profSelectItems');
			}
		}
	
		for(let val of arrayControl) {
			val.get('category').valueChanges.subscribe(data => {
				console.log("Change happened", arrayControl.indexOf(val)+': ', val.get('category').value.name);
				console.log("Need to update items with new type");
				console.log("type is", val.get('category').value.type);
				console.log("select control index=", arrayControl.indexOf(val));
				var inputArrayControl;
				if (claimType === 'Institutional') {
					inputArrayControl = form.get('instInputItems') as FormArray;
				} else {	
					inputArrayControl = form.get('profInputItems') as FormArray;
				}
				var item = inputArrayControl.at(arrayControl.indexOf(val));
				item.get("type").setValue(val.get('category').value.type);
			})
		}
			if (this.instSelectItems != null && this.instSelectItems.length === 6) {
				this.showInstButton = false;
			}

			if (this.profSelectItems != null && this.profSelectItems.length === 6) {
				this.showProfButton = false;
			}
  }
  
  clickTab(tab:string) {
	  this.viewMode = tab;
	  this.selectedActiveInstitutionalClaimIds = [];
		this.selectedActiveProfessionalClaimIds = [];
		setTimeout(() => this.initDatatable1(),0);
		setTimeout(() => this.initDatatable2(),0);
		setTimeout(() => this.initDatatable3(),0);
  }
  
  openEditModal(claimType:string, claims: Claims) {
	  if ( claimType === 'Institutional' && this.selectedActiveInstitutionalClaimIds.length > 0 || 
		claimType === 'Professional' && this.selectedActiveProfessionalClaimIds.length > 0) {
		  const modalRef = this.modalService.open(EditModalComponent, { size: 'lg', backdrop: 'static' });
		  
		  modalRef.componentInstance.title = 'Edit ' + claimType + ' Claim(s)';
		  console.log("Inside openEditModal, claimType=" + claimType);
		  if (claimType === 'Institutional') {
			console.log("Inside openEditModal, this.selectedActiveInstitutionalClaimIds[0]=" + this.selectedActiveInstitutionalClaimIds[0]);
			modalRef.componentInstance.selectedActiveInstitutionalClaimIds = this.selectedActiveInstitutionalClaimIds;
		  } else {
			console.log("Inside openEditModal, this.selectedActiveProfessionalClaimIds[0]=" + this.selectedActiveProfessionalClaimIds[0]);
			modalRef.componentInstance.selectedActiveProfessionalClaimIds = this.selectedActiveProfessionalClaimIds;
		  }
		  
		  modalRef.result.then((result) => {
			window.location.reload();
		  }).catch((error) => {
			console.log(error);
		  });
	  } else {
		  this.openErrorModal();
	  }
  }

  openErrorModal() {
	  const modalRef = this.modalService.open(ErrorModalComponent, {});
	  modalRef.componentInstance.title = 'Error';
	  modalRef.componentInstance.message = 'Please select a claim to edit';
	  modalRef.result.then((result) => {
			console.log(result);
	  }).catch((error) => {
			console.log(error);
	  });
  }

	private buildDtOptions(): DataTables.Settings {
    return {
      searching:false
    };
	}
	
  getAll() {
    this.claimsService.getClaimsList().subscribe((data: Claims[]) => {
      this.activeInstitutionalClaims = data.filter(claim => claim.claimtype === '11' || claim.claimtype === '12');
		this.activeProfessionalClaims = data.filter(claim => claim.claimtype === '20');
		setTimeout(() => this.initDatatable1(),0);
		setTimeout(() => this.initDatatable2(),0);
    });
  }

  getModifiedClaims() {
	  this.claimsService.getModifiedClaimsList().subscribe((data: Claims[]) => {
			this.modifiedClaims = data.filter(claim => claim);
			//setTimeout(() => this.initDatatable3(),0);
		  });	
  }
	
	searchActiveInstitutionalClaims(strFormData) {
		this.claimsService.getSearchResults(strFormData).subscribe((data: Claims[]) => {
			this.activeInstitutionalClaims = data.filter(claim => claim);
			this.reInitDatatable1();
		});
	}

	searchActiveProfessionalClaims(strFormData) {
		this.claimsService.getSearchResults(strFormData).subscribe((data: Claims[]) => {
			this.activeProfessionalClaims = data.filter(claim => claim);
			this.reInitDatatable2();
		});
	}
  
  toggleActiveInstitutionalClaims(id:string, isChecked: boolean){
	console.log("Institutional id=" + id + "isChecked=" + isChecked);
	this.toggleClaims(id, isChecked, 'Institutional');
  }
  
  toggleActiveProfessionalClaims(id:string, isChecked: boolean){
	console.log("Professional id=" + id + "isChecked=" + isChecked);
	this.toggleClaims(id, isChecked, 'Professional');
  }
  
  toggleClaims(id:string, isChecked: boolean, claimType:string) {
	  console.log("isChecked=" + isChecked + ", claimType=" + claimType);
	if (claimType === 'Institutional') {
		this.selectedActiveProfessionalClaimIds = [];
		if (isChecked && this.selectedActiveInstitutionalClaimIds.includes(id) === false) {
			console.log('adding Institutional id');
			this.selectedActiveInstitutionalClaimIds.push(id);
		} else {
			const index: number = this.selectedActiveInstitutionalClaimIds.indexOf(id);
			console.log('index is ' + index);
			if (index !== -1) {
				console.log('removing Institutional id');
				this.selectedActiveInstitutionalClaimIds.splice(index, 1);
			}  
		}
	} else {
		this.selectedActiveInstitutionalClaimIds = [];
		if (isChecked && this.selectedActiveProfessionalClaimIds.includes(id) === false) {
			console.log('adding Professional id');
			this.selectedActiveProfessionalClaimIds.push(id);
		} else {
			const index: number = this.selectedActiveProfessionalClaimIds.indexOf(id);
			if (index !== -1) {
				console.log('removing Professional id');
				this.selectedActiveProfessionalClaimIds.splice(index, 1);
			}  
		}
	}
  }
  
  getControls(frmGrp: FormGroup, key: string) {
	return (<FormArray>frmGrp.controls[key]).controls;
  }
  removeObject(index, claimType) {
		console.log("removing index->", index);
		if (claimType === "Institutional") {
	  	this.instInputItems.removeAt(index);
			this.instSelectItems.removeAt(index);
			if (this.instSelectItems.length < 6) {
				this.showInstButton = true;
			}
		} else {
	  	this.profInputItems.removeAt(index);
			this.profSelectItems.removeAt(index);
			if (this.profSelectItems.length < 6) {
				this.showProfButton = true;
			}
		}
  }
	
	deleteClaims(claims: Claims) {
    this.claimsService.deleteClaims(claims).subscribe(() => {
			this.getModifiedClaims();
			this.reInitDatatable3();
    })
	}
	
  onSubmit(claimType:string, model: any, isValid: boolean, e: any) {
	    e.preventDefault();
			//alert('Form data are: '+JSON.stringify(model));
			let strFormData = JSON.stringify(model);
			
			if (claimType === 'Institutional') {
				this.searchActiveInstitutionalClaims(strFormData);
			} else {
				this.searchActiveProfessionalClaims(strFormData);
			}
	 }
	 
	 clearForm(type:string) {
		 if (type === "Institutional") {
			 this.instSearchForm.reset();
			 this.createInstForm(this.options[0]);
			 this.selectedInstOption = 0;
		 } else {
			 this.profSearchForm.reset();
			 this.createProfForm(this.options[0]);
			 this.selectedProfOption = 0;
		 }
	 }
}