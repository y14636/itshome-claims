import { Component, OnInit, ViewChild, ComponentFactoryResolver, ViewContainerRef } from '@angular/core';
import { ClaimsService, Claims } from '../claims.service';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { EditModalComponent } from '../edit-modal/edit-modal.component';
import { ErrorModalComponent } from '../error-modal/error-modal.component';
import { FormGroup, FormArray, FormBuilder, FormControl, Validators } from '@angular/forms';

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
  
  dtOptions: DataTables.Settings[] = [];
  
  activeInstitutionalClaims: Claims[];
  activeProfessionalClaims: Claims[];
  modifiedClaims: Claims[];
  
  selectedActiveInstitutionalClaimIds: Array<any> = [];
  selectedActiveProfessionalClaimIds: Array<any> = [];
 
  options = [
    { name: "Make a selection", value: 0, type: "default" },
    { name: "Receipt Date", value: 1, type: "ReceiptDate" },
    { name: "Claims Threshold", value: 2, type: "ClaimsThreshold" },
		{ name: "Provider ID", value: 3, type: "ProviderId" },
		{ name: "Provider Type", value: 4, type: "ProviderType" },
		{ name: "Provider Specialty", value: 5, type: "ProviderSpecialty" },
		{ name: "Procedure Code", value: 6, type: "ProcedureCode" },
		{ name: "Diagnosis Code", value: 7, type: "DiagnosisCode" },
		{ name: "Subscriber ID", value: 8, type: "SubscriberId" },
		{ name: "Subscriber Prefix", value: 9, type: "SubscriberPrefix" },
		{ name: "Subscriber Suffix", value: 10, type: "SubscriberSuffix" },
		{ name: "SCCF Number", value: 11, type: "SCCFNumber" },
		{ name: "Revenue Code", value: 12, type: "RevenueCode" },
		{ name: "Bill Type", value: 13, type: "BillType" },
		{ name: "Modifier", value: 14, type: "Modifier" },
		{ name: "Plan Code", value: 15, type: "PlanCode" },
		{ name: "SF Message Code", value: 16, type: "SFMessageCode" },
		{ name: "Pricing Method", value: 17, type: "PricingMethod" },
		{ name: "Pricing Rule", value: 18, type: "PricingRule" },
		{ name: "Delivery Method", value: 19, type: "DeliveryMethod" },
		{ name: "From Date (DOS)", value: 20, type: "FromDate" },
		{ name: "To Date (DOS)", value: 21, type: "ToDate" },
		{ name: "Patient Account Number", value: 22, type: "PatientAccountNumber" }
  ]

  selectedInstOption: number;
  selectedProfOption: number;
	selectedType: string;
	showInstButton: boolean;
	showProfButton: boolean;
  
  constructor(
  private claimsService: ClaimsService, 
  private modalService: NgbModal, 
  private _cfr: ComponentFactoryResolver,
  private formBuilder: FormBuilder
	) { 
		this.createInstForm();
		this.createProfForm();
	}
  
  private createInstForm() {
    this.instSearchForm = this.formBuilder.group({
	  instSelectItems: this.formBuilder.array([ this.createSelectItems() ]),
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
  private createProfForm() {
    this.profSearchForm = this.formBuilder.group({
	  profSelectItems: this.formBuilder.array([ this.createSelectItems() ]),
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
  
  createSelectItems(): FormGroup {
	  return this.formBuilder.group({
		category: [this.options[0]]
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
  
  addSelectItems(form:FormGroup, claimType:string): void {
		var arrayControl;
		if (claimType === 'Institutional') {
			this.instSelectItems = form.get('instSelectItems') as FormArray;
			if (this.instSelectItems.length < 6) {
				this.instSelectItems.push(this.createSelectItems());
				this.instInputItems = form.get('instInputItems') as FormArray;
				this.addInputItem(this.instInputItems.length, "default", form, "Institutional");
				arrayControl = this.getControls(form, 'instSelectItems');
			}
		} else {
			this.profSelectItems = form.get('profSelectItems') as FormArray;
			if (this.profSelectItems.length < 6) {
				this.profSelectItems.push(this.createSelectItems());
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
			//console.log("select items size=", this.instSelectItems.length);
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

  ngOnInit() {
		this.showInstButton = true;
		this.showProfButton = true;
    this.dtOptions[0] = this.buildDtOptions();
    this.dtOptions[1] = this.buildDtOptions();
    this.dtOptions[2] = this.buildDtOptions();
		this.getAll();
	  this.getModifiedClaims();
  }

	private buildDtOptions(): DataTables.Settings {
    return {
      searching:false
    };
	}
	
  getAll() {
    this.claimsService.getClaimsList().subscribe((data: Claims[]) => {
      this.activeInstitutionalClaims = data.filter(claim => claim.claimtype === '11');
	  this.activeProfessionalClaims = data.filter(claim => claim.claimtype === '20');
    });
  }

  getModifiedClaims() {
	  this.claimsService.getModifiedClaimsList().subscribe((data: Claims[]) => {
		  this.modifiedClaims = data.filter(claim => claim);
	  });	
  }
	
	searchActiveInstitutionalClaims(strFormData) {
		this.claimsService.getSearchResults(strFormData).subscribe((data: Claims[]) => {
			this.activeInstitutionalClaims = data.filter(claim => claim);
		});
	}

	searchActiveProfessionalClaims(strFormData) {
		this.claimsService.getSearchResults(strFormData).subscribe((data: Claims[]) => {
			this.activeProfessionalClaims = data.filter(claim => claim);
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
  
  onSubmit(claimType:string, model: any, isValid: boolean, e: any) {
	  e.preventDefault();
			alert('Form data are: '+JSON.stringify(model));
			let strFormData = JSON.stringify(model);
			
			if (claimType === 'Institutional') {
				this.searchActiveInstitutionalClaims(strFormData);
			} else {
				this.searchActiveProfessionalClaims(strFormData);
			}
   }
}