import { Component, OnInit, Output, EventEmitter, Input } from '@angular/core';
import { FormGroup, FormBuilder, FormControl, Validators } from '@angular/forms';
import { NgbActiveModal } from '@ng-bootstrap/ng-bootstrap';
import { ClaimsService, Claims } from '../claims.service';

@Component({
  selector: 'app-edit-modal',
  templateUrl: './edit-modal.component.html',
  styleUrls: ['./edit-modal.component.css']
})
export class EditModalComponent implements OnInit {
	
  @Input ()selectedActiveInstitutionalClaimIds: Array<any>;
  @Input ()selectedActiveProfessionalClaimIds: Array<any>;
  @Input ()subscriberId: 						string;
  @Input ()patientAccountNumber: 				string;
  @Input ()suffix: 								string;
  @Input ()procedureCode: 						string;
  @Input ()diagnosisCode: 						string;
  @Input ()modifier: 							string;
  prefix:										string;
  fromDate: 									string;
  toDate: 										string;
  claimtype: 									string;
  serviceId:									string;
  receiptDate:									string;
  providerType: 								string;
  providerId: 									string;
  providerSpecialty: 							string;
  sccfNumber:									string;
  planCode:										string;
  sfMessageCode:								string;
  pricingMethod:								string;
  pricingRule:									string;
  deliveryMethod:								string;
  claimsForm: 									FormGroup;
  selectedActiveInstitutionalClaims: 			Claims[];
	selectedActiveProfessionalClaims: 			Claims[];
	isMultiple: 									boolean;
  
  constructor(
	public activeModal: NgbActiveModal,
	private claimsService: ClaimsService,
	private formBuilder: FormBuilder
	) { 
		this.createForm();
	}

  ngOnInit() {
		this.isMultiple = true;
	  console.log("selected Institutional Id is " + this.selectedActiveInstitutionalClaimIds);
	  console.log("selected Professional Id is " + this.selectedActiveProfessionalClaimIds);
	  this.getClaimsListByIds();
  }

  getClaimsListByIds() {
    this.claimsService.getClaimsList().subscribe((data: Claims[]) => {
		if (this.selectedActiveInstitutionalClaimIds !== undefined && this.selectedActiveInstitutionalClaimIds.length === 1) {
			this.isMultiple = false;
			this.selectedActiveInstitutionalClaims = data.filter(claim => claim.id === this.selectedActiveInstitutionalClaimIds[0]);
			this.populateForm(this.selectedActiveInstitutionalClaims);
		} else if (this.selectedActiveProfessionalClaimIds !== undefined && this.selectedActiveProfessionalClaimIds.length === 1) {
			this.isMultiple = false;
			this.selectedActiveProfessionalClaims = data.filter(claim => claim.id === this.selectedActiveProfessionalClaimIds[0]);
			this.populateForm(this.selectedActiveProfessionalClaims);
		}
		this.updateClaimForm();
    });
  }
  
  isaNumber(text) {
	var isNumber = false;
	if (!isNaN(parseInt(text, 10))) {
		isNumber = true;
	}
	return isNumber;
  }

  private createForm() {
    this.claimsForm = this.formBuilder.group({
      subscriberId: '',
	  suffix: '',
      patientAccountNumber: '',
	  procedureCode: '',
	  diagnosisCode: '',
	  modifier: '',
	  prefix: [{value: '', disabled: true}],
	  fromDate: [{value: '', disabled: true}],
	  toDate: [{value: '', disabled: true}],
	  claimtype: [{value: '', disabled: true}],
	  serviceId: [{value: '', disabled: true}],
	  receiptDate: [{value: '', disabled: true}],
	  providerType: [{value: '', disabled: true}],
	  providerId: [{value: '', disabled: true}],
	  providerSpecialty: [{value: '', disabled: true}],
	  sccfNumber: [{value: '', disabled: true}],
	  planCode: [{value: '', disabled: true}],
	  sfMessageCode: [{value: '', disabled: true}],
	  pricingMethod: [{value: '', disabled: true}],
	  pricingRule: [{value: '', disabled: true}],
		deliveryMethod: [{value: '', disabled: true}],
		selectedActiveInstitutionalClaimIds: '',
		selectedActiveProfessionalClaimIds: ''

    });
  }

