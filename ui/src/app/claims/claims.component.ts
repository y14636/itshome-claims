import { Component, OnInit, ViewChild, ComponentFactoryResolver, ViewContainerRef } from '@angular/core';
import { ClaimsService, Claims } from '../claims.service';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { EditModalComponent } from '../edit-modal/edit-modal.component';
import { ErrorModalComponent } from '../error-modal/error-modal.component';
import { SearchComponent } from '../search/search.component';
import { FormGroup, FormBuilder, FormControl, Validators } from '@angular/forms';

@Component({
  selector: 'app-claims',
  templateUrl: './claims.component.html',
  styleUrls: ['./claims.component.css']
})
export class ClaimsComponent implements OnInit {
	
  viewMode = 'tab1';
  
  searchForm: FormGroup;
   
  dtOptions: DataTables.Settings = {};
  
  activeInstitutionalClaims: Claims[];
  activeProfessionalClaims: Claims[];
  modifiedClaims: Claims[];
  
  selectedActiveInstitutionalClaimIds: Array<any> = [];
  selectedActiveProfessionalClaimIds: Array<any> = [];

  claimtype: string;
  fromDate: string;
  toDate: string;
  placeOfService: string;
  providerId: string;
  providerType: string;
  providerSpecialty: string;
  procedureCode: string;
  diagnosisCode: string;
  networkIndicator: string;
  subscriberId: string;
  patientAccountNumber: string;
  sccfNumber: string;
  revenueCode: string;
  billType: string;
  modifier: string;
  planCode: string;
  sfMessageCode: string;
  pricingMethod: string;
  pricingRule: string;
  deliveryMethod: string;
  inputDate: string;
  fileName: string;
  
  @ViewChild('parent', { read: ViewContainerRef }) container: ViewContainerRef;
  
  options = [{ id: 1, label: 'Category One' }, { id: 2, label: 'Category Two' }];

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
      category: new FormControl(null, Validators.required)
    });
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
    this.dtOptions = {
	  searching:false
    };
    this.getAll();
	this.getModifiedClaims();
  }

  getAll() {
    this.claimsService.getClaimsList().subscribe((data: Claims[]) => {
      this.activeInstitutionalClaims = data.filter(claim => claim.claimtype === '11' || claim.claimtype === '12');
	  this.activeProfessionalClaims = data.filter(claim => claim.claimtype === '20');
    });
  }

  getModifiedClaims() {
	this.claimsService.getModifiedClaimsList().subscribe((data: Claims[]) => {
		this.modifiedClaims = data.filter(claim => claim);
	});	
  }
  
  addClaims() {
    var newClaims : Claims = {
	  id: '',
      claimtype: this.claimtype,
	  fromDate: this.fromDate,
      toDate: this.toDate,
      placeOfService: this.placeOfService,
      providerId: this.providerId,
      providerType: this.providerType,
      providerSpecialty: this.providerSpecialty,
      procedureCode: this.procedureCode,
      diagnosisCode: this.diagnosisCode,
      networkIndicator: this.networkIndicator,
      subscriberId: this.subscriberId,
      patientAccountNumber: this.patientAccountNumber,
      sccfNumber: this.sccfNumber,
      revenueCode: this.revenueCode,
      billType: this.billType,
      modifier: this.modifier,
      planCode: this.planCode,
      sfMessageCode: this.sfMessageCode,
      pricingMethod: this.pricingMethod,
      pricingRule: this.pricingRule,
      deliveryMethod: this.deliveryMethod,
      inputDate: this.inputDate,
      fileName: this.fileName
    };

    this.claimsService.addClaims(newClaims).subscribe(() => {
    this.getAll();
	this.claimtype = '';
	this.fromDate = '';
	this.toDate = '';
	this.placeOfService = '';
	this.providerId = '';
	this.providerType = '';
	this.providerSpecialty = '';
	this.procedureCode = '';
	this.diagnosisCode = '';
	this.networkIndicator = '';
	this.subscriberId = '';
	this.patientAccountNumber = '';
	this.sccfNumber = '';
	this.revenueCode = '';
	this.billType = '';
	this.modifier = '';
	this.planCode = '';
	this.sfMessageCode = '';
	this.pricingMethod = '';
	this.pricingRule = '';
	this.deliveryMethod = '';
	this.inputDate = '';
	this.fileName = '';
    });
  }

  deleteClaims(claims: Claims) {
    this.claimsService.deleteClaims(claims).subscribe(() => {
      this.getAll();
    })
	window.location.reload();
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
  
  addComponent(){
      var comp = this._cfr.resolveComponentFactory(SearchComponent);
      var searchComponent = this.container.createComponent(comp);
      searchComponent.instance._ref = searchComponent;
  }
}