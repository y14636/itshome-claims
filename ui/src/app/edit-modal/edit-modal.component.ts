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
  
  constructor(
	public activeModal: NgbActiveModal,
	private claimsService: ClaimsService,
	private formBuilder: FormBuilder
	) { 
		this.createForm();
	}

  ngOnInit() {
	  console.log("selected Institutional Id is " + this.selectedActiveInstitutionalClaimIds);
	  console.log("selected Professional Id is " + this.selectedActiveProfessionalClaimIds);
	  this.getClaimsListByIds();
  }

  getClaimsListByIds() {
    this.claimsService.getClaimsList().subscribe((data: Claims[]) => {
		if (this.selectedActiveInstitutionalClaimIds !== undefined && this.selectedActiveInstitutionalClaimIds.length > 0) {
		    this.selectedActiveInstitutionalClaims = data.filter(claim => claim.id === this.selectedActiveInstitutionalClaimIds[0]);
		    console.log("This institutional claim sub id= " + this.selectedActiveInstitutionalClaims[0].subscriberId);
		    this.subscriberId = this.selectedActiveInstitutionalClaims[0].subscriberId.length > 0 ? this.selectedActiveInstitutionalClaims[0].subscriberId.slice(0, 9) : this.selectedActiveInstitutionalClaims[0].subscriberId;
		    this.suffix = this.selectedActiveInstitutionalClaims[0].subscriberId.length > 9 ? this.selectedActiveInstitutionalClaims[0].subscriberId.slice(-2) : 'N/A';
		    this.patientAccountNumber = this.selectedActiveInstitutionalClaims[0].patientAccountNumber;
		    this.procedureCode = this.selectedActiveInstitutionalClaims[0].procedureCode;
		    this.diagnosisCode = this.selectedActiveInstitutionalClaims[0].diagnosisCode;
		    this.modifier = this.selectedActiveInstitutionalClaims[0].modifier;
		    this.fromDate = this.selectedActiveInstitutionalClaims[0].fromDate;
		    this.toDate = this.selectedActiveInstitutionalClaims[0].toDate;
			this.claimtype = this.selectedActiveInstitutionalClaims[0].claimtype;
			this.serviceId = this.selectedActiveInstitutionalClaims[0].serviceId;
			this.receiptDate = this.selectedActiveInstitutionalClaims[0].receiptDate;
			this.providerType = this.selectedActiveInstitutionalClaims[0].providerType;
			this.providerId = this.selectedActiveInstitutionalClaims[0].providerId;
			this.providerSpecialty = this.selectedActiveInstitutionalClaims[0].providerSpecialty;
			this.sccfNumber = this.selectedActiveInstitutionalClaims[0].sccfNumber;
			this.planCode = this.selectedActiveInstitutionalClaims[0].planCode;
			this.sfMessageCode = this.selectedActiveInstitutionalClaims[0].sfMessageCode;
			this.pricingMethod = this.selectedActiveInstitutionalClaims[0].pricingMethod;
			this.pricingRule = this.selectedActiveInstitutionalClaims[0].pricingRule;
			this.deliveryMethod = this.selectedActiveInstitutionalClaims[0].deliveryMethod;
		} else if (this.selectedActiveProfessionalClaimIds !== undefined && this.selectedActiveProfessionalClaimIds) {
			this.selectedActiveProfessionalClaims = data.filter(claim => claim.id === this.selectedActiveProfessionalClaimIds[0]);
			console.log("This profession claim sub id= " + this.selectedActiveProfessionalClaims[0].subscriberId);
			this.subscriberId = this.selectedActiveProfessionalClaims[0].subscriberId.length > 0 ? this.selectedActiveProfessionalClaims[0].subscriberId.slice(0, 9) : this.selectedActiveProfessionalClaims[0].subscriberId;
			this.suffix = this.selectedActiveProfessionalClaims[0].subscriberId.length > 9 ? this.selectedActiveProfessionalClaims[0].subscriberId.slice(-2) : 'N/A';
			this.patientAccountNumber = this.selectedActiveProfessionalClaims[0].patientAccountNumber;
			this.procedureCode = this.selectedActiveProfessionalClaims[0].procedureCode;
			this.diagnosisCode = this.selectedActiveProfessionalClaims[0].diagnosisCode;
			this.modifier = this.selectedActiveProfessionalClaims[0].modifier;
			this.fromDate = this.selectedActiveProfessionalClaims[0].fromDate;
			this.toDate = this.selectedActiveProfessionalClaims[0].toDate;
			this.claimtype = this.selectedActiveProfessionalClaims[0].claimtype;
			this.serviceId = this.selectedActiveProfessionalClaims[0].serviceId;
			this.receiptDate = this.selectedActiveProfessionalClaims[0].receiptDate;
			this.providerType = this.selectedActiveProfessionalClaims[0].providerType;
			this.providerId = this.selectedActiveProfessionalClaims[0].providerId;
			this.providerSpecialty = this.selectedActiveProfessionalClaims[0].providerSpecialty;
			this.sccfNumber = this.selectedActiveProfessionalClaims[0].sccfNumber;
			this.planCode = this.selectedActiveProfessionalClaims[0].planCode;
			this.sfMessageCode = this.selectedActiveProfessionalClaims[0].sfMessageCode;
			this.pricingMethod = this.selectedActiveProfessionalClaims[0].pricingMethod;
			this.pricingRule = this.selectedActiveProfessionalClaims[0].pricingRule;
			this.deliveryMethod = this.selectedActiveProfessionalClaims[0].deliveryMethod;
		}
	  this.updateClaimForm();
    });
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
	  deliveryMethod: [{value: '', disabled: true}]

    });
  }
  onSubmit(form) {
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
		  subscriberId: form.getRawValue().subscriberId + form.getRawValue().suffix,
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

    this.claimsService.addModifiedClaims(newClaims).subscribe(() => {
		//this.getAll();
		this.claimtype = '';
		//this.serviceId = '';
		//this.receiptDate = '';
		this.fromDate = '';
		this.toDate = '';
		//this.placeOfService = '';
		this.providerId = '';
		this.providerType = '';
		this.providerSpecialty = '';
		this.procedureCode = '';
		this.diagnosisCode = '';
		//this.networkIndicator = '';
		this.subscriberId = '';
		this.suffix = '';
		this.patientAccountNumber = '';
		this.sccfNumber = '';
		//this.revenueCode = '';
		//this.billType = '';
		this.modifier = '';
		this.planCode = '';
		this.sfMessageCode = '';
		this.pricingMethod = '';
		this.pricingRule = '';
		this.deliveryMethod = '';
		//this.inputDate = '';
		//this.fileName = '';
    });
	
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
		  deliveryMethod: this.deliveryMethod
		  
	  });
  }

}
