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
  
  getSearchResults(search: string) {
    return this.httpClient.post(environment.gateway + '/searchclaims', search);
  }

  getModifiedClaimsList() {
    return this.httpClient.get(environment.gateway + '/modifiedclaims');
  }

  addClaims(claims: Claims) {
    return this.httpClient.post(environment.gateway + '/claims', claims);
  }
  
  addModifiedClaims(claims: Claims) {
    return this.httpClient.post(environment.gateway + '/modifiedclaims', claims);
  }

  deleteClaims(claims: Claims) {
    return this.httpClient.delete(environment.gateway + '/modifiedclaims/' + claims.id);
  }
}

export class Claims {
  id: string;
  claimtype: string;
  serviceId: string;
  receiptDate: string;
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
}