import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { AppRoutingModule } from './app-routing.module';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
import { DataTablesModule } from 'angular-datatables';
import { AppComponent } from './app.component';
import { HomeComponent } from './home/home.component';
import { AuthGuardService } from './auth-guard.service';
//import { AuthService } from './auth.service';
import { CallbackComponent } from './callback/callback.component';
import { ClaimsComponent } from './claims/claims.component';
import { ClaimsService } from './claims.service';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { EditModalComponent } from './edit-modal/edit-modal.component';
import { ErrorModalComponent } from './error-modal/error-modal.component';
import { SearchComponent } from './search/search.component';
//import { TokenInterceptor } from './token.interceptor';

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    CallbackComponent,
    ClaimsComponent,
    EditModalComponent,
    ErrorModalComponent,
    SearchComponent
  ],
  imports: [
    AppRoutingModule,
    BrowserModule,
    FormsModule,
    ReactiveFormsModule,
    HttpClientModule,
	DataTablesModule,
	NgbModule.forRoot()
  ],
  providers: [AuthGuardService, ClaimsService],
  bootstrap: [AppComponent],
  entryComponents: [EditModalComponent, ErrorModalComponent, SearchComponent]
})
export class AppModule { }
