import { Component, OnInit, ViewChild, ComponentFactoryResolver, ViewContainerRef } from '@angular/core';
import { ClaimsService, Claims } from '../claims.service';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { EditModalComponent } from '../edit-modal/edit-modal.component';
import { ErrorModalComponent } from '../error-modal/error-modal.component';
//import { SearchComponent } from '../search/search.component';
import { FormGroup, FormArray, FormBuilder, FormControl, Validators } from '@angular/forms';

@Component({
  selector: 'app-claims',
  templateUrl: './claims.component.html',
  styleUrls: ['./claims.component.css']
})
export class ClaimsComponent implements OnInit {
	
  viewMode = 'tab1';
  
  searchForm: FormGroup;
  //procedureCodeItems: FormArray;
  inputItems: FormArray;
  selectItems: FormArray;
  
  dtOptions: DataTables.Settings = {};
  
  activeInstitutionalClaims: Claims[];
  activeProfessionalClaims: Claims[];
  modifiedClaims: Claims[];
  
  selectedActiveInstitutionalClaimIds: Array<any> = [];
  selectedActiveProfessionalClaimIds: Array<any> = [];

	// claimtype: string;
	// serviceId: string;
	// receiptDate: string;
  // fromDate: string;
  // toDate: string;
  // placeOfService: string;
  // providerId: string;
  // providerType: string;
  // providerSpecialty: string;
  // procedureCode: string;
  // diagnosisCode: string;
  // networkIndicator: string;
  // subscriberId: string;
  // patientAccountNumber: string;
  // sccfNumber: string;
  // revenueCode: string;
  // billType: string;
  // modifier: string;
  // planCode: string;
  // sfMessageCode: string;
  // pricingMethod: string;
  // pricingRule: string;
  // deliveryMethod: string;
  // inputDate: string;
  // fileName: string;
  
  //@ViewChild('parent', { read: ViewContainerRef }) container: ViewContainerRef;
  
  options = [
    { name: "Make a selection", value: 0, type: "default" },
    { name: "Receipt Date", value: 1, type: "receiptDate" },
    { name: "Claims Threshold", value: 2, type: "claimsThreshold" },
		{ name: "Provider ID", value: 3, type: "providerId" },
		{ name: "Provider Type", value: 4, type: "providerType" },
		{ name: "Provider Specialty", value: 5, type: "providerSpecialty" },
		{ name: "Procedure Code", value: 6, type: "procedureCode" },
		{ name: "Diagnosis Code", value: 7, type: "diagnosisCode" },
		{ name: "Subscriber ID", value: 8, type: "subscriberId" },
		{ name: "Subscriber Prefix", value: 9, type: "subscriberPrefix" },
		{ name: "Subscriber Suffix", value: 10, type: "subscriberSuffix" },
		{ name: "SCCF Number", value: 11, type: "sccfNumber" },
		{ name: "Revenue Code", value: 12, type: "revenueCode" },
		{ name: "Bill Type", value: 13, type: "billType" },
		{ name: "Modifier", value: 14, type: "modifier" },
		{ name: "Plan Code", value: 15, type: "planCode" },
		{ name: "SF Message Code", value: 16, type: "sfMessageCode" },
		{ name: "Pricing Method", value: 17, type: "pricingMethod" },
		{ name: "Pricing Rule", value: 18, type: "pricingRule" },
		{ name: "Delivery Method", value: 19, type: "deliveryMethod" },
		{ name: "From Date (DOS)", value: 20, type: "fromDate" },
		{ name: "To Date (DOS)", value: 21, type: "toDate" },
		{ name: "Patient Account Number", value: 22, type: "patientAccountNumber" }
  ]

  selectedOption: number;
	selectedType: string;
	showButton: boolean;
  
  constructor(
  private claimsService: ClaimsService, 
  private modalService: NgbModal, 
  private _cfr: ComponentFactoryResolver,
  private formBuilder: FormBuilder
	) { 
		this.createForm();
	}
  
  private createForm() {
    this.searchForm = this.formBuilder.group({
      //category: [this.options[0]],
	  selectItems: this.formBuilder.array([ this.createSelectItems() ]),
	  inputItems: this.formBuilder.array([ this.createInputItem(0, "default") ])
	  //procedureCodeItems: this.formBuilder.array([ this.createProcedureCodeItem() ])
    });
	
	console.log("Inside createForm()");
	var arrayControl = this.searchForm.get('inputItems') as FormArray;
	var item = arrayControl.at(0);
	this.searchForm.get('selectItems').valueChanges.subscribe(data => {
		this.selectedOption = data[0].category.value;
		item.get("type").setValue(data[0].category.type);
	})
	//this.searchForm.get('selectItems').valueChanges.subscribe(changes => console.log('select value has changed:', changes));
  }
  
  // createProcedureCodeItem(): FormGroup {
	  // return this.formBuilder.group({
		// procedureCode: ''
	  // });
  // }
  
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
  
  // addProcedureCodeItem(): void {
	// this.procedureCodeItems = this.searchForm.get('procedureCodeItems') as FormArray;
	// this.procedureCodeItems.push(this.createProcedureCodeItem());
  // }
  
  addInputItem(index:number, type:string): void {
		console.log("addInputItem() index=", index);
	this.inputItems = this.searchForm.get('inputItems') as FormArray;
	this.inputItems.push(this.createInputItem(index, type));
  }
  
