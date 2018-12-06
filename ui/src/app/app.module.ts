import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { RouterModule, Routes } from '@angular/router';
//import { AppRoutingModule } from './app-routing.module';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
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
//import { TokenInterceptor } from './token.interceptor';
import { LoggerModule, NgxLoggerLevel } from 'ngx-logger';
import { environment } from '../environments/environment';
import { LoginComponent } from './login/login.component';
import { AuthenticationService } from './authentication/authentication.service';

const appRoutes: Routes = [
  { path: 'login', component: LoginComponent },
  { path: 'claims', component: ClaimsComponent, canActivate: [AuthGuardService]},
  { path: '',
    redirectTo: '/claims',
    pathMatch: 'full'
  }

];

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    CallbackComponent,
    ClaimsComponent,
    EditModalComponent,
    ErrorModalComponent,
    LoginComponent
  ],
  imports: [
//    AppRoutingModule,
    BrowserModule,
    FormsModule,
    ReactiveFormsModule,
    HttpClientModule,
    NgbModule.forRoot(),
    RouterModule.forRoot(appRoutes),
    LoggerModule.forRoot({serverLoggingUrl: environment.gateway + "/logging", 
                          level: NgxLoggerLevel.DEBUG, 
                          serverLogLevel: NgxLoggerLevel.DEBUG})
  ],
  providers: [AuthGuardService, ClaimsService, AuthenticationService],
  bootstrap: [AppComponent],
  entryComponents: [EditModalComponent, ErrorModalComponent]
})
export class AppModule { }
