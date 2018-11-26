import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class ClaimsService {

  constructor(private httpClient: HttpClient) {}

  getClaimsList() {
    return this.httpClient.get(environment.gateway + '/claims');
  }
  
  getClaimsListByIds(claimsId: string) {
    return this.httpClient.get(environment.gateway + '/claims/' + claimsId);
  }

  getSearchResults(search: string) {
    return this.httpClient.get(environment.gateway + '/searchclaims/' + search);
  }

  getModifiedClaimsList() {
    return this.httpClient.get(environment.gateway + '/modifiedclaims');
  }

  // addClaims(claims: Claims) {
  //   return this.httpClient.post(environment.gateway + '/claims', claims);
  // }
  
  addModifiedClaims(claims: Claims) {
    return this.httpClient.post(environment.gateway + '/modifiedclaims', claims);
  }

  addClaims(claimsData: string) {
    return this.httpClient.get(environment.gateway + '/modifiedclaims/' + claimsData);
  }

  deleteClaims(mClaims: ModifiedClaims) {
    return this.httpClient.delete(environment.gateway + '/modifiedclaims/' + mClaims.id);
  }
}

export class Claims {
  id: string;
  claimtype: string;
  serviceId: string;
  receiptDate: string;
  fromDate: string;
  toDate: string;
  providerId: string;
  providerType: string;
  providerSpecialty: string;
  diagnosisCode: string;
  networkIndicator: string;
  subscriberId: string;
  patientAccountNumber: string;
  sccfNumber: string;
  billType: string;
  planCode: string;
  sfMessageCode: string;
  deliveryMethod: string;
  inputDate: string;
  fileName: string;
  createDate: string;
  createdBy: string;
  pSfMessageCode: string;
  pricingMethod: string;
  pricingRule: string;
  procedureCode: string;
  revenueCode: string;
  modifier: string;
  dosFrom: string;
  dosTo: string;
  placeOfService: string;
}

export class ModifiedClaims {
  id:string;
  subscriberId: string;
	originalClaimID: string;
	sccfNumber: string;
	procedureCode: string;
	diagnosisCode: string;
	modifier: string;
	patientAccountNumber: string;
	networkIndicator: string;
	fromDate: string;
	toDate: string;
	status: string;
	dosFrom: string;
	dosTo: string;
	createDate: string;
	createdBy: string;
}