  addSelectItems(): void {
		this.selectItems = this.searchForm.get('selectItems') as FormArray;
		this.inputItems = this.searchForm.get('inputItems') as FormArray;
			if (this.selectItems.length < 6) {
				this.selectItems.push(this.createSelectItems());
				this.inputItems = this.searchForm.get('inputItems') as FormArray;
				this.addInputItem(this.inputItems.length, "default");
				var arrayControl = this.getControls(this.searchForm, 'selectItems');
				for(let val of arrayControl) {
						val.get('category').valueChanges.subscribe(data => {
							console.log("Change happened", arrayControl.indexOf(val)+': ', val.get('category').value.name);
							console.log("Need to update items with new type");
							console.log("type is", val.get('category').value.type);
							console.log("select control index=", arrayControl.indexOf(val));
							var inputArrayControl = this.searchForm.get('inputItems') as FormArray;
							var item = inputArrayControl.at(arrayControl.indexOf(val));
							item.get("type").setValue(val.get('category').value.type);
					})
				}
			}
			console.log("select items size=", this.selectItems.length);
			if (this.selectItems.length === 6) {
				this.showButton = false;
			}
  }
  
  clickTab(tab:string) {
	  console.log("clicked " + tab);
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
			//console.log(result);
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
		this.showButton = true;
    this.dtOptions = {
	  searching:false
    };
		this.getAll();
	  this.getModifiedClaims();
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
  // addClaims() {
  //   var newClaims : Claims = {
	//   id: '',
	// 		claimtype: this.claimtype,
	// 		serviceId: this.serviceId,
	// 		receiptDate: this.receiptDate,
	//     fromDate: this.fromDate,
  //     toDate: this.toDate,
  //     placeOfService: this.placeOfService,
  //     providerId: this.providerId,
  //     providerType: this.providerType,
  //     providerSpecialty: this.providerSpecialty,
  //     procedureCode: this.procedureCode,
  //     diagnosisCode: this.diagnosisCode,
  //     networkIndicator: this.networkIndicator,
  //     subscriberId: this.subscriberId,
  //     patientAccountNumber: this.patientAccountNumber,
  //     sccfNumber: this.sccfNumber,
  //     revenueCode: this.revenueCode,
  //     billType: this.billType,
  //     modifier: this.modifier,
  //     planCode: this.planCode,
  //     sfMessageCode: this.sfMessageCode,
  //     pricingMethod: this.pricingMethod,
  //     pricingRule: this.pricingRule,
  //     deliveryMethod: this.deliveryMethod,
  //     inputDate: this.inputDate,
  //     fileName: this.fileName
  //   };

  //   this.claimsService.addClaims(newClaims).subscribe(() => {
  //   this.getAll();
	// this.claimtype = '';
	// this.serviceId = '';
	// this.receiptDate = '';
	// this.fromDate = '';
	// this.toDate = '';
	// this.placeOfService = '';
	// this.providerId = '';
	// this.providerType = '';
	// this.providerSpecialty = '';
	// this.procedureCode = '';
	// this.diagnosisCode = '';
	// this.networkIndicator = '';
	// this.subscriberId = '';
	// this.patientAccountNumber = '';
	// this.sccfNumber = '';
	// this.revenueCode = '';
	// this.billType = '';
	// this.modifier = '';
	// this.planCode = '';
	// this.sfMessageCode = '';
	// this.pricingMethod = '';
	// this.pricingRule = '';
	// this.deliveryMethod = '';
	// this.inputDate = '';
	// this.fileName = '';
  //   });
  // }

  // deleteClaims(claims: Claims) {
  //   this.claimsService.deleteClaims(claims).subscribe(() => {
  //     this.getAll();
  //   })
	// window.location.reload();
  // }
  
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
  // addComponent(){
      // var comp = this._cfr.resolveComponentFactory(SearchComponent);
      // var searchComponent = this.container.createComponent(comp);
      // searchComponent.instance._ref = searchComponent;
  // }
  
  removeObject(index) {
		console.log("removing index->", index);
	  this.inputItems.removeAt(index);
		this.selectItems.removeAt(index);
		if (this.selectItems.length < 6) {
			this.showButton = true;
		}
  }
  
  onSubmit(model: any, isValid: boolean, e: any) {
	  e.preventDefault();
			//alert('Form data are: '+JSON.stringify(model));
			let strFormData = JSON.stringify(model);
			this.searchActiveInstitutionalClaims(strFormData);
	//   let strFormData = JSON.stringify(model);
	//   let jsonObj = JSON.parse(strFormData);
	//   let items = jsonObj.inputItems;
	//   let procedureCodes = [];
	//   let diagnosisCodes = [];
	  
	//   for (let i = 0; i < items.length; i++) {
		  
	// 	  if (items[i].type === 'procedureCode') {
	// 		  var values = Object.values(items[i]);
	// 			console.log(values[0]);
	// 			procedureCodes.push(values[0]);
	// 	  } else if (items[i].type === 'diagnosisCode') {
	// 		  	console.log(values[0]);
	// 			diagnosisCodes.push(values[0]);
	// 	  }
	//   }
	//   console.log("procedureCodes.length=" + procedureCodes.length);
	//   console.log("diagnosisCodes.length=" + diagnosisCodes.length);
   }
}