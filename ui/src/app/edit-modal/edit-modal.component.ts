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
		var selectedId;
		if (this.selectedActiveInstitutionalClaimIds !== undefined && this.selectedActiveInstitutionalClaimIds.length === 1) {
			this.isMultiple = false;
			selectedId = this.selectedActiveInstitutionalClaimIds[0];
		} else if (this.selectedActiveProfessionalClaimIds !== undefined && this.selectedActiveProfessionalClaimIds.length === 1) {
			this.isMultiple = false;
			selectedId = this.selectedActiveProfessionalClaimIds[0];
		}
		if (!this.isMultiple) {
			console.log("selectedId=", selectedId);
			this.claimsService.getClaimsListByIds(selectedId).subscribe((data: Claims[]) => {
					this.isMultiple = false;
					this.populateForm(data.filter(claim =>claim));
					this.updateClaimForm();
			});
		} else {
			this.updateClaimForm();
		}
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
	var subscriberId = selectedClaims[0].subscriberId.trim();
	console.log("length=", subscriberId.length);
	var sub = subscriberId.slice(0, 1);
	if (sub != null && this.isaNumber(sub)) {
		this.subscriberId = subscriberId.slice(0, 9);
	} else {
		this.subscriberId = subscriberId.length > 0 ? subscriberId.slice(0, 12) : subscriberId;
	}
	
	this.suffix = subscriberId.length > 9 ? subscriberId.slice(-2) : 'N/A';
	console.log("suffix=", this.suffix);

	this.patientAccountNumber = selectedClaims[0].patientAccountNumber.trim();
	this.procedureCode = selectedClaims[0].procedureCode.trim();
	this.diagnosisCode = selectedClaims[0].diagnosisCode.trim();
	this.modifier = selectedClaims[0].modifier.trim();
	this.fromDate = selectedClaims[0].fromDate.trim();
	this.toDate = selectedClaims[0].toDate.trim();
	this.claimtype = selectedClaims[0].claimtype.trim();
	this.serviceId = selectedClaims[0].serviceId.trim();
	this.receiptDate = selectedClaims[0].receiptDate.trim();
	this.providerType = selectedClaims[0].providerType.trim();
	this.providerId = selectedClaims[0].providerId.trim();
	this.providerSpecialty = selectedClaims[0].providerSpecialty.trim();
	this.sccfNumber = selectedClaims[0].sccfNumber.trim();
	this.planCode = selectedClaims[0].planCode.trim();
	this.sfMessageCode = selectedClaims[0].sfMessageCode.trim();
	this.pricingMethod = selectedClaims[0].pricingMethod.trim();
	this.pricingRule = selectedClaims[0].pricingRule.trim();
	this.deliveryMethod = selectedClaims[0].deliveryMethod.trim();
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
				providerId: form.getRawValue().providerId,
				providerType: form.getRawValue().providerType,
				providerSpecialty: form.getRawValue().providerSpecialty,
				diagnosisCode: form.getRawValue().diagnosisCode,
				networkIndicator: 'na',
				subscriberId: form.getRawValue().subscriberId + form.getRawValue().suffix,
				patientAccountNumber: form.getRawValue().patientAccountNumber,
				sccfNumber: form.getRawValue().sccfNumber,
				billType: 'na',
				planCode: form.getRawValue().planCode,
				sfMessageCode: form.getRawValue().sfMessageCode,
				deliveryMethod: form.getRawValue().deliveryMethod,
				inputDate: '',
				fileName: '',
				createDate: '',
				createdBy: '',
				pSfMessageCode: '',
				pricingMethod: form.getRawValue().pricingMethod,
				pricingRule: form.getRawValue().pricingRule,
				procedureCode: form.getRawValue().procedureCode,
				revenueCode: 'na',
				modifier: form.getRawValue().modifier,
				dosFrom: '',
				dosTo: '',
				placeOfService: 'na'

				};
		} 
		//alert('Form data are: '+JSON.stringify(this.claimsForm.value));
		let strFormData = JSON.stringify(this.claimsForm.value);
		this.claimsService.addMultipleClaims(strFormData).subscribe(() => {});

		this.activeModal.close(this.claimsForm.value);
  }
  
  closeModal() {
	this.activeModal.close('Modal Closed');
  }

  updateClaimForm() {
	  this.claimsForm.patchValue({
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
