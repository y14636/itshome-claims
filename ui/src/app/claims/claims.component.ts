import { Component, OnInit, ViewContainerRef } from '@angular/core';
import { ClaimsService, Claims, ModifiedClaims } from '../claims.service';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { EditModalComponent } from '../edit-modal/edit-modal.component';
import { ErrorModalComponent } from '../error-modal/error-modal.component';
import { FormGroup, FormArray, FormBuilder, FormControl, Validators } from '@angular/forms';
import { HttpClient } from '@angular/common/http';
import { HttpErrorResponse } from '@angular/common/http';
import * as $ from 'jquery';
import 'datatables.net';
import 'datatables.net-bs4';
import { NGXLogger } from 'ngx-logger';

@Component({
  selector: 'app-claims',
  templateUrl: './claims.component.html',
	styleUrls: ['./claims.component.css'],
	providers: [NGXLogger]
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
  modifiedClaims: ModifiedClaims[];
  
  selectedActiveInstitutionalClaimIds: Array<any> = [];
	selectedActiveProfessionalClaimIds: Array<any> = [];
	options: string [];
  selectedInstOption: number;
  selectedProfOption: number;
	selectedType: string;
	showInstButton: boolean;
	showProfButton: boolean;
	dataTable: any;
	public instTable: any;
	public profTable: any;
	public modTable: any;
	submittedInstForm = false;
	submittedProfForm = false;

  constructor(
	private logger: NGXLogger,
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
        this.logger.error (err.message);
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
				array.splice(index, 1);
			}    
		});
	}
	private initInstTable(): void {
    let instTable: any = $('#instTable');
    this.instTable = instTable.DataTable({
			searching: false,
			"pagingType": "full_numbers"
    });
	}
	
	private initProfTable(): void {
		let profTable: any = $('#profTable');
		this.profTable = profTable.DataTable({
			searching: false,
			"pagingType": "full_numbers"
		});
	}
	
	private initModTable(): void {
		let modTable: any = $('#modTable');
		this.modTable = modTable.DataTable({
			searching: false,
			"pagingType": "full_numbers",
			"ordering": true,
			"order": [[7, "desc"]]
		});
	}
	
	private reInitInstTable(): void {
    if (this.instTable) {
      this.instTable.destroy();
      this.instTable=null;
    }
    setTimeout(() => this.initInstTable(),0);
	}
	
	private reInitProfTable(): void {
    if (this.profTable) {
      this.profTable.destroy();
      this.profTable=null;
    }
    setTimeout(() => this.initProfTable(),0);
	}
	
	private reInitModTable(): void {
    if (this.modTable) {
			this.logger.debug("modTable exists")
      this.modTable.destroy();
      this.modTable=null;
    }
    setTimeout(() => this.initModTable(),1000);
	}
	
  private createInstForm(option:string) {
    this.instSearchForm = this.formBuilder.group({
	  instSelectItems: this.formBuilder.array([ this.createSelectItems(option) ]),
		instInputItems: this.formBuilder.array([ this.createInputItem(0, "default") ]),
		claimsThreshold: 100,
	  hiddenInputItems: this.formBuilder.array([ this.createHiddenInputItem(0, "ClaimType", "11"), this.createHiddenInputItem(1, "ClaimType", "12") ])
    });
	
		this.logger.debug("Inside createInstForm()");
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
		claimsThreshold: 100,
	  hiddenInputItems: this.formBuilder.array([ this.createHiddenInputItem(0, "ClaimType", "20")])
    });
	
		this.logger.debug("Inside createProfForm()");
		var arrayControl = this.profSearchForm.get('profInputItems') as FormArray;
		var item = arrayControl.at(0);
		this.profSearchForm.get('profSelectItems').valueChanges.subscribe(data => {
			this.selectedProfOption = data[0].category.value;
			item.get("type").setValue(data[0].category.type);
			this.logger.debug("select prof event registered...");
		})
  }
  
  createHiddenInputItem(index:number, type:string, value:string): FormGroup {
	  var inputName = "inputName" + index;
	  this.logger.debug("hiddenInputName=", inputName);
	  return this.formBuilder.group({
		[inputName]: [value],
		type: [type]
	  });
  }
  
  createInputItem(index:number, type:string): FormGroup {
	  var inputName = "inputName" + index;
	  this.logger.debug("inputName=", inputName);
	  return this.formBuilder.group({
		[inputName]: ['', [Validators.required, Validators.minLength(1)]],
		type: [type]
	  });
  }
  
  createSelectItems(option:string): FormGroup {
	  return this.formBuilder.group({
			category: [option]
	  });
  }
  
  addInputItem(index:number, type:string, form:FormGroup, claimType:string): void {
		this.logger.debug("addInputItem() index=", index);
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
				this.logger.debug("Change happened", arrayControl.indexOf(val)+': ', val.get('category').value.name);
				this.logger.debug("Need to update items with new type");
				this.logger.debug("type is", val.get('category').value.type);
				this.logger.debug("select control index=", arrayControl.indexOf(val));
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
		setTimeout(() => this.initInstTable(),0);
		setTimeout(() => this.initProfTable(),0);
		setTimeout(() => this.initModTable(),0);
  }
  
  openEditModal(claimType:string, claims: Claims) {
	  if ( claimType === 'Institutional' && this.selectedActiveInstitutionalClaimIds.length > 0 || 
		claimType === 'Professional' && this.selectedActiveProfessionalClaimIds.length > 0) {
		  const modalRef = this.modalService.open(EditModalComponent, { size: 'lg', backdrop: 'static' });
		  
		  modalRef.componentInstance.title = 'Edit ' + claimType + ' Claim(s)';
		  this.logger.debug("Inside openEditModal, claimType=" + claimType);
		  if (claimType === 'Institutional') {
			this.logger.debug("Inside openEditModal, this.selectedActiveInstitutionalClaimIds[0]=" + this.selectedActiveInstitutionalClaimIds[0]);
			modalRef.componentInstance.selectedActiveInstitutionalClaimIds = this.selectedActiveInstitutionalClaimIds;
		  } else {
			this.logger.debug("Inside openEditModal, this.selectedActiveProfessionalClaimIds[0]=" + this.selectedActiveProfessionalClaimIds[0]);
			modalRef.componentInstance.selectedActiveProfessionalClaimIds = this.selectedActiveProfessionalClaimIds;
		  }
		  
		  modalRef.result.then((result) => {
			window.location.reload();
		  }).catch((error) => {
			this.logger.error(error);
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
			this.logger.debug(result);
	  }).catch((error) => {
			this.logger.error(error);
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
			if (this.instTable) {
				this.reInitInstTable();
			} else {
				setTimeout(() => this.initInstTable(),0);
			}
			if (this.profTable) {
				this.reInitProfTable();
			} else {
				setTimeout(() => this.initProfTable(),0);
			}
    });
  }

  getModifiedClaims() {
	  this.claimsService.getModifiedClaimsList().subscribe((data: ModifiedClaims[]) => {
			this.modifiedClaims = data.filter(claim => claim);
		  });	
  }
	
	searchActiveInstitutionalClaims(strFormData) {
		this.claimsService.getSearchResults(strFormData).subscribe((data: Claims[]) => {
			this.activeInstitutionalClaims = data.filter(claim => claim);
			this.reInitInstTable();
		});
	}

	searchActiveProfessionalClaims(strFormData) {
		this.claimsService.getSearchResults(strFormData).subscribe((data: Claims[]) => {
			this.activeProfessionalClaims = data.filter(claim => claim);
			this.reInitProfTable();
		});
	}
  
  toggleActiveInstitutionalClaims(id:string, isChecked: boolean){
	this.logger.debug("Institutional id=" + id + "isChecked=" + isChecked);
	this.toggleClaims(id, isChecked, 'Institutional');
  }
  
  toggleActiveProfessionalClaims(id:string, isChecked: boolean){
	this.logger.debug("Professional id=" + id + "isChecked=" + isChecked);
	this.toggleClaims(id, isChecked, 'Professional');
  }
  
  toggleClaims(id:string, isChecked: boolean, claimType:string) {
	  this.logger.debug("isChecked=" + isChecked + ", claimType=" + claimType);
	if (claimType === 'Institutional') {
		this.selectedActiveProfessionalClaimIds = [];
		if (isChecked && this.selectedActiveInstitutionalClaimIds.includes(id) === false) {
			this.logger.debug('adding Institutional id');
			this.selectedActiveInstitutionalClaimIds.push(id);
		} else {
			const index: number = this.selectedActiveInstitutionalClaimIds.indexOf(id);
			this.logger.debug('index is ' + index);
			if (index !== -1) {
				this.logger.debug('removing Institutional id');
				this.selectedActiveInstitutionalClaimIds.splice(index, 1);
			}  
		}
	} else {
		this.selectedActiveInstitutionalClaimIds = [];
		if (isChecked && this.selectedActiveProfessionalClaimIds.includes(id) === false) {
			this.logger.debug('adding Professional id');
			this.selectedActiveProfessionalClaimIds.push(id);
		} else {
			const index: number = this.selectedActiveProfessionalClaimIds.indexOf(id);
			if (index !== -1) {
				this.logger.debug('removing Professional id');
				this.selectedActiveProfessionalClaimIds.splice(index, 1);
			}  
		}
	}
  }
  
  getControls(frmGrp: FormGroup, key: string) {
	return (<FormArray>frmGrp.controls[key]).controls;
  }
  removeObject(index, claimType) {
		this.logger.debug("removing index->", index);
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
	
	deleteClaims(mClaims: ModifiedClaims) {
		this.logger.debug("deleted id=", mClaims.id);
    this.claimsService.deleteClaims(mClaims).subscribe(() => {
			this.getModifiedClaims();
			this.reInitModTable();
    })
	}
	
  onSubmit(claimType:string, model: any) {
			this.selectedActiveInstitutionalClaimIds = [];
			this.selectedActiveProfessionalClaimIds = [];
	
			let strFormData = JSON.stringify(model);
			
			if (claimType === 'Institutional') {
				this.submittedInstForm = true;
				if (this.instSearchForm.invalid) {
					return;
				} else {
					this.searchActiveInstitutionalClaims(strFormData);
				}
			} else {
				this.submittedProfForm = true;
				if (this.profSearchForm.invalid) {
					return;
				} else {
					this.searchActiveProfessionalClaims(strFormData);
				}
			}
	 }
	 
	 clearForm(type:string) {
		 if (type === "Institutional") {
			 this.createInstForm(this.options[0]);
			 this.selectedInstOption = 0;
			 this.submittedInstForm = false;
		 } else {
			 this.createProfForm(this.options[0]);
			 this.selectedProfOption = 0;
			 this.submittedProfForm = false;
		 }
		 this.getAll();
	 }
}