  populateForm(selectedClaims: Claims[]) {
	console.log("This institutional claim sub id= " + selectedClaims[0].subscriberId);
	var sub = selectedClaims[0].subscriberId.slice(0, 1);
	if (sub != null && this.isaNumber(sub)) {
		this.subscriberId = selectedClaims[0].subscriberId.slice(0, 9);
	} else {
		this.subscriberId = selectedClaims[0].subscriberId.length > 0 ? selectedClaims[0].subscriberId.slice(3, 12) : selectedClaims[0].subscriberId;
	}
	this.prefix = selectedClaims[0].subscriberId.length > 12 ? selectedClaims[0].subscriberId.slice(0, 3) : 'N/A';
	this.suffix = selectedClaims[0].subscriberId.length > 9 ? selectedClaims[0].subscriberId.slice(-2) : 'N/A';
	this.patientAccountNumber = selectedClaims[0].patientAccountNumber;
	this.procedureCode = selectedClaims[0].procedureCode;
	this.diagnosisCode = selectedClaims[0].diagnosisCode;
	this.modifier = selectedClaims[0].modifier;
	this.fromDate = selectedClaims[0].fromDate;
	this.toDate = selectedClaims[0].toDate;
	this.claimtype = selectedClaims[0].claimtype;
	this.serviceId = selectedClaims[0].serviceId;
	this.receiptDate = selectedClaims[0].receiptDate;
	this.providerType = selectedClaims[0].providerType;
	this.providerId = selectedClaims[0].providerId;
	this.providerSpecialty = selectedClaims[0].providerSpecialty;
	this.sccfNumber = selectedClaims[0].sccfNumber;
	this.planCode = selectedClaims[0].planCode;
	this.sfMessageCode = selectedClaims[0].sfMessageCode;
	this.pricingMethod = selectedClaims[0].pricingMethod;
	this.pricingRule = selectedClaims[0].pricingRule;
	this.deliveryMethod = selectedClaims[0].deliveryMethod;
  }

  onSubmit(form) {
		if (!this.isMultiple) {
				var newClaims : Claims = {
					id: '',
					claimtype: form.getRawValue().claimtype,
					serviceId: form.getRawValue().serviceId,
					receiptDate: form.getRawValue().receiptDate,
					fromDate: form.getRawValue().fromDate,
					toDate: form.getRawValue().toDate,
					placeOfService: 'na',
					providerId: form.getRawValue().providerId,
					providerType: form.getRawValue().providerType,
					providerSpecialty: form.getRawValue().providerSpecialty,
					procedureCode: form.getRawValue().procedureCode,
					diagnosisCode: form.getRawValue().diagnosisCode,
					networkIndicator: 'na',
					subscriberId: form.getRawValue().prefix + form.getRawValue().subscriberId + form.getRawValue().suffix,
					patientAccountNumber: form.getRawValue().patientAccountNumber,
					sccfNumber: form.getRawValue().sccfNumber,
					revenueCode: 'na',
					billType: 'na',
					modifier: form.getRawValue().modifier,
					planCode: form.getRawValue().planCode,
					sfMessageCode: form.getRawValue().sfMessageCode,
					pricingMethod: form.getRawValue().pricingMethod,
					pricingRule: form.getRawValue().pricingRule,
					deliveryMethod: form.getRawValue().deliveryMethod,
					inputDate: '',
					fileName: 'fromgui'
				};

				// this.claimsService.addModifiedClaims(newClaims).subscribe(() => {
				// //this.getAll();
				// this.claimtype = '';
				// //this.serviceId = '';
				// //this.receiptDate = '';
				// this.fromDate = '';
				// this.toDate = '';
				// //this.placeOfService = '';
				// this.providerId = '';
				// this.providerType = '';
				// this.providerSpecialty = '';
				// this.procedureCode = '';
				// this.diagnosisCode = '';
				// //this.networkIndicator = '';
				// this.subscriberId = '';
				// this.prefix = '';
				// this.suffix = '';
				// this.patientAccountNumber = '';
				// this.sccfNumber = '';
				// //this.revenueCode = '';
				// //this.billType = '';
				// this.modifier = '';
				// this.planCode = '';
				// this.sfMessageCode = '';
				// this.pricingMethod = '';
				// this.pricingRule = '';
				// this.deliveryMethod = '';
				// //this.inputDate = '';
				// //this.fileName = '';
				// });
		} 
		//else {
				//alert('Form data are: '+JSON.stringify(this.claimsForm.value));
				let strFormData = JSON.stringify(this.claimsForm.value);
				this.claimsService.addMultipleClaims(strFormData).subscribe(() => {});
		//}

		this.activeModal.close(this.claimsForm.value);
  }
  
  closeModal() {
	this.activeModal.close('Modal Closed');
  }

  updateClaimForm() {
	  this.claimsForm.patchValue({
		  prefix: this.prefix,
		  subscriberId: this.subscriberId,
		  suffix: this.suffix,
		  patientAccountNumber: this.patientAccountNumber,
		  procedureCode: this.procedureCode,
		  diagnosisCode: this.diagnosisCode,
		  modifier: this.modifier,
		  fromDate: this.fromDate,
		  toDate: this.toDate,
		  claimtype: this.claimtype,
		  serviceId: this.serviceId,
		  receiptDate: this.receiptDate,
		  providerType: this.providerType,
		  providerId: this.providerId,
		  providerSpecialty: this.providerSpecialty,
		  sccfNumber: this.sccfNumber,
		  planCode: this.planCode,
		  sfMessageCode: this.sfMessageCode,
		  pricingMethod: this.pricingMethod,
		  pricingRule: this.pricingRule,
			deliveryMethod: this.deliveryMethod,
			selectedActiveInstitutionalClaimIds: this.selectedActiveInstitutionalClaimIds,
			selectedActiveProfessionalClaimIds: this.selectedActiveProfessionalClaimIds
	  });
  }

}
